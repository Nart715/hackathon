package service

import (
	"component-master/infra/grpc/client"
	proto "component-master/proto/account"
	"context"
	"fmt"
	"time"
)

var ttl = time.Duration(24) * time.Hour

type AccountService interface {
	CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error)
	BalanceChange(ctx context.Context, req *proto.BalanceChangeRequest) (*proto.BalanceChangeResponse, error)
}

type accountService struct {
	grpcclient client.AccountClient
}

func NewAccountService(client client.AccountClient) AccountService {
	return &accountService{grpcclient: client}
}

func (s *accountService) CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error) {
	fmt.Println("ACCOUNT SERVICE >> ", req.AccountId)
	return s.grpcclient.CreateAccount(ctx, req)
}

func (s *accountService) BalanceChange(ctx context.Context, req *proto.BalanceChangeRequest) (*proto.BalanceChangeResponse, error) {
	return s.grpcclient.BalanceChange(ctx, req)
}
