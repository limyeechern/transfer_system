package dal

import (
	"context"

	"transfer_system/biz/model"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, accountID int64, balance int64) (*model.Account, error)
	GetAccount(ctx context.Context, accountID int64) (*model.Account, error)
}
