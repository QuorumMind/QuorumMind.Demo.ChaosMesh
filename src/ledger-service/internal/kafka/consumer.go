package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/QuorumMind/ledger-service/internal/db"
	"github.com/QuorumMind/ledger-service/internal/models"
	"github.com/segmentio/kafka-go"
)

type TransactionEvent struct {
	ID          string  `json:"id"`
	FromAccount string  `json:"fromAccount"`
	ToAccount   string  `json:"toAccount"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"createdAt"`
}

func StartConsumer() {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:9092"
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{broker},
		Topic:          "transactions",
		GroupID:        "ledger-service",
		MinBytes:       10e3,
		MaxBytes:       10e6,
		MaxWait:        1 * time.Second,
		CommitInterval: 0,
	})

	log.Println("Kafka consumer started...")

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("Kafka read error:", err)
			continue
		}

		var evt models.TransactionEvent

		if err := json.Unmarshal(msg.Value, &evt); err != nil {
			log.Println("Invalid JSON:", err)
			continue
		}

		log.Printf("📩 Received transaction event: %s\n", evt.ID)

		time.Sleep(2 * time.Second) // Let just simulate processing

		if err := db.InsertLedgerEntries(evt); err != nil {
			log.Println("Ledger insert failed:", err)
		} else {
			log.Println("Ledger updated for:", evt.ID)
		}
	}
}
