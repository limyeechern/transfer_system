CREATE TABLE IF NOT EXISTS accounts (
    account_id BIGSERIAL PRIMARY KEY,
    balance BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS transactions (
    transaction_id BIGSERIAL PRIMARY KEY,
    source_account_id BIGINT NOT NULL REFERENCES accounts(account_id),
    destination_account_id BIGINT NOT NULL REFERENCES accounts(account_id),
    amount BIGINT NOT NULL CHECK (amount > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_transactions_source_account_id
ON transactions(source_account_id);

CREATE INDEX IF NOT EXISTS idx_transactions_destination_account_id
ON transactions(destination_account_id);