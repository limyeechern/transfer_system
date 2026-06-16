package dal

import (
	"context"
	"transfer_system/biz/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresTransactionRepository struct {
	db *pgxpool.Pool
}

func NewPostgresTransactionRepository(db *pgxpool.Pool) TransactionRepository {
	return &PostgresTransactionRepository{db: db}
}

// remember to do double entry ledger

func (r *PostgresTransactionRepository) CreateTransaction(ctx context.Context, sourceAccountId int64, destinationAccountId int64, amount int64) (*model.Transaction, error) {
	return nil, nil
}
