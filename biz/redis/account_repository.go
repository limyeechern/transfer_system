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
		logs.CtxInfo(ctx, "redis cache hit", logs.Fields{"accountID": accountID})
		return account, nil
	}

	account, err := r.inner.GetAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}
	r.set(ctx, account)
	logs.CtxInfo(ctx, "redis cache set", logs.Fields{"accountID": accountID})
	return account, nil
}

func (r *AccountRepository) get(ctx context.Context, accountID int64) (*model.Account, bool) {
	raw, err := r.client.Get(ctx, accountKey(accountID)).Result()
	if errors.Is(err, goredis.Nil) || err != nil {
		return nil, false
	}

	account := &model.Account{}
	if err := json.Unmarshal([]byte(raw), account); err != nil {
		return nil, false
	}
	return account, true
}

func (r *AccountRepository) set(ctx context.Context, account *model.Account) {
	payload, err := json.Marshal(account)
	if err != nil {
		return
	}
	r.client.Set(ctx, accountKey(account.AccountID), payload, r.ttl)
}

func (r *AccountRepository) delete(ctx context.Context, accountID int64) {
	r.client.Del(ctx, accountKey(accountID))
}

func accountKey(accountID int64) string {
	return fmt.Sprintf("account:%d", accountID)
}
