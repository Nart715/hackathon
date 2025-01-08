package model

type AccountBalanceEntity struct {
	Id        int64 `gorm:"primary_key"`
	Balance   int64 `gorm:"column:balance"`
	CreatedAt int64 `gorm:"column:created_at"`
	UpdatedAt int64 `gorm:"column:updated_at"`
}

func (e *AccountBalanceEntity) TableName() string {
	return "account_balances"
}

type AccountTransactionEntity struct {
	Id              int64 `gorm:"primary_key"`
	AccountId       int64 `gorm:"column:account_id"`
	AmountBefore    int64 `gorm:"column:amount_before"`
	Change          int64 `gorm:"column:change"`
	AmountAfter     int64 `gorm:"column:amount_after"`
	TransactionType int64 `gorm:"column:transaction_type"`
	ObjectId        int64 `gorm:"column:object_id"`
	CreatedAt       int64 `gorm:"column:created_at"`
}

func (e *AccountTransactionEntity) TableName() string {
	return "account_transactions"
}
