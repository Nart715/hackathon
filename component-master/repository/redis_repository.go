package repository

import (
	"component-master/infra/redis"
	proto "component-master/proto/account"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	redisv9 "github.com/redis/go-redis/v9"
)

var (
	pipelineSize    = int64(1)
	pipelineTimeout = 100 * time.Millisecond

	count       int
	redisCtx    = context.Background()
	safeCounter *SafeCounter

	AccountIdField = "1"
	BalanceField   = "2"
	CreatedAtField = "3"
	UpdatedAtField = "4"

	TransactionSetKeyPrefix = "account:%d:transactions"

	TransactionExpiryDuration = 24 * time.Hour
)

type SafeCounter struct {
	mu    sync.RWMutex
	count int
}

func (c *SafeCounter) GetCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.count
}

func tryLock(m *SafeCounter, timeout time.Duration) bool {
	ch := make(chan bool, 1)
	go func() {
		m.mu.Lock()
		ch <- true
	}()

	select {
	case <-ch:
		return true
	case <-time.After(timeout):
		return false
	}
}

type RedisRepository interface {
	GetRedis() *redisv9.ClusterClient
	CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error)
	BalanceChange(ctx context.Context, accountId, amount, transactionId int64) (*proto.BalanceChangeResponse, error)
}

type redisRepository struct {
	rd           *redis.RedisClient
	pipeLineTx   redisv9.Pipeliner
	pipelineMu   sync.Mutex
	pipeline     redisv9.Pipeliner
	opCounter    atomic.Int64
	lastExecTime atomic.Int64
}

func NewRedisRepository(rd *redis.RedisClient) RedisRepository {
	safeCounter = &SafeCounter{count: 0}
	repo := &redisRepository{
		rd:         rd,
		pipeline:   rd.GetRd().Pipeline(),
		pipeLineTx: rd.GetRd().TxPipeline(),
	}
	go repo.pipelineFlusher()
	return repo
}
func (r *redisRepository) GetRedis() *redisv9.ClusterClient {
	return r.rd.GetRd()
}

func (r *redisRepository) CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error) {
	slog.Info("Create account")
	account := map[string]interface{}{
		AccountIdField: req.AccountId,
		BalanceField:   int64(0),
		CreatedAtField: time.Now().UnixMicro(),
		UpdatedAtField: time.Now().UnixMicro(),
	}

	// Add operation to pipeline with mutex protection
	r.pipelineMu.Lock()
	accountKey := mapKeyInt64toString(req.AccountId)
	err := r.pipeline.HSet(ctx, accountKey, account).Err()
	r.pipelineMu.Unlock()
	if err != nil {
		return nil, err
	}

	// Increment operation counter

	count := r.opCounter.Add(1)

	// Execute pipeline if threshold reached
	if count%pipelineSize == 0 {
		if err := r.flushPipeline(ctx); err != nil {
			return nil, err
		}
	}

	accountData := proto.AccountData{
		AccountId: account[AccountIdField].(int64),
		Balance:   account[BalanceField].(int64),
		CreatedAt: account[CreatedAtField].(int64),
		UpdatedAt: account[UpdatedAtField].(int64),
	}
	return &proto.CreateAccountResponse{Code: 0, Message: "Successs", Data: &accountData}, nil
}

func (r *redisRepository) ValidateAndRecordTransaction(ctx context.Context, accountId, transactionId int64) (bool, error) {
	transactionKey := fmt.Sprintf(TransactionSetKeyPrefix, accountId)

	r.pipelineMu.Lock()
	defer r.pipelineMu.Unlock()

	tx := r.pipeline.TxPipeline()

	// use lua script for atomic check and add
	script := `
        -- Check if transaction exists
        if redis.call('SISMEMBER', KEYS[1], ARGV[1]) == 1 then
            return 0  -- Transaction already exists
        end

        -- Add transaction and set expiry
        redis.call('SADD', KEYS[1], ARGV[1])
        redis.call('EXPIRE', KEYS[1], ARGV[2])
        return 1  -- New transaction recorded
    `

	// execute script with transaction id and expiry time
	result := tx.Eval(ctx, script,
		[]string{transactionKey},
		transactionId,
		int(TransactionExpiryDuration.Seconds()))

	_, err := tx.Exec(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to execute transaction: %w", err)
	}

	success, err := result.Int64()
	if err != nil {
		return false, fmt.Errorf("failed to get result: %w", err)
	}

	return success == 1, nil
}

func (r *redisRepository) BalanceChange(ctx context.Context, accountId, amount, transactionId int64) (*proto.BalanceChangeResponse, error) {
	isValid, err := r.ValidateAndRecordTransaction(ctx, accountId, transactionId)
	if err != nil {
		return nil, fmt.Errorf("failed to validate transaction: %w", err)
	}

	if !isValid {
		slog.Error("duplicate transaction detected",
			"accountId", accountId,
			"transactionId", transactionId)
		return nil, errors.New("duplicate transaction")
	}

	accountKey := mapKeyInt64toString(accountId)

	r.pipelineMu.Lock()
	defer r.pipelineMu.Unlock()

	tx := r.pipeline.TxPipeline()

	// Lua script to check and increase balance atomically
	script := `
        local balance = redis.call('HGET', KEYS[1], ARGV[1])
        local amount = tonumber(ARGV[2])
        balance = tonumber(balance)

        -- Check if balance would go negative
        if (balance + amount) < 0 then
            return 0  -- Return 0 if balance would be negative
        end

        -- If not negative, increase balance and update timestamp
        redis.call('HINCRBY', KEYS[1], ARGV[1], amount)
        redis.call('HSET', KEYS[1], ARGV[3], ARGV[4])
        return 1  -- Return 1 for successful increase
    `

	// Execute the script with parameters
	result := tx.Eval(ctx, script, []string{accountKey},
		BalanceField,
		amount,
		UpdatedAtField,
		time.Now().UnixMicro())

	_, err = tx.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute transaction: %w", err)
	}

	// Check if operation was successful
	success, err := result.Int64()
	if err != nil {
		return nil, fmt.Errorf("failed to get result: %w", err)
	}

	if success == 0 {
		slog.Error("balance is not enough")
		return nil, errors.New("balance is not enough")
	}

	slog.Info("Balance updated successfully",
		"accountId", accountId,
		"amount", amount)

	return &proto.BalanceChangeResponse{Code: 0, Message: "Success"}, nil
}

func mapKeyInt64toString(key int64) string {
	return strconv.Itoa(int(key))
}

func (r *redisRepository) flushPipeline(ctx context.Context) error {
	r.pipelineMu.Lock()
	defer r.pipelineMu.Unlock()

	if r.opCounter.Load() == 0 {
		return nil
	}

	_, err := r.pipeline.Exec(ctx)
	if err != nil {
		return err
	}

	r.pipeline = r.rd.GetRd().Pipeline()
	r.opCounter.Store(0)
	r.lastExecTime.Store(time.Now().UnixNano())

	return nil
}

func (r *redisRepository) pipelineFlusher() {
	ticker := time.NewTicker(pipelineTimeout)
	defer ticker.Stop()

	for range ticker.C {
		if time.Since(time.Unix(0, r.lastExecTime.Load())) >= pipelineTimeout {
			if err := r.flushPipeline(context.Background()); err != nil {
				// Handle error (log it)
				continue
			}
		}
	}
}
