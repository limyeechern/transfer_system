package handler

import "transfer_system/biz/service/create_account"

type App struct {
	Dependencies
}

type Dependencies struct {
	CreateAccountService create_account.CreateAccountService
}

func NewApp(deps Dependencies) *App {
	return &App{
		Dependencies: deps,
	}
}
