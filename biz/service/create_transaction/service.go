package create_transaction

import (
	"context"

	"transfer_system/biz/apperror"
	"transfer_system/biz/dal"
	"transfer_system/biz/model"
	"transfer_system/biz/service"
	"transfer_system/biz/util"
)

type CreateTransactionService interface {
	service.Creator[model.Transaction, model.Transaction]
}

type CreateTransactionServiceImpl struct {
	transactions dal.TransactionRepository
}

func NewCreateTransactionService(repos ...dal.TransactionRepository) CreateTransactionService {
	var transactions dal.TransactionRepository
	if len(repos) > 0 {
		transactions = repos[0]
	}
	return &CreateTransactionServiceImpl{
		transactions: transactions,
	}
}

func (s *CreateTransactionServiceImpl) Create(ctx context.Context, req *model.Transaction) (*model.Transaction, error) {
	if err := s.Validate(ctx, req); err != nil {
		return nil, err
	}

	amount, err := util.ParseAmount5DP(req.Amount)
	if err != nil {
		return nil, err
	}

	return s.transactions.CreateTransaction(ctx, req.SourceAccountID, req.DestinationAccountID, amount)
}

func (s *CreateTransactionServiceImpl) Validate(ctx context.Context, req *model.Transaction) error {
	if req == nil {
		return apperror.ErrInvalidRequest
	}
	if s.transactions == nil {
		return apperror.ErrInternalError
	}
	if req.SourceAccountID <= 0 || req.DestinationAccountID <= 0 {
		return apperror.ErrInvalidAccount
	}
	if _, err := util.ParseAmount5DP(req.Amount); err != nil {
		return err
	}
	return nil
}
