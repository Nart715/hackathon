package service

import (
	"component-master/infra/grpc/client"
	proto "component-master/proto/account"
	"component-master/repository"
	"context"
	"fmt"
	"strconv"
	"time"
)

var ttl = time.Duration(24) * time.Hour

type AccountService interface {
	CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error)
	BalanceChange(ctx context.Context, req *proto.BalanceChangeRequest) (*proto.BalanceChangeResponse, error)
}

type accountService struct {
	grpcclient client.AccountClient
	redis      repository.RedisRepository
}

func NewAccountService(client client.AccountClient, redis repository.RedisRepository) AccountService {
	return &accountService{grpcclient: client, redis: redis}
}

func (s *accountService) CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error) {
	fmt.Println("ACCOUNT SERVICE >> ", req.AccountId)
	return s.grpcclient.CreateAccount(ctx, req)
}

func (s *accountService) BalanceChange(ctx context.Context, req *proto.BalanceChangeRequest) (*proto.BalanceChangeResponse, error) {
	s.redis.GetRedis().Set(context.Background(), strconv.Itoa(int(req.GetTx())), "processing", ttl)
	return s.grpcclient.BalanceChange(ctx, req)
}
