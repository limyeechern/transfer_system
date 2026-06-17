package dal

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"transfer_system/biz/apperror"
	"transfer_system/biz/util"

	"github.com/jackc/pgx/v5/pgxpool"
)

var integrationIDCounter int64

func TestPostgresAccountRepositoryIntegration(t *testing.T) {
	ctx, pool := newIntegrationTestPool(t)
	accountRepo := NewPostgresAccountRepository(pool)

	accountID := uniqueIntegrationAccountID()
	cleanupIntegrationAccounts(t, pool, accountID)

	account, err := accountRepo.CreateAccount(ctx, accountID, 10023344)
	if err != nil {
		t.Fatalf("CreateAccount returned error: %v", err)
	}
	if account.AccountID != accountID {
		t.Fatalf("account id = %d, want %d", account.AccountID, accountID)
	}
	if account.Balance != "100.23344" {
		t.Fatalf("account balance = %q, want %q", account.Balance, "100.23344")
	}

	got, err := accountRepo.GetAccount(ctx, accountID)
	if err != nil {
		t.Fatalf("GetAccount returned error: %v", err)
	}
	if got.AccountID != accountID || got.Balance != "100.23344" {
		t.Fatalf("got account = %+v, want account_id=%d balance=%q", got, accountID, "100.23344")
	}

	_, err = accountRepo.CreateAccount(ctx, accountID, 10023344)
	if !errors.Is(err, apperror.ErrAccountIdAlreadyExists) {
		t.Fatalf("duplicate CreateAccount error = %v, want %v", err, apperror.ErrAccountIdAlreadyExists)
	}

	_, err = accountRepo.GetAccount(ctx, accountID+1)
	if !errors.Is(err, apperror.ErrAccountNotFound) {
		t.Fatalf("missing GetAccount error = %v, want %v", err, apperror.ErrAccountNotFound)
	}
}

func TestPostgresTransactionRepositoryIntegration(t *testing.T) {
	ctx, pool := newIntegrationTestPool(t)
	accountRepo := NewPostgresAccountRepository(pool)
	transactionRepo := NewPostgresTransactionRepository(pool)

	sourceAccountID := uniqueIntegrationAccountID()
	destinationAccountID := sourceAccountID + 1
	cleanupIntegrationAccounts(t, pool, sourceAccountID, destinationAccountID)

	mustCreateIntegrationAccount(t, ctx, accountRepo, sourceAccountID, 10000000)
	mustCreateIntegrationAccount(t, ctx, accountRepo, destinationAccountID, 0)

	txID := uniqueIntegrationTxID("success")
	transaction, err := transactionRepo.CreateTransaction(ctx, txID, sourceAccountID, destinationAccountID, 12345)
	if err != nil {
		t.Fatalf("CreateTransaction returned error: %v", err)
	}
	if transaction.TransactionID != txID {
		t.Fatalf("transaction id = %q, want %q", transaction.TransactionID, txID)
	}
	if transaction.Amount != "0.12345" {
		t.Fatalf("transaction amount = %q, want %q", transaction.Amount, "0.12345")
	}

	assertIntegrationAccountBalance(t, ctx, accountRepo, sourceAccountID, "99.87655")
	assertIntegrationAccountBalance(t, ctx, accountRepo, destinationAccountID, "0.12345")
	assertIntegrationLedgerEntries(t, ctx, pool, txID, map[int64]int64{
		sourceAccountID:      -12345,
		destinationAccountID: 12345,
	})
}

func TestPostgresTransactionRepositoryRollsBackOnInsufficientBalance(t *testing.T) {
	ctx, pool := newIntegrationTestPool(t)
	accountRepo := NewPostgresAccountRepository(pool)
	transactionRepo := NewPostgresTransactionRepository(pool)

	sourceAccountID := uniqueIntegrationAccountID()
	destinationAccountID := sourceAccountID + 1
	cleanupIntegrationAccounts(t, pool, sourceAccountID, destinationAccountID)

	mustCreateIntegrationAccount(t, ctx, accountRepo, sourceAccountID, 500)
	mustCreateIntegrationAccount(t, ctx, accountRepo, destinationAccountID, 0)

	txID := uniqueIntegrationTxID("insufficient")
	_, err := transactionRepo.CreateTransaction(ctx, txID, sourceAccountID, destinationAccountID, 1000)
	if !errors.Is(err, apperror.ErrInsufficientBalance) {
		t.Fatalf("CreateTransaction error = %v, want %v", err, apperror.ErrInsufficientBalance)
	}

	assertIntegrationAccountBalance(t, ctx, accountRepo, sourceAccountID, "0.00500")
	assertIntegrationAccountBalance(t, ctx, accountRepo, destinationAccountID, "0.00000")
	assertIntegrationLedgerCount(t, ctx, pool, txID, 0)
}

