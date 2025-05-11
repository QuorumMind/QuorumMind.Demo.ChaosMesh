CREATE TABLE IF NOT EXISTS ledger_entries (
    id UUID PRIMARY KEY,
    account TEXT NOT NULL,
    amount NUMERIC(18,2) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    transaction_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL
);
