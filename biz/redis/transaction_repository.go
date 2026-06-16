package redis

import (
	"context"

	"transfer_system/biz/dal"
	"transfer_system/biz/model"
	"transfer_system/logs"

	goredis "github.com/redis/go-redis/v9"
)

type TransactionRepository struct {
	inner  dal.TransactionRepository
	client *goredis.Client
}

func NewTransactionRepository(inner dal.TransactionRepository, client *goredis.Client) dal.TransactionRepository {
	return &TransactionRepository{
		inner:  inner,
		client: client,
	}
}

func (r *TransactionRepository) CreateTransaction(ctx context.Context, txID string, sourceAccountID int64, destinationAccountID int64, amount int64) (*model.Transaction, error) {
	transaction, err := r.inner.CreateTransaction(ctx, txID, sourceAccountID, destinationAccountID, amount)
	if err != nil {
		return nil, err
	}

	r.client.Del(ctx, accountKey(sourceAccountID), accountKey(destinationAccountID))
	logs.CtxInfo(ctx, "redis cache key invalidated", logs.Fields{"sourceAccountID": sourceAccountID, "destinationAccountID": destinationAccountID})
	return transaction, nil
}
