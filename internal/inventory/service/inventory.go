package service

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/orderpulse/internal/inventory/domain"
)

type inventoryService struct {
	postgresRepo domain.PostgresRepository
	redisRepo    domain.RedisRepository
	logger       *zap.Logger
}

var _ domain.Service = (*inventoryService)(nil)

func NewInventoryService(
	pgRepo domain.PostgresRepository,
	rRepo domain.RedisRepository,
	logger *zap.Logger,
) domain.Service {
	return &inventoryService{
		postgresRepo: pgRepo,
		redisRepo:    rRepo,
		logger:       logger,
	}
}

func (s *inventoryService) CheckStock(ctx context.Context, productID string, quantity int32) (*domain.StockStatus, error) {
	status, err := s.redisRepo.GetStock(ctx, productID)
	if err == nil {
		s.logger.Debug("Stock status served from Redis cache", zap.String("product_id", productID))
		status.IsAvailable = status.AvailableQuantity >= quantity
		return status, nil
	}

	s.logger.Info("Redis cache miss, querying PostgreSQL database", zap.String("product_id", productID))
	item, err := s.postgresRepo.GetItem(ctx, productID)
	if err != nil {
		if errors.Is(err, domain.ErrItemNotFound) {
			return &domain.StockStatus{IsAvailable: false, AvailableQuantity: 0, Price: 0}, nil
		}
		return nil, fmt.Errorf("service failed to fetch item from DB: %w", err)
	}

	stockStatus := &domain.StockStatus{
		IsAvailable:       item.AvailableQuantity >= quantity,
		AvailableQuantity: item.AvailableQuantity,
		Price:             item.Price,
	}

	if err := s.redisRepo.SetStock(ctx, productID, stockStatus); err != nil {
		s.logger.Warn("Failed to cache stock status in Redis", zap.Error(err))
	}

	return stockStatus, nil
}

func (s *inventoryService) ReserveStock(ctx context.Context, orderID, productID string, quantity int32) (bool, string, error) {
	item, err := s.postgresRepo.ReserveStock(ctx, orderID, productID, quantity)
	if err != nil {
		if errors.Is(err, domain.ErrItemNotFound) {
			return false, "Product not found", nil
		}
		if errors.Is(err, domain.ErrInsufficientStock) {
			return false, "Insufficient stock available", nil
		}
		return false, "", fmt.Errorf("failed to reserve stock in DB: %w", err)
	}

	s.logger.Info("Successfully reserved stock in PostgreSQL",
		zap.String("order_id", orderID),
		zap.String("product_id", productID),
		zap.Int32("remaining_stock", item.AvailableQuantity),
	)

	updatedStatus := &domain.StockStatus{
		IsAvailable:       item.AvailableQuantity > 0,
		AvailableQuantity: item.AvailableQuantity,
		Price:             item.Price,
	}

	if err := s.redisRepo.SetStock(ctx, productID, updatedStatus); err != nil {
		s.logger.Warn("Failed to update Redis cache after stock reservation", zap.Error(err))
	}

	return true, "Stock successfully reserved", nil
}
