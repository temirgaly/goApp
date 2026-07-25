package grpc

import (
	"context"

	"go.uber.org/zap"

	"github.com/orderpulse/internal/inventory/domain"
	inventoryv1 "github.com/orderpulse/proto/inventory/v1"
)

type Handler struct {
	inventoryv1.UnimplementedInventoryServiceServer
	svc    domain.Service
	logger *zap.Logger
}

func NewHandler(svc domain.Service, logger *zap.Logger) *Handler {
	return &Handler{
		svc:    svc,
		logger: logger,
	}
}

func (h *Handler) CheckStock(ctx context.Context, req *inventoryv1.CheckStockRequest) (*inventoryv1.CheckStockResponse, error) {
	h.logger.Info("Received CheckStock gRPC request",
		zap.String("product_id", req.GetProductId()),
		zap.Int32("quantity", req.GetQuantity()),
	)

	status, err := h.svc.CheckStock(ctx, req.GetProductId(), req.GetQuantity())
	if err != nil {
		h.logger.Error("Failed to check stock", zap.Error(err))
		return nil, err
	}

	return &inventoryv1.CheckStockResponse{
		IsAvailable:       status.IsAvailable,
		AvailableQuantity: status.AvailableQuantity,
		Price:             status.Price,
	}, nil
}

func (h *Handler) ReserveStock(ctx context.Context, req *inventoryv1.ReserveStockRequest) (*inventoryv1.ReserveStockResponse, error) {
	h.logger.Info("Received ReserveStock gRPC request",
		zap.String("order_id", req.GetOrderId()),
		zap.String("product_id", req.GetProductId()),
		zap.Int32("quantity", req.GetQuantity()),
	)

	success, msg, err := h.svc.ReserveStock(ctx, req.GetOrderId(), req.GetProductId(), req.GetQuantity())
	if err != nil {
		h.logger.Error("Failed to reserve stock", zap.Error(err))
		return nil, err
	}

	return &inventoryv1.ReserveStockResponse{
		Success: success,
		Message: msg,
	}, nil
}
