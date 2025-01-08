package service

import (
	"component-master/infra/grpc/client"
	proto "component-master/proto/message"
	"context"
	"log/slog"
)

type BettingService interface {
	CreateBetting(ctx context.Context, req *proto.BettingMessageRequest) (*proto.BettingResponse, error)
}

type bettingService struct {
	grpcclient client.BettingClient
}

func NewBettingService(client client.BettingClient) BettingService {
	return &bettingService{grpcclient: client}
}

func (s *bettingService) CreateBetting(ctx context.Context, req *proto.BettingMessageRequest) (*proto.BettingResponse, error) {
	slog.Info("betting service - send create betting")
	return s.grpcclient.CreateBetting(ctx, req)
}
