package handler

import (
	"transfer_system/biz/service/create_account"
	"transfer_system/biz/service/create_transaction"
	"transfer_system/biz/service/get_account"
)

type App struct {
	Dependencies
}

type Dependencies struct {
	CreateAccountService     create_account.CreateAccountService
	GetAccountService        get_account.GetAccountService
	CreateTransactionService create_transaction.CreateTransactionService
}

func NewApp(deps Dependencies) *App {
	return &App{
		Dependencies: deps,
	}
}
