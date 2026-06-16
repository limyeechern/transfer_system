package apperror

import "errors"

var (
	ErrInvalidRequest         = errors.New("invalid request")
	ErrInvalidAccount         = errors.New("invalid account id")
	ErrInvalidTransaction     = errors.New("invalid transaction, transfers cannot occur between same accounts")
	ErrAccountIdAlreadyExists = errors.New("account id already exists")
	ErrAccountNotFound        = errors.New("account not found")
	ErrInvalidAmount          = errors.New("invalid amount passed, amount must be more than 0 or at most 5 decimal places")
	ErrInsufficientBalance    = errors.New("source account has insufficient balance to transfer")
	ErrInternalError          = errors.New("internal error")
)
