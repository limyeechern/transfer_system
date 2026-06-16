package app

import (
	"context"
	"os"
	"time"

	"transfer_system/biz/dal"
	"transfer_system/biz/handler"
	cacheredis "transfer_system/biz/redis"
	"transfer_system/biz/service/create_account"
	"transfer_system/biz/service/create_transaction"
	"transfer_system/biz/service/get_account"
)

const defaultDatabaseURL = "postgres://transfer_system:transfer_system@127.0.0.1:15432/transfer_system?sslmode=disable"
const defaultRedisAddr = "127.0.0.1:6379"
const defaultAccountCacheTTL = 1 * time.Minute

type Config struct {
	DatabaseURL string
	RedisAddr   string
}

func ConfigFromEnv() Config {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = defaultDatabaseURL
	}
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = defaultRedisAddr
	}

	return Config{
		DatabaseURL: databaseURL,
		RedisAddr:   redisAddr,
	}
}

func New(ctx context.Context, cfg Config) (*handler.App, func(), error) {
	pool, err := dal.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, nil, err
	}

	redisClient := cacheredis.NewClient(cfg.RedisAddr)
	accountRepository := cacheredis.NewAccountRepository(
		dal.NewPostgresAccountRepository(pool),
		redisClient,
		defaultAccountCacheTTL,
	)
	transactionRepository := cacheredis.NewTransactionRepository(
		dal.NewPostgresTransactionRepository(pool),
		redisClient,
	)
	createAccountService := create_account.NewCreateAccountService(accountRepository)
	getAccountService := get_account.NewGetAccountService(accountRepository)
	createTransactionService := create_transaction.NewCreateTransactionService(transactionRepository)

	cleanup := func() {
		redisClient.Close()
		pool.Close()
	}

	return handler.NewApp(handler.Dependencies{
		CreateAccountService:     createAccountService,
		GetAccountService:        getAccountService,
		CreateTransactionService: createTransactionService,
	}), cleanup, nil
}
