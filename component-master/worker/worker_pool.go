package worker

import (
	proto "component-master/proto/account"
	"context"
	"sync"
)

type Job func(ctx context.Context, input *proto.BalanceChangeRequest) (*proto.BalanceChangeResponse, error)

type WorkerPool struct {
	numWorkers int
	jobs       chan Job
	results    chan error
	wg         sync.WaitGroup
}

func NewWorkerPool(numWorkers int) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		jobs:       make(chan Job, numWorkers),
		results:    make(chan error, numWorkers),
	}
}

func (wp *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(ctx, i)
	}
}

func (wp *WorkerPool) worker(ctx context.Context, id int) {
	defer wp.wg.Done()

	for {
		select {
		case job, ok := <-wp.jobs:
			if !ok {
				return
			}
			// Execute the job function
			err := job()
			wp.results <- err

		case <-ctx.Done():
			return
		}
	}
}

// Submit adds a job to the pool
func (wp *WorkerPool) Submit(job Job) {
	wp.jobs <- job
}

// Close shuts down the worker pool
func (wp *WorkerPool) Close() {
	close(wp.jobs)
	wp.wg.Wait()
	close(wp.results)
}
