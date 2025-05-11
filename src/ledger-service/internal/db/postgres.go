package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/QuorumMind/ledger-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func InitPostgres() error {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		return fmt.Errorf("POSTGRES_DSN not set")
	}

	var err error
	Pool, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		return err
	}

	return nil
}

func EnsureSchema() error {
	_, err := Pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS ledger_entries (
			id UUID PRIMARY KEY,
			account TEXT NOT NULL,
			amount NUMERIC(18,2) NOT NULL,
			currency VARCHAR(10) NOT NULL,
			transaction_id UUID NOT NULL,
			created_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create ledger_entries table: %w", err)
	}

	return nil
}

func InsertLedgerEntries(evt models.TransactionEvent) error {

	tx, err := Pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	timestamp, _ := time.Parse(time.RFC3339, evt.CreatedAt)

	query := `
		INSERT INTO ledger_entries (id, account, amount, currency, transaction_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	// Debit from sender
	if _, err := tx.Exec(context.Background(), query,
		uuid.New(), evt.FromAccount, -evt.Amount, evt.Currency, evt.ID, timestamp); err != nil {
		return fmt.Errorf("sender ledger entry failed: %w", err)
	}

	// Credit to receiver
	if _, err := tx.Exec(context.Background(), query,
		uuid.New(), evt.ToAccount, evt.Amount, evt.Currency, evt.ID, timestamp); err != nil {
		return fmt.Errorf("receiver ledger entry failed: %w", err)
	}

	return tx.Commit(context.Background())
}
