package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrOrderNotFound    = errors.New("order not found")
	ErrStockReservation = errors.New("failed to reserve inventory stock")
)

type Order struct {
	ID          string
	UserID      string
	ProductID   string
	Quantity    int32
	TotalAmount float64
	Status      string
	CreatedAt   time.Time
}

type OrderEvent struct {
	OrderID     string    `json:"order_id"`
	UserID      string    `json:"user_id"`
	ProductID   string    `json:"product_id"`
	Quantity    int32     `json:"quantity"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"`
	TraceID     string    `json:"trace_id"`
	Timestamp   time.Time `json:"timestamp"`
}

type InvoiceTask struct {
	OrderID     string    `json:"order_id"`
	UserID      string    `json:"user_id"`
	TotalAmount float64   `json:"total_amount"`
	TraceID     string    `json:"trace_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type PostgresRepository interface {
	CreateOrder(ctx context.Context, order *Order) error
	GetOrder(ctx context.Context, orderID string) (*Order, error)
}

type InventoryClient interface {
	ReserveStock(ctx context.Context, orderID, productID string, quantity int32) (bool, float64, string, error)
}

type EventProducer interface {
	PublishOrderPlaced(ctx context.Context, event *OrderEvent) error
	Close() error
}

type TaskPublisher interface {
	PublishInvoiceTask(ctx context.Context, task *InvoiceTask) error
	Close() error
}

type Service interface {
	CreateOrder(ctx context.Context, userID, productID string, quantity int32) (*Order, error)
	GetOrder(ctx context.Context, orderID string) (*Order, error)
}
