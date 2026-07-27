package worker

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/orderpulse/backend/internal/worker/config"
	"github.com/orderpulse/backend/internal/worker/consumer"
	"github.com/orderpulse/backend/pkg/logger"
)

var Module = fx.Options(
	fx.Provide(
		config.Load,
		ProvideLogger,
		consumer.NewInvoiceConsumer,
	),
	fx.Invoke(RegisterWorkerConsumer),
)

func ProvideLogger() (*zap.Logger, error) {
	log, _, err := logger.NewLogger("worker-service")
	return log, err
}

func RegisterWorkerConsumer(
	lc fx.Lifecycle,
	logger *zap.Logger,
	cons *consumer.InvoiceConsumer,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting Worker Service Consumer")
			return cons.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping Worker Service Consumer")
			return cons.Stop(ctx)
		},
	})
}
