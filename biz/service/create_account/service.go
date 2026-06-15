package create_account

import (
	"context"

	"transfer_system/biz/apperror"
	"transfer_system/biz/dal"
	"transfer_system/biz/model"
	"transfer_system/biz/service"
	"transfer_system/biz/util"
)

type CreateAccountService interface {
	service.Creator[model.NewAccount, model.Account]
}

type CreateAccountServiceImpl struct {
	accounts dal.AccountRepository
}

func NewCreateAccountService(repos ...dal.AccountRepository) CreateAccountService {
	var accounts dal.AccountRepository
	if len(repos) > 0 {
		accounts = repos[0]
	}
	return &CreateAccountServiceImpl{
		accounts: accounts,
	}
}

func (s *CreateAccountServiceImpl) Create(ctx context.Context, req *model.NewAccount) (*model.Account, error) {
	if err := s.Validate(ctx, req); err != nil {
		return nil, err
	}

	balance, err := util.ParseAmount5DP(req.InitialBalance)
	if err != nil {
		return nil, err
	}

	return s.accounts.CreateAccount(ctx, req.AccountID, balance)
}

func (s *CreateAccountServiceImpl) Validate(ctx context.Context, req *model.NewAccount) error {
	if req == nil || s.accounts == nil {
		return apperror.ErrInvalidAccount
	}
	if req.AccountID <= 0 {
		return apperror.ErrInvalidAccount
	}
	if _, err := util.ParseAmount5DP(req.InitialBalance); err != nil {
		return err
	}
	return nil
}
