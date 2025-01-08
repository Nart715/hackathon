package repository

import (
	"component-master/infra/database"
	"component-master/model"
	"context"
	"log/slog"
)

type BettingRepository interface {
	CreateBetting(ctx context.Context, betting *model.AccountBalanceEntity) error
}

type bettingRepository struct {
	db *database.DataSource
}

func NewBettingRepository(db *database.DataSource) BettingRepository {
	return &bettingRepository{db: db}
}

func (r *bettingRepository) CreateBetting(ctx context.Context, betting *model.AccountBalanceEntity) error {
	slog.InfoContext(ctx, "CreateBetting", "betting", betting)
	//return r.db.GetDB().WithContext(ctx).Create(betting).Error
	return nil
}
