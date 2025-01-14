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
}

func NewAccountClient(conf config.GrpcConfigClient) AccountClient {
	if accountGrpcClient == nil {
		InitAccountGrpcClient(conf)
	}
	if accountGrpcClient == nil {
		slog.Error("bettingGrpcClient is nil")
		return nil
	}
	return accountGrpcClient
}

func InitAccountGrpcClient(cgf config.GrpcConfigClient) {
	doOne.Do(func() {
		conn, err := InitConnection(cgf)
		if err != nil {
			slog.Error("InitConnection error", "err", err)
			return
		}
		initConnectionClient(conn)
	})
}

func initConnectionClient(conn *grpc.ClientConn) {
	accountGrpcClient = &accountClient{
		client: rpc.NewAccountServiceClient(conn),
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
	if c == nil || c.client == nil {
		slog.Error("Grpc account client has an error")
		return nil, errors.New("grpc account client is error")
	}
	return c.client.BalanceChange(ctx, req)
}
