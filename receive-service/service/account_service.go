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

func (s accountService) BalanceChange(ctx context.Context, req *proto.BalanceChangeRequest) (*proto.BalanceChangeResponse, error) {
	_, err := s.accountBiz.BalanceChange(ctx, req)
	if err != nil {
		return &proto.BalanceChangeResponse{
			Code:    -1,
			Message: err.Error(),
		}, nil
	}
	return &proto.BalanceChangeResponse{
		Code:    0,
		Message: "success",
	}, nil
}
