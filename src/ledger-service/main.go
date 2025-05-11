package main

import (
	"log"

	"github.com/QuorumMind/ledger-service/internal/db"
	"github.com/QuorumMind/ledger-service/internal/kafka"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	if err := db.InitPostgres(); err != nil {
		log.Fatal("Postgres init failed:", err)
	}

	if err := db.EnsureSchema(); err != nil {
		log.Fatal("Schema creation failed:", err)
	}

	go kafka.StartConsumer()

	log.Println("Ledger service running.")
	select {}
}
