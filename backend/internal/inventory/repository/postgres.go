package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/orderpulse/backend/internal/inventory/config"
	"github.com/orderpulse/backend/internal/inventory/domain"
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
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return repo, nil
}

func (r *postgresRepository) initSchema(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS inventory_items (
		id VARCHAR(64) PRIMARY KEY,
		product_id VARCHAR(64) UNIQUE NOT NULL,
		available_quantity INT NOT NULL,
		price NUMERIC(10, 2) NOT NULL,
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);
	INSERT INTO inventory_items (id, product_id, available_quantity, price, updated_at)
	VALUES ('item-101', 'prod-101', 100, 49.99, NOW())
	ON CONFLICT (product_id) DO NOTHING;
	`
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("schema execution failed: %w", err)
	}
	r.logger.Info("PostgreSQL inventory schema initialized with seed data")
	return nil
}

func (r *postgresRepository) GetItem(ctx context.Context, productID string) (*domain.Item, error) {
	query := `
		SELECT id, product_id, available_quantity, price, updated_at
		FROM inventory_items
		WHERE product_id = $1
	`
	row := r.db.QueryRowContext(ctx, query, productID)

	var item domain.Item
	err := row.Scan(&item.ID, &item.ProductID, &item.AvailableQuantity, &item.Price, &item.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrItemNotFound
		}
		return nil, fmt.Errorf("failed to query item: %w", err)
	}

	return &item, nil
}

func (r *postgresRepository) ReserveStock(ctx context.Context, orderID, productID string, quantity int32) (*domain.Item, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	querySelect := `
		SELECT id, product_id, available_quantity, price, updated_at
		FROM inventory_items
		WHERE product_id = $1
		FOR UPDATE
	`
	var item domain.Item
	err = tx.QueryRowContext(ctx, querySelect, productID).Scan(
		&item.ID, &item.ProductID, &item.AvailableQuantity, &item.Price, &item.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrItemNotFound
		}
		return nil, fmt.Errorf("failed to lock item for update: %w", err)
	}

	if item.AvailableQuantity < quantity {
		return nil, domain.ErrInsufficientStock
	}

	newQuantity := item.AvailableQuantity - quantity
	queryUpdate := `
		UPDATE inventory_items
		SET available_quantity = $1, updated_at = NOW()
		WHERE product_id = $2
		RETURNING updated_at
	`
	err = tx.QueryRowContext(ctx, queryUpdate, newQuantity, productID).Scan(&item.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to update stock: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	item.AvailableQuantity = newQuantity
	return &item, nil
}
