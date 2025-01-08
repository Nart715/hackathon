package repository

import (
	"component-master/infra/database"
	"component-master/model"
	"context"

	"gorm.io/gorm"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, account *model.AccountBalanceEntity) error
	Deposit(ctx context.Context, accountId int64, amount int64) error
}

type accountRepository struct {
	db *database.DataSource
}

func NewAccountRepository(db *database.DataSource) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) CreateAccount(ctx context.Context, account *model.AccountBalanceEntity) error {
	return r.db.GetDB().WithContext(ctx).Create(account).Error
}

func (r *accountRepository) Deposit(ctx context.Context, accountId int64, amount int64) error {
	return r.db.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		account := &model.AccountBalanceEntity{}
		result := tx.Set("gorm:pessimistic_lock", true).
			First(&account, accountId)
		if result.Error != nil {
			return result.Error
		}

		if err := tx.Where("id = ?", accountId).First(account).Error; err != nil {
			return err
		}

		accountTranaction := &model.AccountTransactionEntity{
			AccountId:    accountId,
			AmountBefore: account.Balance,
			AmountAfter:  account.Balance + amount,
		}

		account.Balance += amount
		if err := tx.Save(account).Error; err != nil {
			return err
		}
		return tx.Create(accountTranaction).Error
	})
}
