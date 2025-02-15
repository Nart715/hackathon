package client

import (
	"component-master/config"
	tproto "component-master/proto/transaction"
	"context"
	"errors"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
)

var (
	grpcClient *transactionClient
)

type TransactionClient interface {
	UpdateTransaction(ctx context.Context, req *tproto.TransactionUpdate) (*tproto.Empty, error)
}

type transactionClient struct {
	client       tproto.TransactionServiceClient
	clientConfig config.GrpcConfigClient
}

func NewTransactionClient(conf config.GrpcConfigClient) TransactionClient {
	InitTransactionGrpcClient(conf)
	return grpcClient
}

func InitTransactionGrpcClient(conf config.GrpcConfigClient) {
	doOne.Do(func() {
		conn, err := InitConnection(conf)
		if err != nil {
			slog.Error("init connection error ", "error", err)
			return
		}
		initConnectionTransactionClient(conn, conf)
	})
}

func initConnectionTransactionClient(conn *grpc.ClientConn, conf config.GrpcConfigClient) {
	grpcClient = &transactionClient{
		client:       tproto.NewTransactionServiceClient(conn),
		clientConfig: conf,
	}
	slog.Info(fmt.Sprintf("InitConnectionClient = %v", grpcClient != nil))
}

func (c *transactionClient) UpdateTransaction(ctx context.Context, req *tproto.TransactionUpdate) (*tproto.Empty, error) {
	if grpcClient == nil {
		InitTransactionGrpcClient(c.clientConfig)
	}
	if c == nil || c.client == nil {
		slog.Error("transaction client is nil")
		return nil, errors.New("transaction client is nil")
	}
	return c.client.UpdateTransaction(ctx, req)
}
