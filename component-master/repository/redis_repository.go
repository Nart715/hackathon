package repository

import (
	"component-master/infra/grpc/client"
	"component-master/infra/redis"
	proto "component-master/proto/account"
	"component-master/proto/transaction"
	"component-master/util"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	redisv9 "github.com/redis/go-redis/v9"
)

var (
	pipelineSize    = int64(1)
	pipelineTimeout = 100 * time.Millisecond

	MAP_BALANCE_CHANGE = map[int32]int32{
		int32(proto.Action_DEBIT.Number()):  -1,
		int32(proto.Action_CREDIT.Number()): 1,
	}

	count       int
	redisCtx    = context.Background()
	safeCounter *SafeCounter

	AccountIdField = "1"
	BalanceField   = "2"
	UpdatedAtField = "3"

	TranactionIdField         = "2"
	InternalTranactionIdField = "2_1"
	AmountField               = "3"
	ActionField               = "4"
	CreatedAtField            = "5"

	TransactionSetKeyValidate = "tx:validate"

	AccountSetKeyPrefix = "ac:%d"

	TransactionExpiryDuration = 24 * time.Hour

	TransactionSeqKeyAccount = "ac:%d:tx"

	TransactionSequenceAutoIncrementKey = "tx:sequence"

	// Lua script to check and increase balance atomically
	scriptDebit = `
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

	scriptCredit = `
        local balance = redis.call('HGET', KEYS[1], ARGV[1])
        local amount = tonumber(ARGV[2])
        balance = tonumber(balance)

        redis.call('HINCRBY', KEYS[1], ARGV[1], amount)
        redis.call('HSET', KEYS[1], ARGV[3], ARGV[4])
        return 1  -- Return 1 for successful increase
    `

	mapActionScript = map[int32]string{
		0: scriptCredit,
		1: scriptDebit,
	}

	scriptAddTransactionAccount = `
        local key = KEYS[1]          -- tx:transactionValidate
        local transactionId = ARGV[1] -- transaction ID
        local jsonData = ARGV[2]      -- JSON data
        local expiry = ARGV[3]        -- expiry time

        -- Add transaction and set expiry
        redis.call('HSET', key, transactionId, jsonData)
        redis.call('EXPIRE', key, expiry)
        return 1  -- New transaction recorded
    `

	// Lua script for atomic check and add
	scriptValidateAndAddTranaction = `
        local key = KEYS[1]          -- tx:transactionValidate
        local transactionId = ARGV[1] -- transaction ID
        local jsonData = ARGV[2]      -- JSON data
        local expiry = ARGV[3]        -- expiry time

      -- Check if transaction exists
      --  if redis.call('HEXISTS', key, transactionId) == 1 then
      --      return 0  -- Transaction already exists
      --  end

        -- Add transaction and set expiry
        redis.call('HSET', key, transactionId, jsonData)
        redis.call('EXPIRE', key, expiry)
        return 1  -- New transaction recorded
    `

	scriptGenerateTransactionBySequence = `
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
	RedisPubSubChannel = "account-balance-change"
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

type messageProcess func(string, RedisRepository)
type RedisRepository interface {
	GetRedis() *redisv9.ClusterClient
	CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error)
	BalanceChange(ctx context.Context, req *proto.BalanceChangeRequest) (*proto.BalanceChangeResponse, error)
	GenerateTransactionSeq(ctx context.Context) (int64, error)
	PublishMessage(context.Context, string) error
	SubscribeMessage(context.Context, messageProcess) error
	Close() error
}

type redisRepository struct {
	rd                    *redis.RedisClient
	pipeLineTx            redisv9.Pipeliner
	pipelineMu            sync.Mutex
	pipeline              redisv9.Pipeliner
	opCounter             atomic.Int64
	lastExecTime          atomic.Int64
	channel               string
	mutex                 sync.RWMutex
	done                  chan struct{}
	pubsub                *redisv9.PubSub
	isConnected           bool
	transactionGrpcClient client.TransactionClient
	transactionMap        map[int64]bool
}

type Message struct {
	Channel string
	Payload string
}

