package dal

import (
	"context"

	"transfer_system/biz/model"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, accountID int64, balance int64) (*model.Account, error)
	GetAccount(ctx context.Context, accountID int64) (*model.Account, error)
}

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, sourceAccountId int64, destinationAccountId int64, amount int64) (*model.Transaction, error)
}
