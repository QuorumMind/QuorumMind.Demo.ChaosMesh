package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

var Writer *kafka.Writer

func InitKafkaWriter() {
	kafkaURL := os.Getenv("KAFKA_BROKER")
	if kafkaURL == "" {
		kafkaURL = "localhost:9092"
	}

	Writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaURL},
		Topic:    "transactions",
		Balancer: &kafka.LeastBytes{},
	})
}

func SendTransactionEvent(event interface{}) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(time.Now().Format(time.RFC3339Nano)),
		Value: payload,
	}

	log.Println("Sending event to Kafka:", string(payload))
	return Writer.WriteMessages(context.Background(), msg)
}
