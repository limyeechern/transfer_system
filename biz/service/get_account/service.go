package get_account

import (
	"context"

	"transfer_system/biz/apperror"
	"transfer_system/biz/dal"
	"transfer_system/biz/model"
	"transfer_system/biz/service"
)

type GetAccountService interface {
	service.Reader[model.GetAccount, model.Account]
}

type GetAccountServiceImpl struct {
	accounts dal.AccountRepository
}

func NewGetAccountService(repos ...dal.AccountRepository) GetAccountService {
	var accounts dal.AccountRepository
	if len(repos) > 0 {
		accounts = repos[0]
	}
	return &GetAccountServiceImpl{
		accounts: accounts,
	}
}

func (s *GetAccountServiceImpl) Read(ctx context.Context, req *model.GetAccount) (*model.Account, error) {
	if err := s.Validate(ctx, req); err != nil {
		return nil, err
	}

	return s.accounts.GetAccount(ctx, req.AccountID)
}

func (s *GetAccountServiceImpl) Validate(ctx context.Context, req *model.GetAccount) error {
	if req == nil {
		return apperror.ErrInvalidRequest
	}
	if s.accounts == nil {
		return apperror.ErrInternalError
	}
	if req.AccountID <= 0 {
		return apperror.ErrInvalidAccount
	}
	return nil
}
