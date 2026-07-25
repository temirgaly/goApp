package config

import (
	"os"
	"strings"
)

type Config struct {
	PostgresDSN           string
	InventoryGRPCEndpoint string
	KafkaBrokers          []string
	RabbitMQURL           string
	GRPCPort              string
}

func Load() (*Config, error) {
	postgresDSN := os.Getenv("POSTGRES_DSN")
	if postgresDSN == "" {
		postgresDSN = "postgres://appuser:apppassword@localhost:5432/orderpulse_db?sslmode=disable"
	}

	inventoryEndpoint := os.Getenv("INVENTORY_GRPC_ENDPOINT")
	if inventoryEndpoint == "" {
		inventoryEndpoint = "localhost:50051"
	}

	brokersStr := os.Getenv("KAFKA_BROKERS")
	var kafkaBrokers []string
	if brokersStr != "" {
		kafkaBrokers = strings.Split(brokersStr, ",")
	} else {
		kafkaBrokers = []string{"localhost:9092"}
	}

	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://guest:guest@localhost:5672/"
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = ":50052"
	}

	return &Config{
		PostgresDSN:           postgresDSN,
		InventoryGRPCEndpoint: inventoryEndpoint,
		KafkaBrokers:          kafkaBrokers,
		RabbitMQURL:           rabbitmqURL,
		GRPCPort:              grpcPort,
	}, nil
}
