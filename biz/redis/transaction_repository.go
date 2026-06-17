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

	if err := r.client.Del(ctx, accountKey(sourceAccountID), accountKey(destinationAccountID)).Err(); err != nil {
		logs.CtxError(ctx, "failed to invalidate redis cache keys", err, logs.Fields{
			"source_account_id":      sourceAccountID,
			"destination_account_id": destinationAccountID,
		})
		return transaction, nil
	}
	logs.CtxInfo(ctx, "redis cache keys invalidated", logs.Fields{
		"source_account_id":      sourceAccountID,
		"destination_account_id": destinationAccountID,
	})
	return transaction, nil
}
