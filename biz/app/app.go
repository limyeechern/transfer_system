package app

import (
	"context"
	"os"

	"transfer_system/biz/dal"
	"transfer_system/biz/handler"
	"transfer_system/biz/service/create_account"
)

const defaultDatabaseURL = "postgres://transfer_system:transfer_system@127.0.0.1:15432/transfer_system?sslmode=disable"

type Config struct {
	DatabaseURL string
}

func ConfigFromEnv() Config {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = defaultDatabaseURL
	}

	return Config{
		DatabaseURL: databaseURL,
	}
}

func New(ctx context.Context, cfg Config) (*handler.App, func(), error) {
	pool, err := dal.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, nil, err
	}

	accountRepository := dal.NewPostgresAccountRepository(pool)
	createAccountService := create_account.NewCreateAccountService(accountRepository)

	cleanup := func() {
		pool.Close()
	}

	return handler.NewApp(handler.Dependencies{
		CreateAccountService: createAccountService,
	}), cleanup, nil
}
