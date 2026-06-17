package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"transfer_system/biz/dal"
	"transfer_system/biz/model"
	"transfer_system/logs"

	goredis "github.com/redis/go-redis/v9"
)

type AccountRepository struct {
	inner  dal.AccountRepository
	client *goredis.Client
	ttl    time.Duration
}

func NewAccountRepository(inner dal.AccountRepository, client *goredis.Client, ttl time.Duration) dal.AccountRepository {
	return &AccountRepository{
		inner:  inner,
		client: client,
		ttl:    ttl,
	}
}

func (r *AccountRepository) CreateAccount(ctx context.Context, accountID int64, balance int64) (*model.Account, error) {
	account, err := r.inner.CreateAccount(ctx, accountID, balance)
	if err != nil {
		return nil, err
	}
	// All writes should only invalidate
	r.delete(ctx, accountID)
	return account, nil
}

func (r *AccountRepository) GetAccount(ctx context.Context, accountID int64) (*model.Account, error) {
	account, found := r.get(ctx, accountID)
	if found {
		logs.CtxInfo(ctx, "redis cache hit", logs.Fields{"account_id": accountID})
		return account, nil
	}
	logs.CtxInfo(ctx, "redis cache miss", logs.Fields{"account_id": accountID})

	account, err := r.inner.GetAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if r.set(ctx, account) {
		logs.CtxInfo(ctx, "redis cache set", logs.Fields{"account_id": accountID})
	}
	return account, nil
}

func (r *AccountRepository) get(ctx context.Context, accountID int64) (*model.Account, bool) {
	raw, err := r.client.Get(ctx, accountKey(accountID)).Result()
	if errors.Is(err, goredis.Nil) {
		return nil, false
	}
	if err != nil {
		logs.CtxError(ctx, "failed to get account from redis", err, logs.Fields{
			"account_id": accountID,
		})
		return nil, false
	}

	account := &model.Account{}
	if err := json.Unmarshal([]byte(raw), account); err != nil {
		logs.CtxError(ctx, "failed to decode account from redis", err, logs.Fields{
			"account_id": accountID,
		})
		return nil, false
	}
	return account, true
}

func (r *AccountRepository) set(ctx context.Context, account *model.Account) bool {
	payload, err := json.Marshal(account)
	if err != nil {
		logs.CtxError(ctx, "failed to encode account for redis", err, logs.Fields{
			"account_id": account.AccountID,
		})
		return false
	}
	if err := r.client.Set(ctx, accountKey(account.AccountID), payload, r.ttl).Err(); err != nil {
		logs.CtxError(ctx, "failed to set account in redis", err, logs.Fields{
			"account_id": account.AccountID,
		})
		return false
	}
	return true
}

func (r *AccountRepository) delete(ctx context.Context, accountID int64) {
	if err := r.client.Del(ctx, accountKey(accountID)).Err(); err != nil {
		logs.CtxError(ctx, "failed to delete account from redis", err, logs.Fields{
			"account_id": accountID,
		})
		return
	}
	logs.CtxInfo(ctx, "redis cache key invalidated", logs.Fields{"account_id": accountID})
}

func accountKey(accountID int64) string {
	return fmt.Sprintf("account:%d", accountID)
}
