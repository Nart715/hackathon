package service

import (
	proto "component-master/proto/account"
	"context"
	"receive-service/biz"
)

type accountService struct {
	proto.UnimplementedAccountServiceServer
	accountBiz biz.AccountBiz
}

func NewAccountService(accountBiz biz.AccountBiz) proto.AccountServiceServer {
	return &accountService{accountBiz: accountBiz}
}

func (s accountService) CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error) {
	_, err := s.accountBiz.CreateAccount(ctx, req)
	if err != nil {
		return &proto.CreateAccountResponse{
			Code:    -1,
			Message: err.Error(),
		}, nil
	}
	return &proto.CreateAccountResponse{
		Code:    0,
		Message: "success",
	}, nil
}

func (s accountService) DepositAccount(ctx context.Context, req *proto.DepositAccountRequest) (*proto.DepositAccountResponse, error) {
	_, err := s.accountBiz.DepositAccount(ctx, req)
	if err != nil {
		return &proto.DepositAccountResponse{
			Code:    -1,
			Message: err.Error(),
		}, nil
	}
	return &proto.DepositAccountResponse{
		Code:    0,
		Message: "success",
	}, nil
}
