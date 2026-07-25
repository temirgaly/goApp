package grpc

import (
	"context"

	"go.uber.org/zap"

	"github.com/orderpulse/internal/order/domain"
	orderv1 "github.com/orderpulse/proto/order/v1"
)

type Handler struct {
	orderv1.UnimplementedOrderServiceServer
	svc    domain.Service
	logger *zap.Logger
}

func NewHandler(svc domain.Service, logger *zap.Logger) *Handler {
	return &Handler{
		svc:    svc,
		logger: logger,
	}
}

func (h *Handler) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (*orderv1.CreateOrderResponse, error) {
	h.logger.Info("Received CreateOrder gRPC request",
		zap.String("user_id", req.GetUserId()),
		zap.String("product_id", req.GetProductId()),
		zap.Int32("quantity", req.GetQuantity()),
	)

	order, err := h.svc.CreateOrder(ctx, req.GetUserId(), req.GetProductId(), req.GetQuantity())
	if err != nil {
		h.logger.Error("Failed to process CreateOrder", zap.Error(err))
		return nil, err
	}

	return &orderv1.CreateOrderResponse{
		OrderId:     order.ID,
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
	}, nil
}

func (h *Handler) GetOrderStatus(ctx context.Context, req *orderv1.GetOrderStatusRequest) (*orderv1.GetOrderStatusResponse, error) {
	h.logger.Info("Received GetOrderStatus gRPC request", zap.String("order_id", req.GetOrderId()))

	order, err := h.svc.GetOrder(ctx, req.GetOrderId())
	if err != nil {
		h.logger.Error("Failed to fetch order status", zap.Error(err))
		return nil, err
	}

	return &orderv1.GetOrderStatusResponse{
		OrderId:   order.ID,
		Status:    order.Status,
		CreatedAt: order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}
