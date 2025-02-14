package client

import (
	"component-master/config"
	rpc "component-master/proto/account"
	"context"
	"errors"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
)

var (
	accountGrpcClient *accountClient
)

type AccountClient interface {
	CreateAccount(ctx context.Context, req *rpc.CreateAccountRequest) (*rpc.CreateAccountResponse, error)
	BalanceChange(ctx context.Context, req *rpc.BalanceChangeRequest) (*rpc.BalanceChangeResponse, error)
}

type accountClient struct {
	client rpc.AccountServiceClient
	conf   config.GrpcConfigClient
}

func NewAccountClient(conf config.GrpcConfigClient) AccountClient {
	InitAccountGrpcClient(conf)
	return accountGrpcClient
}

func InitAccountGrpcClient(cfg config.GrpcConfigClient) {
	doOne.Do(func() {
		conn, err := InitConnection(cfg)

		if err != nil {
			slog.Error("InitConnection error", "err", err)
			return
		}
		initConnectionClient(conn, cfg)
	})
}

func initConnectionClient(conn *grpc.ClientConn, cfg config.GrpcConfigClient) {
	accountGrpcClient = &accountClient{
		client: rpc.NewAccountServiceClient(conn),
		conf:   cfg,
	}
	slog.Info(fmt.Sprintf("InitConnectionClient = %v", accountGrpcClient != nil))
}

func (c *accountClient) CreateAccount(ctx context.Context, req *rpc.CreateAccountRequest) (*rpc.CreateAccountResponse, error) {
	if c == nil || c.client == nil {
		slog.Error("account client is nil")
		return nil, errors.New("account client is nil")
	}
	return c.client.CreateAccount(ctx, req)
}

func (c *accountClient) BalanceChange(ctx context.Context, req *rpc.BalanceChangeRequest) (*rpc.BalanceChangeResponse, error) {
	if accountGrpcClient == nil {
		InitTransactionGrpcClient(c.conf)
	}

	if c == nil || c.client == nil {
		slog.Error("Grpc account client has an error")
		return nil, errors.New("grpc account client is error")
	}
	return c.client.BalanceChange(ctx, req)
}
