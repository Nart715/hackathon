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
	go a.rd.CreateAccount(ctx, req)
	return &proto.CreateAccountResponse{Code: 0, Message: "success"}, nil
}

func (a *accountBiz) BalanceChange(ctx context.Context, accountId int64, amount int64) (*proto.BalanceChangeResponse, error) {
	go a.rd.BalanceChange(ctx, accountId, amount)
	return &proto.BalanceChangeResponse{Code: 0, Message: "success"}, nil
}
