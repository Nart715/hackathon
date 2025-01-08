package service

import (
	proto "component-master/proto/betting"
	"context"
	"log/slog"
	"receive-service/biz"
)

type bettingService struct {
	proto.UnimplementedBettingServiceServer
	bettingBiz biz.BettingBiz
}

func NewBettingService(bettingBiz biz.BettingBiz) proto.BettingServiceServer {
	return bettingService{bettingBiz: bettingBiz}
}

func (s bettingService) CreateBetting(ctx context.Context, req *proto.BettingMessageRequest) (*proto.BettingResponse, error) {
	slog.Info("Betting server receive message")
	_, err := s.bettingBiz.CreateBetting(ctx, req)
	if err != nil {
		slog.Error("CreateBetting error", "err", err)
		return &proto.BettingResponse{
			Code:    -1,
			Message: err.Error(),
		}, nil
	}
	return &proto.BettingResponse{
		Message: "Hello, world!",
	}, nil
}
