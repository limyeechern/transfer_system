package create_account

import (
	"context"

	"transfer_system/biz/model"
	"transfer_system/biz/service"
)

type CreateAccountService interface {
	service.Creator[model.NewAccount, model.Account]
}

type CreateAccountServiceImpl struct{}

func NewCreateAccountService() CreateAccountService {
	return &CreateAccountServiceImpl{}
}

func (s *CreateAccountServiceImpl) Create(ctx context.Context, req *model.NewAccount) (*model.Account, error) {
	if err := s.Validate(ctx, req); err != nil {
		return nil, err
	}

	return &model.Account{
		AccountID: req.AccountID,
		Balance:   req.InitialBalance,
	}, nil
}

func (s *CreateAccountServiceImpl) Validate(ctx context.Context, req *model.NewAccount) error {
	if req == nil {
		return model.ErrInvalidAccount
	}
	return nil
}
