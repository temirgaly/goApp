package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/orderpulse/backend/internal/inventory/config"
	"github.com/orderpulse/backend/internal/inventory/domain"
)

type redisRepository struct {
	client *redis.Client
	logger *zap.Logger
}

var _ domain.RedisRepository = (*redisRepository)(nil)

func NewRedisRepository(cfg *config.Config, logger *zap.Logger) (domain.RedisRepository, error) {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return &redisRepository{
		client: client,
		logger: logger,
	}, nil
}

func (r *redisRepository) GetStock(ctx context.Context, productID string) (*domain.StockStatus, error) {
	key := fmt.Sprintf("stock:%s", productID)
	vals, err := r.client.HMGet(ctx, key, "available_quantity", "price").Result()
	if err != nil {
		return nil, fmt.Errorf("redis HMGet error: %w", err)
	}

	if vals[0] == nil || vals[1] == nil {
		return nil, domain.ErrItemNotFound
	}

	qtyStr, ok1 := vals[0].(string)
	priceStr, ok2 := vals[1].(string)
	if !ok1 || !ok2 {
		return nil, domain.ErrItemNotFound
	}

	qty, err := strconv.ParseInt(qtyStr, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid quantity in redis: %w", err)
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid price in redis: %w", err)
	}

	return &domain.StockStatus{
		IsAvailable:       qty > 0,
		AvailableQuantity: int32(qty),
		Price:             price,
	}, nil
}

func (r *redisRepository) SetStock(ctx context.Context, productID string, status *domain.StockStatus) error {
	key := fmt.Sprintf("stock:%s", productID)
	pipe := r.client.Pipeline()
	pipe.HSet(ctx, key, map[string]interface{}{
		"available_quantity": status.AvailableQuantity,
		"price":              status.Price,
	})
	pipe.Expire(ctx, key, 10*time.Minute)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to set stock cache in redis: %w", err)
	}
	return nil
}

func (r *redisRepository) Invalidate(ctx context.Context, productID string) error {
	key := fmt.Sprintf("stock:%s", productID)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to invalidate redis cache key %s: %w", key, err)
	}
	return nil
}
