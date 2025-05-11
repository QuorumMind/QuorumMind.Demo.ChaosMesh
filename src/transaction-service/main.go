package main

import (
	"log"
	"os"

	"github.com/QuorumMind/transaction-service/internal/db"
	"github.com/QuorumMind/transaction-service/internal/handlers"
	"github.com/QuorumMind/transaction-service/internal/kafka"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env variables")
	}

	log.Println("POSTGRES_DSN =", os.Getenv("POSTGRES_DSN"))

	if err := db.InitPostgres(); err != nil {
		log.Fatal("Postgres init failed:", err)
	}

	if err := db.EnsureSchema(); err != nil {
		log.Fatal("Schema creation failed:", err)
	}

	kafka.InitKafkaWriter()

	r := gin.Default()

	r.POST("/transfer", handlers.TransferHandler)

	log.Println("Transaction Service running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
