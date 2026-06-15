package dal

import (
	"context"
	"errors"

	"transfer_system/biz/model"

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
	account := &model.Account{}
	err := r.db.QueryRow(ctx, `
		INSERT INTO accounts (account_id, balance)
		VALUES ($1, $2)
		RETURNING account_id, balance
	`, accountID, balance).Scan(&account.AccountID, &account.Balance)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "accounts_pkey" {
			return nil, model.ErrAccountIdAlreadyExists
		}
		return nil, err
	}
	return account, nil
}
