package dal

import (
	"context"
	"errors"

	"transfer_system/biz/apperror"
	"transfer_system/biz/model"
	"transfer_system/biz/util"
	"transfer_system/logs"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresTransactionRepository struct {
	db *pgxpool.Pool
}

func NewPostgresTransactionRepository(db *pgxpool.Pool) TransactionRepository {
	return &PostgresTransactionRepository{db: db}
}

func (r *PostgresTransactionRepository) CreateTransaction(ctx context.Context, txID string, sourceAccountID int64, destinationAccountID int64, amount int64) (*model.Transaction, error) {
	dbTx, err := r.db.Begin(ctx)
	if err != nil {
		logs.CtxError(ctx, "failed to begin transaction", err, logs.Fields{
			"transaction_id":         txID,
			"source_account_id":      sourceAccountID,
			"destination_account_id": destinationAccountID,
			"amount":                 amount,
		})
		return nil, err
	}
	defer dbTx.Rollback(ctx)

	if err := lockTransferAccounts(ctx, dbTx, sourceAccountID, destinationAccountID); err != nil {
		logs.CtxError(ctx, "failed to lock transfer accounts", err, logs.Fields{
			"transaction_id":         txID,
			"source_account_id":      sourceAccountID,
			"destination_account_id": destinationAccountID,
		})
		return nil, err
	}

	if err := insertLedgerEntry(ctx, dbTx, txID, sourceAccountID, -amount); err != nil {
		return nil, err
	}
	if err := insertLedgerEntry(ctx, dbTx, txID, destinationAccountID, amount); err != nil {
		return nil, err
	}
	// Potential optimisation: keep this write path ledger-only by inserting the
	// two entries above, then let a background job project cached account
	// balances. That needs a separate overdraft design because the synchronous
	// source balance update is what currently enforces sufficient funds.
	if err := updateAccount(ctx, dbTx, sourceAccountID, -amount); err != nil {
		logs.CtxError(ctx, "failed to update source account", err, logs.Fields{
			"transaction_id":    txID,
			"source_account_id": sourceAccountID,
			"amount":            amount,
		})
		return nil, err
	}
	if err := updateAccount(ctx, dbTx, destinationAccountID, amount); err != nil {
		logs.CtxError(ctx, "failed to update destination account", err, logs.Fields{
			"transaction_id":         txID,
			"destination_account_id": destinationAccountID,
			"amount":                 amount,
		})
		return nil, err
	}
	if err := dbTx.Commit(ctx); err != nil {
		logs.CtxError(ctx, "failed to commit transaction", err, logs.Fields{
			"transaction_id":         txID,
			"source_account_id":      sourceAccountID,
			"destination_account_id": destinationAccountID,
			"amount":                 amount,
		})
		return nil, err
	}

	return &model.Transaction{
		TransactionID:        txID,
		SourceAccountID:      sourceAccountID,
		DestinationAccountID: destinationAccountID,
		Amount:               util.FormatAmount5DP(amount),
	}, nil
}

func lockTransferAccounts(ctx context.Context, tx pgx.Tx, sourceAccountID int64, destinationAccountID int64) error {
	firstAccountID, secondAccountID := sourceAccountID, destinationAccountID
	if secondAccountID < firstAccountID {
		firstAccountID, secondAccountID = secondAccountID, firstAccountID
	}

	if err := lockAccount(ctx, tx, firstAccountID); err != nil {
		return err
	}
	return lockAccount(ctx, tx, secondAccountID)
}

func lockAccount(ctx context.Context, tx pgx.Tx, accountID int64) error {
	var lockedAccountID int64
	err := tx.QueryRow(ctx, `
		SELECT account_id
		FROM accounts
		WHERE account_id = $1
		FOR UPDATE
	`, accountID).Scan(&lockedAccountID)
	if errors.Is(err, pgx.ErrNoRows) {
		logs.CtxError(ctx, "account not found while locking transfer account", apperror.ErrAccountNotFound, logs.Fields{
			"account_id": accountID,
		})
		return apperror.ErrAccountNotFound
	}
	if err != nil {
		logs.CtxError(ctx, "failed to lock transfer account", err, logs.Fields{
			"account_id": accountID,
		})
	}
	return err
}

func insertLedgerEntry(ctx context.Context, tx pgx.Tx, txID string, accountID int64, amount int64) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO ledger_entries (transaction_id, account_id, amount)
		VALUES ($1, $2, $3)
	`, txID, accountID, amount)
	if err != nil {
		mappedErr := mapTransactionError(err)
		logs.CtxError(ctx, "failed to insert ledger entry", mappedErr, logs.Fields{
			"transaction_id": txID,
			"account_id":     accountID,
			"amount":         amount,
		})
		return mappedErr
	}
	return nil
}

func mapTransactionError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" && pgErr.ConstraintName == "ledger_entries_account_id_fkey" {
		return apperror.ErrAccountNotFound
	}
	if errors.As(err, &pgErr) && pgErr.Code == "23514" {
		return apperror.ErrInvalidAmount
	}
	return err
}
