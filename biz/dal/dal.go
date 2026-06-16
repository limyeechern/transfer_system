package dal

import (
	"context"

	"transfer_system/biz/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, accountID int64, balance int64) (*model.Account, error)
	GetAccount(ctx context.Context, accountID int64) (*model.Account, error)
	UpdateAccount(ctx context.Context, accountID int64, balance int64) error
}

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, txId string, sourceAccountID int64, destinationAccountID int64, amount int64) (*model.Transaction, error)
}

type pgxExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}
