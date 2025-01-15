package repository

import (
	"component-master/infra/redis"
	proto "component-master/proto/account"
	"context"
	"errors"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	redisv9 "github.com/redis/go-redis/v9"
)

var (
	pipelineSize    = int64(1000)
	pipelineTimeout = 100 * time.Millisecond

	count       int
	redisCtx    = context.Background()
	safeCounter *SafeCounter

	AccountIdField = "account_id"
	BalanceField   = "balance"
	CreatedAtField = "created_at"
	UpdatedAtField = "updated_at"
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
	BalanceChange(ctx context.Context, accountId int64, amount int64) (*proto.BalanceChangeResponse, error)
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
	account := map[string]interface{}{
		AccountIdField: req.AccountId,
		BalanceField:   int64(0),
		CreatedAtField: time.Now().UnixMicro(),
		UpdatedAtField: time.Now().UnixMicro(),
	}

	// Add operation to pipeline with mutex protection
	r.pipelineMu.Lock()
	err := r.pipeline.HSet(ctx, mapKeyInt64toString(req.AccountId), account).Err()
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

func (r *redisRepository) BalanceChange(ctx context.Context, accountId int64, amount int64) (*proto.BalanceChangeResponse, error) {
	tx := r.pipeline.TxPipeline()
	accountKey := mapKeyInt64toString(accountId)
	balance, err := tx.HGet(ctx, accountKey, BalanceField).Result()
	if err != nil {
		return nil, err
	}
	num, _ := strconv.Atoi(balance)
	if (num + int(amount)) < 0 {
		return nil, errors.New("balance is not enough")
	}
	num = num + int(amount)
	tx.HSet(ctx, accountKey, BalanceField, num)
	tx.HSet(ctx, accountKey, UpdatedAtField, time.Now().UnixMicro())
	_, err = tx.Exec(ctx)
	if err != nil {
		return nil, errors.New("failed to update balance")
	}
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
