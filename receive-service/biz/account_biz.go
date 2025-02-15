package biz

import (
	"component-master/config"
	proto "component-master/proto/account"
	"component-master/repository"
	"component-master/worker"
	"context"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

type AccountBiz interface {
	CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error)
	BalanceChange(ctx context.Context, req *proto.BalanceChangeRequest) (*proto.BalanceChangeResponse, error)
	Shutdown() error
}

type accountBiz struct {
	rd          repository.RedisRepository
	workerPool  *worker.BalanceWorkerPool
	workerMutex sync.Mutex
}

func NewAccountBiz(redisRepo repository.RedisRepository, config config.WorkerConfig) AccountBiz {
	fmt.Println(config)
	return &accountBiz{
		rd:         redisRepo,
		workerPool: worker.NewBalanceWorkerPool(config),
	}
}

func (a *accountBiz) ensureWorkerPool(ctx context.Context) {
	a.workerMutex.Lock()
	defer a.workerMutex.Unlock()

	if !a.workerPool.IsRunning.Load() {
		a.workerPool.Start(ctx, a.rd)
	}
}

func (a *accountBiz) CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error) {
	done := make(chan struct {
		resp *proto.CreateAccountResponse
		err  error
	}, 1)

	go func() {
		resp, err := a.rd.CreateAccount(ctx, req)
		done <- struct {
			resp *proto.CreateAccountResponse
			err  error
		}{resp, err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case result := <-done:
		if result.err != nil {
			return nil, result.err
		}
		return result.resp, nil
	}
}

func (a *accountBiz) BalanceChange(ctx context.Context, req *proto.BalanceChangeRequest) (*proto.BalanceChangeResponse, error) {
	a.ensureWorkerPool(ctx)
	g, ctx := errgroup.WithContext(ctx)
	var response *proto.BalanceChangeResponse
	var err error

	g.Go(func() error {
		response, err = a.workerPool.Submit(ctx, req)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	if response != nil {
		return response, nil
	}
	return &proto.BalanceChangeResponse{Code: 0, Message: "success"}, nil
}

func (a *accountBiz) Shutdown() error {
	return a.workerPool.Stop()
}
