package inventory

import (
	"context"
	"net"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/orderpulse/backend/internal/inventory/config"
	inventorygrpc "github.com/orderpulse/backend/internal/inventory/grpc"
	"github.com/orderpulse/backend/internal/inventory/repository"
	inventorysvc "github.com/orderpulse/backend/internal/inventory/service"
	"github.com/orderpulse/backend/pkg/logger"
	inventoryv1 "github.com/orderpulse/proto/inventory/v1"
)

var Module = fx.Options(
	fx.Provide(
		config.Load,
		ProvideLogger,
		repository.NewPostgresRepository,
		repository.NewRedisRepository,
		inventorysvc.NewInventoryService,
		inventorygrpc.NewHandler,
	),
	fx.Invoke(RegisterGRPCServer),
)

func ProvideLogger() (*zap.Logger, error) {
	log, _, err := logger.NewLogger("inventory-service")
	return log, err
}

func RegisterGRPCServer(
	lc fx.Lifecycle,
	logger *zap.Logger,
	cfg *config.Config,
	handler *inventorygrpc.Handler,
) {
	server := grpc.NewServer()
	inventoryv1.RegisterInventoryServiceServer(server, handler)

	// Enable gRPC Server Reflection for Postman / grpcurl auto-discovery
	reflection.Register(server)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen("tcp", cfg.GRPCPort)
			if err != nil {
				logger.Error("Failed to bind gRPC port", zap.String("port", cfg.GRPCPort), zap.Error(err))
				return err
			}
			logger.Info("Starting Inventory gRPC Server", zap.String("port", cfg.GRPCPort))
			go func() {
				if err := server.Serve(lis); err != nil {
					logger.Error("gRPC server encountered an error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping Inventory gRPC Server gracefully")
			server.GracefulStop()
			return nil
		},
	})
}
