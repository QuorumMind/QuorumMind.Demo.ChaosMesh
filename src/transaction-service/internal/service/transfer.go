package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/QuorumMind/transaction-service/internal/db"
	"github.com/QuorumMind/transaction-service/internal/kafka"
	"github.com/QuorumMind/transaction-service/internal/models"

	"github.com/google/uuid"
)

func ProcessTransfer(req models.TransferRequest) error {
	txId := uuid.New().String()
	query := `
		INSERT INTO transactions (id, from_account, to_account, amount, currency, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := db.Pool.Exec(
		context.Background(),
		query,
		txId,
		req.FromAccount,
		req.ToAccount,
		req.Amount,
		req.Currency,
		"pending",
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to insert transaction: %w", err)
	}

	fmt.Println("Inserted transaction:", txId)

	event := map[string]interface{}{
		"id":          txId,
		"fromAccount": req.FromAccount,
		"toAccount":   req.ToAccount,
		"amount":      req.Amount,
		"currency":    req.Currency,
		"status":      "pending",
		"createdAt":   time.Now().Format(time.RFC3339),
	}

	if err := kafka.SendTransactionEvent(event); err != nil {
		log.Println("⚠️ Failed to send Kafka event:", err)
	}

	return nil
}
