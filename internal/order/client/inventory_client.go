package client

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/orderpulse/internal/order/config"
	"github.com/orderpulse/internal/order/domain"
	inventoryv1 "github.com/orderpulse/proto/inventory/v1"
)

type inventoryClient struct {
	client inventoryv1.InventoryServiceClient
	logger *zap.Logger
}

var _ domain.InventoryClient = (*inventoryClient)(nil)

func NewInventoryClient(cfg *config.Config, logger *zap.Logger) (domain.InventoryClient, error) {
	conn, err := grpc.Dial(
		cfg.InventoryGRPCEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial inventory gRPC service at %s: %w", cfg.InventoryGRPCEndpoint, err)
	}

	cli := inventoryv1.NewInventoryServiceClient(conn)
	logger.Info("Connected to Inventory gRPC Service", zap.String("endpoint", cfg.InventoryGRPCEndpoint))

	return &inventoryClient{
		client: cli,
		logger: logger,
	}, nil
}

func (c *inventoryClient) ReserveStock(ctx context.Context, orderID, productID string, quantity int32) (bool, float64, string, error) {
	// First check stock price
	checkResp, err := c.client.CheckStock(ctx, &inventoryv1.CheckStockRequest{
		ProductId: productID,
		Quantity:  quantity,
	})
	if err != nil {
		return false, 0, "", fmt.Errorf("gRPC CheckStock call failed: %w", err)
	}

	if !checkResp.GetIsAvailable() {
		return false, checkResp.GetPrice(), "Product stock is not available", nil
	}

	unitPrice := checkResp.GetPrice()

	// Reserve stock
	reserveResp, err := c.client.ReserveStock(ctx, &inventoryv1.ReserveStockRequest{
		OrderId:   orderID,
		ProductId: productID,
		Quantity:  quantity,
	})
	if err != nil {
		return false, unitPrice, "", fmt.Errorf("gRPC ReserveStock call failed: %w", err)
	}

	return reserveResp.GetSuccess(), unitPrice, reserveResp.GetMessage(), nil
}