func NewRedisRepository(rd *redis.RedisClient, txClient client.TransactionClient, channel string) RedisRepository {
	safeCounter = &SafeCounter{count: 0}
	repo := &redisRepository{
		rd:                    rd,
		pipeline:              rd.GetRd().Pipeline(),
		pipeLineTx:            rd.GetRd().TxPipeline(),
		channel:               channel,
		transactionGrpcClient: txClient,
		transactionMap:        make(map[int64]bool),
	}
	go repo.pipelineFlusher()
	return repo
}

func (r *redisRepository) PublishMessage(ctx context.Context, message string) error {
	return r.pipeline.Publish(ctx, r.channel, message).Err()
}

func (c *redisRepository) SubscribeMessage(ctx context.Context, fn messageProcess) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.pubsub = c.GetRedis().Subscribe(ctx, c.channel)

	go c.handleMessages(ctx, fn)
	go c.monitorConnection(ctx)

	return nil
}

func (c *redisRepository) monitorConnection(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-c.done:
			return
		case <-ticker.C:
			if err := c.GetRedis().Ping(ctx).Err(); err != nil {
				log.Printf("Connection check failed: %v", err)
				c.handleDisconnect(ctx)
			}
		}
	}
}

func (c *redisRepository) handleMessages(ctx context.Context, fn messageProcess) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-c.done:
			return
		default:
			msg, err := c.pubsub.Receive(ctx)
			if err != nil {
				log.Printf("Error receiving message: %v", err)
				c.handleDisconnect(ctx)
				continue
			}

			switch m := msg.(type) {
			case *redisv9.Message:
				// Process message
				fn(m.Payload, c)
			case *redisv9.Subscription:
				log.Printf("Subscription message: %s: %s", m.Channel, m.Kind)
			case error:
				log.Printf("Error in subscription: %v", m)
				c.handleDisconnect(ctx)
			}
		}
	}
}

func (c *redisRepository) handleDisconnect(ctx context.Context) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.isConnected {
		return
	}

	c.isConnected = false
	log.Printf("Handling disconnect, attempting to reconnect...")

	// Close existing pubsub
	if c.pubsub != nil {
		_ = c.pubsub.Close()
	}

	// Attempt to reconnect with exponential backoff
	backoff := time.Second
	maxBackoff := time.Second * 30
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Create new pubsub client
			c.pubsub = c.GetRedis().Subscribe(ctx, c.channel)

			// Test connection
			if _, err := c.pubsub.Receive(ctx); err == nil {
				c.isConnected = true
				log.Printf("Successfully reconnected to Redis cluster")
				return
			}

			log.Printf("Reconnection attempt failed, retrying in %v", backoff)
			time.Sleep(backoff)

			// Increase backoff time
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}
	}
}

func (r *redisRepository) GetRedis() *redisv9.ClusterClient {
	return r.rd.GetRd()
}

func (r *redisRepository) CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error) {
	slog.Info("Create account")
	account := map[string]interface{}{
		AccountIdField: req.AccountId,
		BalanceField:   int32(0),
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
		Ac: account[AccountIdField].(int32),
		Bc: account[BalanceField].(int32),
		Up: account[UpdatedAtField].(int64),
	}
	return &proto.CreateAccountResponse{Code: 0, Message: "Successs", Data: &accountData}, nil
}

func (r *redisRepository) recordTransaction(ctx context.Context, transactionId int64, inputJsonData string) (bool, error) {
	r.pipelineMu.Lock()
	defer r.pipelineMu.Unlock()

	tx := r.pipeline.TxPipeline()

	result := tx.Eval(ctx, scriptAddTransactionAccount,
		[]string{"tx:transactionValidate"},
		transactionId,
		inputJsonData,
		int(TransactionExpiryDuration.Seconds()))

	_, err := tx.Exec(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to execute transaction record: %w", err)
	}

	success, err := result.Int64()
	if err != nil {
		return false, fmt.Errorf("failed to get result: %w", err)
	}

	return success == 1, nil
}

