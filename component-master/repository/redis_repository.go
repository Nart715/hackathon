package repository

import (
	"component-master/infra/redis"
	proto "component-master/proto/account"
	"context"
	"encoding/json"
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
	UpdatedAtField = "3"

	TranactionIdField = "2"
	AmountField       = "3"
	ActionField       = "4"
	CreatedAtField    = "5"

	TransactionSetKeyValidate = "tx:validate"

	AccountSetKeyPrefix = "ac:%d"

	TransactionExpiryDuration = 24 * time.Hour

	TransactionSeqKeyAccount = "ac:%d:tx"
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
	BalanceChange(ctx context.Context, req *proto.BalanceChangeRequest) (*proto.BalanceChangeResponse, error)
	GenerateTransactionSeq(ctx context.Context) (int64, error)
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
		UpdatedAtField: time.Now().UnixMicro(),
	}

	// Add operation to pipeline with mutex protection
	r.pipelineMu.Lock()
	accountKey := mapKeyInt64toString(AccountSetKeyPrefix, req.AccountId)
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
		UpdatedAt: account[UpdatedAtField].(int64),
	}
	return &proto.CreateAccountResponse{Code: 0, Message: "Successs", Data: &accountData}, nil
}

func (r *redisRepository) validateAndRecordTransaction(ctx context.Context, transactionId int64, inputJsonData string) (bool, error) {
	r.pipelineMu.Lock()
	defer r.pipelineMu.Unlock()

	tx := r.pipeline.TxPipeline()

	transactionID := strconv.FormatInt(transactionId, 10)

	// Lua script for atomic check and add
	script := `
        local key = KEYS[1]          -- tx:transactionValidate
        local transactionId = ARGV[1] -- transaction ID
        local jsonData = ARGV[2]      -- JSON data
        local expiry = ARGV[3]        -- expiry time

        -- Check if transaction exists
        if redis.call('HEXISTS', key, transactionId) == 1 then
            return 0  -- Transaction already exists
        end

        -- Add transaction and set expiry
        redis.call('HSET', key, transactionId, jsonData)
        redis.call('EXPIRE', key, expiry)
        return 1  -- New transaction recorded
    `

	result := tx.Eval(ctx, script,
		[]string{"tx:transactionValidate"},
		transactionID,
		inputJsonData,
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

func (r *redisRepository) addTransactionByAccountId(ctx context.Context, accountId int64, transactionId int64, inputJsonData string) (bool, error) {
	r.pipelineMu.Lock()
	defer r.pipelineMu.Unlock()

	tx := r.pipeline.TxPipeline()

	transactionID := strconv.FormatInt(transactionId, 10)
	redisAccountKey := mapKeyInt64toString(TransactionSeqKeyAccount, accountId)

	// Lua script for atomic check and add
	script := `
        local key = KEYS[1]          -- tx:transactionValidate
        local transactionId = ARGV[1] -- transaction ID
        local jsonData = ARGV[2]      -- JSON data
        local expiry = ARGV[3]        -- expiry time

        -- Add transaction and set expiry
        redis.call('HSET', key, transactionId, jsonData)
        redis.call('EXPIRE', key, expiry)
        return 1  -- New transaction recorded
    `

	result := tx.Eval(ctx, script,
		[]string{redisAccountKey},
		transactionID,
		inputJsonData,
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

func (r *redisRepository) BalanceChange(ctx context.Context, input *proto.BalanceChangeRequest) (*proto.BalanceChangeResponse, error) {
	accountId := input.Ac
	amount := input.Am
	transactionId := input.Tx

	mapRequestTransaction := map[string]interface{}{
		AccountIdField:    accountId,
		TranactionIdField: transactionId,
		AmountField:       amount,
		ActionField:       int32(input.Act),
		CreatedAtField:    time.Now().UnixMicro(),
	}
	// Marshal the request to JSON
	jsonData, err := json.Marshal(mapRequestTransaction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction data: %w", err)
	}

	fmt.Println("JsonData: ", string(jsonData))

	inputJsonData := string(jsonData)

	isValid, err := r.validateAndRecordTransaction(ctx, transactionId, inputJsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to validate transaction: %w", err)
	}

	r.addTransactionByAccountId(ctx, accountId, transactionId, inputJsonData)

	if !isValid {
		slog.Error("duplicate transaction detected",
			"accountId", accountId,
			"transactionId", transactionId)
		return nil, errors.New("duplicate transaction")
	}

	accountKey := mapKeyInt64toString(AccountSetKeyPrefix, accountId)

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

func mapKeyInt64toString(prefix string, key int64) string {
	return fmt.Sprintf(prefix, key)
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

func (r *redisRepository) GenerateTransactionSeq(ctx context.Context) (int64, error) {
	r.pipelineMu.Lock()
	defer r.pipelineMu.Unlock()

	tx := r.pipeline.TxPipeline()

	script := `
        -- Get current timestamp in microseconds
        local ts = redis.call('TIME')
        local timestamp = ts[1] * 1000000 + ts[2]

        -- Get sequence number and increment
        local seq = redis.call('INCR', KEYS[1])

        -- Combine timestamp and sequence (last 4 digits)
        -- Format: timestamp(microseconds) + sequence(last 4 digits)
        local sequence = timestamp * 10000 + (seq % 10000)
        return sequence
    `

	result := tx.Eval(ctx, script, []string{""})

	_, err := tx.Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to execute transaction: %w", err)
	}

	sequence, err := result.Int64()
	if err != nil {
		return 0, fmt.Errorf("failed to get sequence: %w", err)
	}

	return sequence, nil
}
