package client

import (
	"component-master/config"
	rpc "component-master/proto/message"
	"context"
	"errors"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
)

var (
	bettingGrpcClient *bettingClient
)

type BettingClient interface {
	CreateBetting(ctx context.Context, req *rpc.BettingMessageRequest) (*rpc.BettingResponse, error)
}

type bettingClient struct {
	client rpc.BettingServiceClient
}

func NewBettingClient(conf config.GrpcConfigClient) BettingClient {
	if bettingGrpcClient == nil {
		InitBettingGrpcClient(conf)
	}
	if bettingGrpcClient == nil {
		slog.Error("bettingGrpcClient is nil")
		return nil
	}
	return bettingGrpcClient
}

func InitBettingGrpcClient(cgf config.GrpcConfigClient) {
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
	bettingGrpcClient = &bettingClient{
		client: rpc.NewBettingServiceClient(conn),
	}
	slog.Info(fmt.Sprintf("InitConnectionClient = %v", bettingGrpcClient != nil))
}

func (c *bettingClient) CreateBetting(ctx context.Context, req *rpc.BettingMessageRequest) (*rpc.BettingResponse, error) {
	if c == nil || c.client == nil {
		slog.Error("bettingGrpcClient is nil")
		return nil, errors.New("bettingGrpcClient is nil")
	}
	fmt.Println(c.client != nil)
	slog.Info("Betting client receive request")
	return c.client.CreateBetting(ctx, req)
}
