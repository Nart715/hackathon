package biz

import (
	proto "component-master/proto/account"
	"component-master/repository"
	"context"
)

type AccountBiz interface {
	CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error)
	BalanceChange(ctx context.Context, accountId int64, amount int64) (*proto.BalanceChangeResponse, error)
}

type accountBiz struct {
	rd repository.RedisRepository
}

func NewAccountBiz(redisRepo repository.RedisRepository) AccountBiz {
	return &accountBiz{rd: redisRepo}
}

func (a *accountBiz) CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error) {
	return a.rd.CreateAccount(ctx, req)
}

func (a *accountBiz) BalanceChange(ctx context.Context, accountId int64, amount int64) (*proto.BalanceChangeResponse, error) {
	return a.rd.BalanceChange(ctx, accountId, amount)
}
