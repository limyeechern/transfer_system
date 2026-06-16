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

func updateAccount(ctx context.Context, exec pgxExecutor, accountID int64, amount int64) error {
	if amount > 0 {
		return creditAccount(ctx, exec, accountID, amount)
	}
	if amount < 0 {
		return debitAccount(ctx, exec, accountID, -amount)
	}
	return apperror.ErrInvalidAmount
}

func debitAccount(ctx context.Context, exec pgxExecutor, accountID int64, amount int64) error {
	// Note that we have passed a positive value, so we will do balance - $1 for clearer SQL
	result, err := exec.Exec(ctx, `
		UPDATE accounts
		SET balance = balance - $1,
			updated_at = NOW()
		WHERE account_id = $2
		AND balance >= $1
	`, amount, accountID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		exists, err := accountExists(ctx, exec, accountID)
		if err != nil {
			return err
		}
		if !exists {
			return apperror.ErrAccountNotFound
		}
		return apperror.ErrInsufficientBalance
	}
	return nil
}

func creditAccount(ctx context.Context, exec pgxExecutor, accountID int64, amount int64) error {
	result, err := exec.Exec(ctx, `
		UPDATE accounts
		SET balance = balance + $1,
			updated_at = NOW()
		WHERE account_id = $2
	`, amount, accountID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return apperror.ErrAccountNotFound
	}
	return nil
}

func accountExists(ctx context.Context, exec pgxExecutor, accountID int64) (bool, error) {
	var exists bool
	err := exec.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM accounts
			WHERE account_id = $1
		)
	`, accountID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
