package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	"github.com/orderpulse/backend/internal/order/config"
	"github.com/orderpulse/backend/internal/order/domain"
)

type kafkaProducer struct {
	writer *kafka.Writer
	logger *zap.Logger
}

var _ domain.EventProducer = (*kafkaProducer)(nil)

func NewKafkaProducer(cfg *config.Config, logger *zap.Logger) (domain.EventProducer, error) {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.KafkaBrokers...),
		Topic:    "order-events",
		Balancer: &kafka.LeastBytes{},
	}

	logger.Info("Kafka Producer initialized for topic order-events", zap.Strings("brokers", cfg.KafkaBrokers))

	return &kafkaProducer{
		writer: writer,
		logger: logger,
	}, nil
}

func (p *kafkaProducer) PublishOrderPlaced(ctx context.Context, event *domain.OrderEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal OrderEvent to JSON: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(event.OrderID),
		Value: payload,
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to write message to Kafka: %w", err)
	}

	p.logger.Info("Published OrderPlaced event to Kafka",
		zap.String("order_id", event.OrderID),
		zap.String("trace_id", event.TraceID),
		zap.String("topic", "order-events"),
	)
	return nil
}

func (p *kafkaProducer) Close() error {
	if err := p.writer.Close(); err != nil {
		return fmt.Errorf("failed to close Kafka writer: %w", err)
	}
	return nil
}