func (r *redisRepository) addTransactionByAccountId(ctx context.Context, accountId int32, transactionId int64, inputJsonData string) (bool, error) {
	r.pipelineMu.Lock()
	defer r.pipelineMu.Unlock()

	tx := r.pipeline.TxPipeline()

	redisAccountKey := mapKeyInt64toString(TransactionSeqKeyAccount, accountId)

	// Lua script for atomic check and add

	result := tx.Eval(ctx, scriptAddTransactionAccount,
		[]string{redisAccountKey},
		transactionId,
		inputJsonData,
		int(TransactionExpiryDuration.Seconds()))

	_, err := tx.Exec(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to execute transaction add account: %w", err)
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

	if r.transactionMap[transactionId] {
		return nil, errors.New("duplicate transaction")
	}
	r.transactionMap[transactionId] = true

	mapRequestTransaction := map[string]interface{}{
		AccountIdField:            accountId,
		TranactionIdField:         transactionId,
		InternalTranactionIdField: transactionId,
		AmountField:               amount,
		ActionField:               int32(input.Act),
		CreatedAtField:            time.Now().UnixMicro(),
	}
	// Marshal the request to JSON
	jsonData, err := json.Marshal(mapRequestTransaction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction data: %w", err)
	}

	inputJsonData := string(jsonData)

	_, err = r.recordTransaction(context.Background(), transactionId, inputJsonData)
	if err != nil {
		slog.Error("record transaction error ", "error", err)
	}

	_, err = r.addTransactionByAccountId(context.Background(), accountId, transactionId, inputJsonData)
	if err != nil {
		slog.Error("add transaction error ", "error", err)
	}

	accountKey := mapKeyInt64toString(AccountSetKeyPrefix, accountId)
	r.pipelineMu.Lock()
	defer r.pipelineMu.Unlock()

	tx := r.pipeline.TxPipeline()

	now := time.Now().UnixMicro()
	script := mapActionScript[input.Act]
	amountBalance := MAP_BALANCE_CHANGE[input.Act] * amount
	result := tx.Eval(ctx, script, []string{accountKey}, BalanceField, amountBalance, UpdatedAtField, now)

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

	go r.updatebackTransaction(context.Background(), input)
	return &proto.BalanceChangeResponse{Code: 0, Message: "Success"}, nil
}

func (r *redisRepository) updatebackTransaction(ctx context.Context, input *proto.BalanceChangeRequest) {
	r.pipelineMu.Lock()
	defer r.pipelineMu.Unlock()

	tx := r.pipeline.TxPipeline()

	accountKey := mapKeyInt64toString(AccountSetKeyPrefix, input.Ac)
	balance, _ := tx.HGet(ctx, accountKey, BalanceField).Result()
	balanceNum := util.StringToInt64(balance)

	change := int64(MAP_BALANCE_CHANGE[input.Act]) * int64(input.Am)
	balanceBefore := balanceNum - change
	update := &transaction.TransactionUpdate{
		Ac:     input.Ac,
		Tx:     input.Tx,
		Itx:    input.Tx,
		Bab:    balanceBefore,
		Change: int64(input.Am),
		Baf:    balanceNum,
		Type:   input.Act,
		Status: "SUCCESS",
	}
	r.transactionGrpcClient.UpdateTransaction(context.Background(), update)
}

func (r *redisRepository) buildMessagePublish(accountId int32, transactionId int64, amount int32, t int64) {
	slog.Info(util.StringPattern("publish message to redis channel %d, %s", transactionId, r.channel))
	msg := &proto.AccountData{
		Ac:  accountId,
		Tx:  transactionId,
		Itx: transactionId,
		Bc:  amount,
		Up:  t,
	}
	data, _ := json.Marshal(msg)
	r.PublishMessage(context.Background(), string(data))
	r.PublishMessage(context.Background(), "{}")
}

func mapKeyInt64toString(prefix string, key int32) string {
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

	result := tx.Eval(ctx, scriptGenerateTransactionBySequence, []string{TransactionSequenceAutoIncrementKey})

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

func (c *redisRepository) Close() error {
	if c.pubsub != nil {
		if err := c.pubsub.Close(); err != nil {
			return fmt.Errorf("error closing pubsub: %v", err)
		}
	}
	if c.rd != nil {
		return c.GetRedis().Close()
	}
	return nil
}
