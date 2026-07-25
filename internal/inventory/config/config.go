package config

import (
	"os"
)

type Config struct {
	PostgresDSN string
	RedisAddr   string
	GRPCPort    string
}

func Load() (*Config, error) {
	postgresDSN := os.Getenv("POSTGRES_DSN")
	if postgresDSN == "" {
		postgresDSN = "postgres://appuser:apppassword@localhost:5432/orderpulse_db?sslmode=disable"
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = ":50051"
	}

	return &Config{
		PostgresDSN: postgresDSN,
		RedisAddr:   redisAddr,
		GRPCPort:    grpcPort,
	}, nil
}
