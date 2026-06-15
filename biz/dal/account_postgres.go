package dal

import (
	"context"
	"errors"

	"transfer_system/biz/apperror"
	"transfer_system/biz/model"
	"transfer_system/biz/util"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresAccountRepository struct {
	db *pgxpool.Pool
}

func NewPostgresAccountRepository(db *pgxpool.Pool) AccountRepository {
	return &PostgresAccountRepository{db: db}
}

func (r *PostgresAccountRepository) CreateAccount(ctx context.Context, accountID int64, balance int64) (*model.Account, error) {
	var storedBalance int64
	account := &model.Account{}
	err := r.db.QueryRow(ctx, `
		INSERT INTO accounts (account_id, balance)
		VALUES ($1, $2)
		RETURNING account_id, balance
	`, accountID, balance).Scan(&account.AccountID, &storedBalance)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "accounts_pkey" {
			return nil, apperror.ErrAccountIdAlreadyExists
		}
		return nil, err
	}
	account.Balance = util.FormatAmount5DP(storedBalance)
	return account, nil
}

func (r *PostgresAccountRepository) GetAccount(ctx context.Context, accountID int64) (*model.Account, error) {
	var storedBalance int64
	account := &model.Account{}
	err := r.db.QueryRow(ctx, `
		SELECT account_id, balance
		FROM accounts
		WHERE account_id = $1
	`, accountID).Scan(&account.AccountID, &storedBalance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrAccountNotFound
		}
		return nil, err
	}
	account.Balance = util.FormatAmount5DP(storedBalance)
	return account, nil
}
