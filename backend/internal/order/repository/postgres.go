package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/orderpulse/backend/internal/order/config"
	"github.com/orderpulse/backend/internal/order/domain"
)

type postgresRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

var _ domain.PostgresRepository = (*postgresRepository)(nil)

func NewPostgresRepository(cfg *config.Config, logger *zap.Logger) (domain.PostgresRepository, error) {
	db, err := sql.Open("postgres", cfg.PostgresDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres connection: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	repo := &postgresRepository{
		db:     db,
		logger: logger,
	}

	if err := repo.initSchema(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to initialize orders schema: %w", err)
	}

	return repo, nil
}

func (r *postgresRepository) initSchema(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS orders (
		id VARCHAR(64) PRIMARY KEY,
		user_id VARCHAR(64) NOT NULL,
		product_id VARCHAR(64) NOT NULL,
		quantity INT NOT NULL,
		total_amount NUMERIC(10, 2) NOT NULL,
		status VARCHAR(32) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW()
	);
	`
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("orders schema execution failed: %w", err)
	}
	r.logger.Info("PostgreSQL orders table schema initialized")
	return nil
}

func (r *postgresRepository) CreateOrder(ctx context.Context, order *domain.Order) error {
	query := `
		INSERT INTO orders (id, user_id, product_id, quantity, total_amount, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query,
		order.ID,
		order.UserID,
		order.ProductID,
		order.Quantity,
		order.TotalAmount,
		order.Status,
		order.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}
	return nil
}

func (r *postgresRepository) GetOrder(ctx context.Context, orderID string) (*domain.Order, error) {
	query := `
		SELECT id, user_id, product_id, quantity, total_amount, status, created_at
		FROM orders
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, orderID)

	var order domain.Order
	err := row.Scan(
		&order.ID,
		&order.UserID,
		&order.ProductID,
		&order.Quantity,
		&order.TotalAmount,
		&order.Status,
		&order.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to query order: %w", err)
	}

	return &order, nil
}
