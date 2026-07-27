package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"

	"github.com/orderpulse/backend/internal/order/config"
	"github.com/orderpulse/backend/internal/order/domain"
)

type rabbitMQPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	logger  *zap.Logger
}

var _ domain.TaskPublisher = (*rabbitMQPublisher)(nil)

func NewRabbitMQPublisher(cfg *config.Config, logger *zap.Logger) (domain.TaskPublisher, error) {
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ at %s: %w", cfg.RabbitMQURL, err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("failed to open RabbitMQ channel: %w", err)
	}

	// Declare exchange and queue
	err = ch.ExchangeDeclare(
		"order_exchange", // name
		"direct",         // type
		true,             // durable
		false,            // auto-deleted
		false,            // internal
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("failed to declare RabbitMQ exchange: %w", err)
	}

	_, err = ch.QueueDeclare(
		"invoice_queue", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("failed to declare RabbitMQ queue: %w", err)
	}

	err = ch.QueueBind(
		"invoice_queue",  // queue name
		"invoice_task",   // routing key
		"order_exchange", // exchange
		false,
		nil,
	)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("failed to bind RabbitMQ queue to exchange: %w", err)
	}

	logger.Info("RabbitMQ Publisher initialized for queue invoice_queue")

	return &rabbitMQPublisher{
		conn:    conn,
		channel: ch,
		logger:  logger,
	}, nil
}

func (p *rabbitMQPublisher) PublishInvoiceTask(ctx context.Context, task *domain.InvoiceTask) error {
	payload, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal InvoiceTask to JSON: %w", err)
	}

	err = p.channel.PublishWithContext(
		ctx,
		"order_exchange", // exchange
		"invoice_task",   // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         payload,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish RabbitMQ task: %w", err)
	}

	p.logger.Info("Enqueued InvoiceTask into RabbitMQ",
		zap.String("order_id", task.OrderID),
		zap.String("trace_id", task.TraceID),
		zap.String("queue", "invoice_queue"),
	)
	return nil
}

func (p *rabbitMQPublisher) Close() error {
	if err := p.channel.Close(); err != nil {
		p.logger.Warn("Error closing RabbitMQ channel", zap.Error(err))
	}
	if err := p.conn.Close(); err != nil {
		p.logger.Warn("Error closing RabbitMQ connection", zap.Error(err))
	}
	return nil
}
