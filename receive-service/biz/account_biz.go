package biz

import (
	proto "component-master/proto/account"
	"context"
	"log/slog"
	"receive-service/repository"
)

type AccountBiz interface {
	CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error)
	DepositAccount(ctx context.Context, req *proto.DepositAccountRequest) (*proto.DepositAccountResponse, error)
}

type accountBiz struct {
	accountRepository repository.AccountRepository
}

func NewAccountBiz(accountRepository repository.AccountRepository) AccountBiz {
	return &accountBiz{accountRepository: accountRepository}
}

func (a *accountBiz) CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error) {
	slog.InfoContext(ctx, "CreateAccount", "req", req)
	return &proto.CreateAccountResponse{
		Code:    0,
		Message: "success",
	}, nil
}

func (a *accountBiz) DepositAccount(ctx context.Context, req *proto.DepositAccountRequest) (*proto.DepositAccountResponse, error) {
	slog.InfoContext(ctx, "DepositAccount", "req", req)
	return &proto.DepositAccountResponse{
		Code:    0,
		Message: "success",
	}, nil
}
