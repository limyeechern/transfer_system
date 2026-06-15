package model

import "errors"

var (
	ErrInvalidAccount         = errors.New("invalid account")
	ErrInvalidTransaction     = errors.New("invalid transaction")
	ErrAccountIdAlreadyExists = errors.New("account id already exists")
)
