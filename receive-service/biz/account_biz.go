package biz

import (
	proto "component-master/proto/account"
	"component-master/repository"
	"context"

	"golang.org/x/sync/errgroup"
)

var (
	MAP_BALANCE_CHANGE = map[proto.Action]int64{
		proto.Action_DEPOSIT: 1,
		proto.Action_CREDIT:  -1,
	}
)

type AccountBiz interface {
	CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error)
	BalanceChange(ctx context.Context, req *proto.BalanceChangeRequest) (*proto.BalanceChangeResponse, error)
}

type accountBiz struct {
	rd repository.RedisRepository
}

func NewAccountBiz(redisRepo repository.RedisRepository) AccountBiz {
	return &accountBiz{rd: redisRepo}
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
	amount := MAP_BALANCE_CHANGE[proto.Action(req.GetAction())] * req.GetAmount()
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		var err error
		_, err = a.rd.BalanceChange(ctx, req.GetAccountId(), amount, req.GetTransactionId())
		return err
	})
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return &proto.BalanceChangeResponse{Code: 0, Message: "success"}, nil
}
