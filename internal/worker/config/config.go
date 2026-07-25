package config

import (
	"os"
)

type Config struct {
	RabbitMQURL string
}

func Load() (*Config, error) {
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://guest:guest@localhost:5672/"
	}

	return &Config{
		RabbitMQURL: rabbitmqURL,
	}, nil
}
