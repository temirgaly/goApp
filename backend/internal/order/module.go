package order

import (
	"context"
	"net"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/orderpulse/backend/internal/order/client"
	"github.com/orderpulse/backend/internal/order/config"
	ordergrpc "github.com/orderpulse/backend/internal/order/grpc"
	"github.com/orderpulse/backend/internal/order/kafka"
	"github.com/orderpulse/backend/internal/order/rabbitmq"
	"github.com/orderpulse/backend/internal/order/repository"
	ordersvc "github.com/orderpulse/backend/internal/order/service"
	"github.com/orderpulse/backend/pkg/logger"
	orderv1 "github.com/orderpulse/proto/order/v1"
)

var Module = fx.Options(
	fx.Provide(
		config.Load,
		ProvideLogger,
		repository.NewPostgresRepository,
		client.NewInventoryClient,
		kafka.NewKafkaProducer,
		rabbitmq.NewRabbitMQPublisher,
		ordersvc.NewOrderService,
		ordergrpc.NewHandler,
	),
	fx.Invoke(RegisterGRPCServer),
)

func ProvideLogger() (*zap.Logger, error) {
	log, _, err := logger.NewLogger("order-service")
	return log, err
}

func RegisterGRPCServer(
	lc fx.Lifecycle,
	logger *zap.Logger,
	cfg *config.Config,
	handler *ordergrpc.Handler,
) {
	server := grpc.NewServer()
	orderv1.RegisterOrderServiceServer(server, handler)

	// Enable gRPC Reflection for Postman / grpcurl auto-discovery
	reflection.Register(server)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen("tcp", cfg.GRPCPort)
			if err != nil {
				logger.Error("Failed to bind Order gRPC port", zap.String("port", cfg.GRPCPort), zap.Error(err))
				return err
			}
			logger.Info("Starting Order gRPC Server", zap.String("port", cfg.GRPCPort))
			go func() {
				if err := server.Serve(lis); err != nil {
					logger.Error("Order gRPC server error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping Order gRPC Server gracefully")
			server.GracefulStop()
			return nil
		},
	})
}