func TestPostgresTransactionRepositoryRollsBackOnMissingAccount(t *testing.T) {
	ctx, pool := newIntegrationTestPool(t)
	accountRepo := NewPostgresAccountRepository(pool)
	transactionRepo := NewPostgresTransactionRepository(pool)

	sourceAccountID := uniqueIntegrationAccountID()
	destinationAccountID := sourceAccountID + 1
	cleanupIntegrationAccounts(t, pool, sourceAccountID, destinationAccountID)

	mustCreateIntegrationAccount(t, ctx, accountRepo, sourceAccountID, 10000000)

	txID := uniqueIntegrationTxID("missing-account")
	_, err := transactionRepo.CreateTransaction(ctx, txID, sourceAccountID, destinationAccountID, 1000)
	if !errors.Is(err, apperror.ErrAccountNotFound) {
		t.Fatalf("CreateTransaction error = %v, want %v", err, apperror.ErrAccountNotFound)
	}

	assertIntegrationAccountBalance(t, ctx, accountRepo, sourceAccountID, "100.00000")
	assertIntegrationLedgerCount(t, ctx, pool, txID, 0)
}

func TestPostgresTransactionRepositoryConcurrentOppositeTransfers(t *testing.T) {
	ctx, pool := newIntegrationTestPool(t)
	accountRepo := NewPostgresAccountRepository(pool)
	transactionRepo := NewPostgresTransactionRepository(pool)

	firstAccountID := uniqueIntegrationAccountID()
	secondAccountID := firstAccountID + 1
	cleanupIntegrationAccounts(t, pool, firstAccountID, secondAccountID)

	mustCreateIntegrationAccount(t, ctx, accountRepo, firstAccountID, 10000000)
	mustCreateIntegrationAccount(t, ctx, accountRepo, secondAccountID, 10000000)

	const transferCountPerDirection = 20
	var wg sync.WaitGroup
	errs := make(chan error, transferCountPerDirection*2)

	for i := 0; i < transferCountPerDirection; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			txID := uniqueIntegrationTxID(fmt.Sprintf("a-to-b-%d", i))
			_, err := transactionRepo.CreateTransaction(ctx, txID, firstAccountID, secondAccountID, 100000)
			errs <- err
		}(i)
		go func(i int) {
			defer wg.Done()
			txID := uniqueIntegrationTxID(fmt.Sprintf("b-to-a-%d", i))
			_, err := transactionRepo.CreateTransaction(ctx, txID, secondAccountID, firstAccountID, 100000)
			errs <- err
		}(i)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			t.Fatalf("concurrent transfer returned error: %v", err)
		}
	}

	assertIntegrationAccountBalance(t, ctx, accountRepo, firstAccountID, "100.00000")
	assertIntegrationAccountBalance(t, ctx, accountRepo, secondAccountID, "100.00000")
}

