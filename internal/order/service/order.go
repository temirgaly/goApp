package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/orderpulse/internal/order/domain"
)

type orderService struct {
	postgresRepo    domain.PostgresRepository
	inventoryClient domain.InventoryClient
	kafkaProducer   domain.EventProducer
	rabbitPublisher domain.TaskPublisher
	logger          *zap.Logger
}

var _ domain.Service = (*orderService)(nil)

func NewOrderService(
	pgRepo domain.PostgresRepository,
	invClient domain.InventoryClient,
	kProducer domain.EventProducer,
	rPublisher domain.TaskPublisher,
	logger *zap.Logger,
) domain.Service {
	return &orderService{
		postgresRepo:    pgRepo,
		inventoryClient: invClient,
		kafkaProducer:   kProducer,
		rabbitPublisher: rPublisher,
		logger:          logger,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, userID, productID string, quantity int32) (*domain.Order, error) {
	orderID := fmt.Sprintf("order-%s", uuid.New().String()[:8])
	traceID := fmt.Sprintf("trace-%s", uuid.New().String()[:8])

	s.logger.Info("Initiating CreateOrder workflow",
		zap.String("order_id", orderID),
		zap.String("trace_id", traceID),
		zap.String("user_id", userID),
		zap.String("product_id", productID),
		zap.Int32("quantity", quantity),
	)

	// Step 1: Call Inventory Service over gRPC
	success, unitPrice, msg, err := s.inventoryClient.ReserveStock(ctx, orderID, productID, quantity)
	if err != nil {
		s.logger.Error("Inventory Service gRPC call failed", zap.String("trace_id", traceID), zap.Error(err))
		return nil, fmt.Errorf("failed to reserve stock: %w", err)
	}

	if !success {
		s.logger.Warn("Stock reservation rejected by Inventory Service",
			zap.String("trace_id", traceID),
			zap.String("reason", msg),
		)
		return nil, domain.ErrStockReservation
	}

	totalAmount := unitPrice * float64(quantity)
	now := time.Now()

	order := &domain.Order{
		ID:          orderID,
		UserID:      userID,
		ProductID:   productID,
		Quantity:    quantity,
		TotalAmount: totalAmount,
		Status:      "COMPLETED",
		CreatedAt:   now,
	}

	// Step 2: Persist Order in PostgreSQL
	if err := s.postgresRepo.CreateOrder(ctx, order); err != nil {
		s.logger.Error("Failed to persist order in PostgreSQL", zap.String("trace_id", traceID), zap.Error(err))
		return nil, fmt.Errorf("failed to store order: %w", err)
	}

	// Step 3: Emit OrderPlaced Event to Kafka
	event := &domain.OrderEvent{
		OrderID:     order.ID,
		UserID:      order.UserID,
		ProductID:   order.ProductID,
		Quantity:    order.Quantity,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
		TraceID:     traceID,
		Timestamp:   now,
	}
	if err := s.kafkaProducer.PublishOrderPlaced(ctx, event); err != nil {
		s.logger.Warn("Failed to emit Kafka event", zap.String("trace_id", traceID), zap.Error(err))
	}

	// Step 4: Enqueue Invoice Generation Task to RabbitMQ
	task := &domain.InvoiceTask{
		OrderID:     order.ID,
		UserID:      order.UserID,
		TotalAmount: order.TotalAmount,
		TraceID:     traceID,
		CreatedAt:   now,
	}
	if err := s.rabbitPublisher.PublishInvoiceTask(ctx, task); err != nil {
		s.logger.Warn("Failed to enqueue RabbitMQ task", zap.String("trace_id", traceID), zap.Error(err))
	}

	s.logger.Info("Successfully completed CreateOrder workflow",
		zap.String("order_id", order.ID),
		zap.String("trace_id", traceID),
		zap.Float64("total_amount", order.TotalAmount),
	)

	return order, nil
}

func (s *orderService) GetOrder(ctx context.Context, orderID string) (*domain.Order, error) {
	order, err := s.postgresRepo.GetOrder(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	return order, nil
}
