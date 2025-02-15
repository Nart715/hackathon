package worker

import (
	"component-master/config"
	proto "component-master/proto/account"
	"component-master/repository"
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type BalanceWorkerPool struct {
	config    config.WorkerConfig
	JobQueue  chan BalanceJob
	Wg        sync.WaitGroup
	IsRunning atomic.Bool
}

type BalanceJob struct {
	ctx     context.Context
	request *proto.BalanceChangeRequest
	result  chan balanceJobResult
}

type balanceJobResult struct {
	response *proto.BalanceChangeResponse
	err      error
}

func NewBalanceWorkerPool(config config.WorkerConfig) *BalanceWorkerPool {
	if config.NumWorkers <= 0 {
		config.NumWorkers = runtime.NumCPU()
	}
	if config.QueueSize <= 0 {
		config.QueueSize = config.NumWorkers * 100
	}
	if config.JobTimeout <= 0 {
		config.JobTimeout = 30 * time.Second
	}
	if config.ShutdownTimeout <= 0 {
		config.ShutdownTimeout = 60 * time.Second
	}

	return &BalanceWorkerPool{
		config:   config,
		JobQueue: make(chan BalanceJob, config.QueueSize),
	}
}

func (wp *BalanceWorkerPool) Start(ctx context.Context, rd repository.RedisRepository) {
	if !wp.IsRunning.CompareAndSwap(false, true) {
		return // Already running
	}

	// Start workers
	for i := 0; i < wp.config.NumWorkers; i++ {
		wp.Wg.Add(1)
		go wp.Worker(ctx, rd, i)
	}
}

func (wp *BalanceWorkerPool) Stop() error {
	if !wp.IsRunning.CompareAndSwap(true, false) {
		return nil // Already stopped
	}

	ctx, cancel := context.WithTimeout(context.Background(), wp.config.ShutdownTimeout)
	defer cancel()

	close(wp.JobQueue)

	done := make(chan struct{})
	go func() {
		wp.Wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("shutdown timeout exceeded")
	}
}

func (wp *BalanceWorkerPool) Worker(ctx context.Context, rd repository.RedisRepository, id int) {
	defer wp.Wg.Done()

	for {
		select {
		case job, ok := <-wp.JobQueue:
			if !ok {
				return // Worker pool is shutting down
			}

			jobCtx, cancel := context.WithTimeout(job.ctx, wp.config.JobTimeout)
			response, err := rd.BalanceChange(jobCtx, job.request)
			cancel()

			// Send result
			select {
			case job.result <- balanceJobResult{response: response, err: err}:
			case <-ctx.Done():
				return
			}

		case <-ctx.Done():
			return
		}
	}
}

func (wp *BalanceWorkerPool) Submit(ctx context.Context, req *proto.BalanceChangeRequest) (*proto.BalanceChangeResponse, error) {
	if !wp.IsRunning.Load() {
		return nil, fmt.Errorf("worker pool is not running")
	}

	resultChan := make(chan balanceJobResult, 1)

	job := BalanceJob{
		ctx:     ctx,
		request: req,
		result:  resultChan,
	}

	select {
	case wp.JobQueue <- job:
	default:
		return nil, fmt.Errorf("job queue is full")
	}

	select {
	case result := <-resultChan:
		return result.response, result.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
