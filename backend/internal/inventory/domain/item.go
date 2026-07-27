package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrItemNotFound     = errors.New("item not found")
	ErrInsufficientStock = errors.New("insufficient stock")
)

type Item struct {
	ID                string
	ProductID         string
	AvailableQuantity int32
	Price             float64
	UpdatedAt         time.Time
}

type StockStatus struct {
	IsAvailable       bool
	AvailableQuantity int32
	Price             float64
}

type PostgresRepository interface {
	GetItem(ctx context.Context, productID string) (*Item, error)
	ReserveStock(ctx context.Context, orderID, productID string, quantity int32) (*Item, error)
}

type RedisRepository interface {
	GetStock(ctx context.Context, productID string) (*StockStatus, error)
	SetStock(ctx context.Context, productID string, status *StockStatus) error
	Invalidate(ctx context.Context, productID string) error
}

type Service interface {
	CheckStock(ctx context.Context, productID string, quantity int32) (*StockStatus, error)
	ReserveStock(ctx context.Context, orderID, productID string, quantity int32) (bool, string, error)
}
