package worker

import (
	"component-master/config"
	"component-master/infra/grpc/client"
	proto "component-master/proto/account"
	"component-master/util"
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"sent-service/filewriter"
	"sync"
	"sync/atomic"
	"time"
)

type BalanceWorker struct {
	config     config.WorkerConfig
	jobQueue   chan BalanceJob
	wg         sync.WaitGroup
	isRunning  atomic.Bool
	fw         *filewriter.FileWriter
	grpcClient client.AccountClient
}

type BalanceJob struct {
	ctx     context.Context
	request *proto.BalanceChangeRequest
	result  chan BalanceJobResult
}

type BalanceJobResult struct {
	response *proto.BalanceChangeResponse
	err      error
}

func NewBalanceWorker(config config.WorkerConfig, fw *filewriter.FileWriter, grpcClient client.AccountClient) *BalanceWorker {
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

	return &BalanceWorker{
		config:     config,
		jobQueue:   make(chan BalanceJob, config.QueueSize),
		fw:         fw,
		grpcClient: grpcClient,
	}
}

func (w *BalanceWorker) Start(ctx context.Context) {
	if !w.isRunning.CompareAndSwap(false, true) {
		return // Already running
	}

	for i := 0; i < w.config.NumWorkers; i++ {
		w.wg.Add(1)
		go w.worker(ctx, i)
	}
}

func (w *BalanceWorker) Stop() error {
	if !w.isRunning.CompareAndSwap(true, false) {
		return nil // Already stopped
	}

	ctx, cancel := context.WithTimeout(context.Background(), w.config.ShutdownTimeout)
	defer cancel()

	close(w.jobQueue)

	done := make(chan struct{})
	go func() {
		w.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("shutdown timeout exceeded")
	}
}

func (w *BalanceWorker) worker(ctx context.Context, id int) {
	defer w.wg.Done()

	for {
		select {
		case job, ok := <-w.jobQueue:
			if !ok {
				return // Worker pool is shutting down
			}

			w.processJob(job, id)

		case <-ctx.Done():
			return
		}
	}
}

func (w *BalanceWorker) processJob(job BalanceJob, workerID int) {
	req := job.request
	jobCtx := job.ctx

	// Create the processing status line
	line := util.StringPattern("%d_%d_%d_%d_%s",
		req.Tx, req.Ac, req.Am, req.Act, "PROCESSING")

	// Write to file
	if err := w.fw.WriteLine(line); err != nil {
		slog.Error("write processing file error",
			"error", err,
			"worker_id", workerID,
			"tx", req.Tx)
	}

	// Start gRPC call in a goroutine
	go func() {
		grpcCtx := util.ContextwithTimeout()
		_, err := w.grpcClient.BalanceChange(grpcCtx, req)
		if err != nil {
			slog.Error("grpc balance change error",
				"error", err,
				"worker_id", workerID,
				"tx", req.Tx)
		}
	}()

	// Send success response
	select {
	case job.result <- BalanceJobResult{
		response: &proto.BalanceChangeResponse{
			Code:    0,
			Message: "Success",
		},
		err: nil,
	}:
	case <-jobCtx.Done():
		slog.Error("context cancelled while sending response",
			"worker_id", workerID,
			"tx", req.Tx)
	}
}

func (w *BalanceWorker) Submit(ctx context.Context, req *proto.BalanceChangeRequest) (*proto.BalanceChangeResponse, error) {
	if !w.isRunning.Load() {
		return nil, fmt.Errorf("worker pool is not running")
	}

	resultChan := make(chan BalanceJobResult, 1)
	job := BalanceJob{
		ctx:     ctx,
		request: req,
		result:  resultChan,
	}

	// Submit job with timeout
	select {
	case w.jobQueue <- job:
	default:
		return nil, fmt.Errorf("job queue is full")
	}

	// Wait for result
	select {
	case result := <-resultChan:
		return result.response, result.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
