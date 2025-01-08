package biz

import (
	"component-master/model"
	proto "component-master/proto/betting"
	"context"
	"log/slog"
	"receive-service/repository"
)

type BettingBiz interface {
	CreateBetting(ctx context.Context, req *proto.BettingMessageRequest) (*proto.BettingResponse, error)
}

type bettingBiz struct {
	bettingRepository repository.BettingRepository
	redisRepository   repository.RedisRepository
}

func NewBettingBiz(bettingRepository repository.BettingRepository, redisRepository repository.RedisRepository) BettingBiz {
	return &bettingBiz{bettingRepository: bettingRepository, redisRepository: redisRepository}
}

func (b *bettingBiz) CreateBetting(ctx context.Context, req *proto.BettingMessageRequest) (*proto.BettingResponse, error) {
	slog.InfoContext(ctx, "CreateBetting", "req", req)
	return &proto.BettingResponse{
		Code:    0,
		Message: "success",
	}, nil
}

func mapToAccountBalance() *model.AccountBalanceEntity {
	return &model.AccountBalanceEntity{
		Id:        1,
		Balance:   100,
		CreatedAt: 1629886400,
		UpdatedAt: 1629886400,
	}
}
