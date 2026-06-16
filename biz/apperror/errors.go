package apperror

import "errors"

var (
	ErrInvalidRequest         = errors.New("invalid request")
	ErrInvalidAccount         = errors.New("invalid account")
	ErrInvalidTransaction     = errors.New("invalid transaction")
	ErrAccountIdAlreadyExists = errors.New("account id already exists")
	ErrAccountNotFound        = errors.New("account not found")
	ErrInvalidAmount          = errors.New("invalid amount")
	ErrInsufficientBalance    = errors.New("insufficient balance")
	ErrInternalError          = errors.New("internal error")
)