func TestPostgresLedgerEntriesSchemaHasPrimaryKey(t *testing.T) {
	ctx, pool := newIntegrationTestPool(t)

	var primaryKeyCount int
	err := pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM information_schema.table_constraints
		WHERE table_schema = current_schema()
		AND table_name = 'ledger_entries'
		AND constraint_type = 'PRIMARY KEY'
	`).Scan(&primaryKeyCount)
	if err != nil {
		t.Fatalf("query ledger_entries primary key: %v", err)
	}
	if primaryKeyCount != 1 {
		t.Fatalf("ledger_entries primary key count = %d, want 1", primaryKeyCount)
	}
}

func newIntegrationTestPool(t *testing.T) (context.Context, *pgxpool.Pool) {
	t.Helper()

	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("set TEST_DATABASE_URL to run Postgres integration tests")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	t.Cleanup(cancel)

	pool, err := waitForIntegrationPostgres(ctx, databaseURL)
	if err != nil {
		t.Fatalf("connect postgres: %v", err)
	}
	t.Cleanup(pool.Close)

	ensureIntegrationSchema(t, ctx, pool)
	return ctx, pool
}

func waitForIntegrationPostgres(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	var lastErr error
	for {
		pool, err := NewPostgresPool(ctx, databaseURL)
		if err == nil {
			return pool, nil
		}
		lastErr = err

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("%w: %v", ctx.Err(), lastErr)
		case <-time.After(250 * time.Millisecond):
		}
	}
}

func ensureIntegrationSchema(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()

	statements := []string{
		`CREATE TABLE IF NOT EXISTS accounts (
			account_id BIGSERIAL PRIMARY KEY,
			balance BIGINT NOT NULL DEFAULT 0,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS ledger_entries (
			ledger_entry_id BIGSERIAL PRIMARY KEY,
			transaction_id TEXT NOT NULL,
			account_id BIGINT NOT NULL REFERENCES accounts(account_id),
			amount BIGINT NOT NULL CHECK (amount <> 0),
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_ledger_entries_transaction_id ON ledger_entries(transaction_id)`,
		`CREATE INDEX IF NOT EXISTS idx_ledger_entries_account_id ON ledger_entries(account_id)`,
		`CREATE INDEX IF NOT EXISTS idx_ledger_entries_created_at ON ledger_entries(created_at)`,
	}

	for _, statement := range statements {
		if _, err := pool.Exec(ctx, statement); err != nil {
			t.Fatalf("apply integration schema statement %q: %v", statement, err)
		}
	}
}

func uniqueIntegrationAccountID() int64 {
	counter := atomic.AddInt64(&integrationIDCounter, 1)
	return time.Now().UnixMilli()*100000 + counter*100
}

func uniqueIntegrationTxID(name string) string {
	counter := atomic.AddInt64(&integrationIDCounter, 1)
	return fmt.Sprintf("integration-%d-%s-%d", time.Now().UnixNano(), name, counter)
}

func cleanupIntegrationAccounts(t *testing.T, pool *pgxpool.Pool, accountIDs ...int64) {
	t.Helper()

	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		for _, accountID := range accountIDs {
			_, _ = pool.Exec(ctx, `DELETE FROM ledger_entries WHERE account_id = $1`, accountID)
		}
		for _, accountID := range accountIDs {
			_, _ = pool.Exec(ctx, `DELETE FROM accounts WHERE account_id = $1`, accountID)
		}
	}

	cleanup()
	t.Cleanup(cleanup)
}

func mustCreateIntegrationAccount(t *testing.T, ctx context.Context, repo AccountRepository, accountID int64, balance int64) {
	t.Helper()

	if _, err := repo.CreateAccount(ctx, accountID, balance); err != nil {
		t.Fatalf("CreateAccount(%d, %d) returned error: %v", accountID, balance, err)
	}
}

func assertIntegrationAccountBalance(t *testing.T, ctx context.Context, repo AccountRepository, accountID int64, want string) {
	t.Helper()

	account, err := repo.GetAccount(ctx, accountID)
	if err != nil {
		t.Fatalf("GetAccount(%d) returned error: %v", accountID, err)
	}
	if account.Balance != want {
		t.Fatalf("account %d balance = %q, want %q", accountID, account.Balance, want)
	}
}

func assertIntegrationLedgerEntries(t *testing.T, ctx context.Context, pool *pgxpool.Pool, txID string, want map[int64]int64) {
	t.Helper()

	rows, err := pool.Query(ctx, `
		SELECT account_id, amount
		FROM ledger_entries
		WHERE transaction_id = $1
	`, txID)
	if err != nil {
		t.Fatalf("query ledger entries: %v", err)
	}
	defer rows.Close()

	got := map[int64]int64{}
	for rows.Next() {
		var accountID int64
		var amount int64
		if err := rows.Scan(&accountID, &amount); err != nil {
			t.Fatalf("scan ledger entry: %v", err)
		}
		got[accountID] = amount
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate ledger entries: %v", err)
	}

	if len(got) != len(want) {
		t.Fatalf("ledger entry count = %d, want %d; got=%v", len(got), len(want), got)
	}

	var sum int64
	for accountID, wantAmount := range want {
		gotAmount, ok := got[accountID]
		if !ok {
			t.Fatalf("missing ledger entry for account %d; got=%v", accountID, got)
		}
		if gotAmount != wantAmount {
			t.Fatalf("ledger amount for account %d = %d, want %d", accountID, gotAmount, wantAmount)
		}
		sum += gotAmount
	}
	if sum != 0 {
		t.Fatalf("ledger entry sum = %s, want %s", util.FormatAmount5DP(sum), "0.00000")
	}
}

func assertIntegrationLedgerCount(t *testing.T, ctx context.Context, pool *pgxpool.Pool, txID string, want int) {
	t.Helper()

	var got int
	if err := pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM ledger_entries
		WHERE transaction_id = $1
	`, txID).Scan(&got); err != nil {
		t.Fatalf("query ledger entry count: %v", err)
	}
	if got != want {
		t.Fatalf("ledger entry count = %d, want %d", got, want)
	}
}
