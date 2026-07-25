package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"

	"github.com/orderpulse/internal/order/domain"
	"github.com/orderpulse/internal/worker/config"
)

type InvoiceConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	logger  *zap.Logger
	stopCh  chan struct{}
	wg      sync.WaitGroup
}

func NewInvoiceConsumer(cfg *config.Config, logger *zap.Logger) (*InvoiceConsumer, error) {
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("worker failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("worker failed to open RabbitMQ channel: %w", err)
	}

	// Ensure queue exists
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
		return nil, fmt.Errorf("worker failed to declare invoice_queue: %w", err)
	}

	return &InvoiceConsumer{
		conn:    conn,
		channel: ch,
		logger:  logger,
		stopCh:  make(chan struct{}),
	}, nil
}

func (c *InvoiceConsumer) Start(ctx context.Context) error {
	msgs, err := c.channel.Consume(
		"invoice_queue", // queue
		"worker-service", // consumer
		false,           // auto-ack (set false for manual ack)
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	if err != nil {
		return fmt.Errorf("worker failed to start consuming invoice_queue: %w", err)
	}

	c.logger.Info("Worker Service listening for tasks on RabbitMQ queue invoice_queue")

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			select {
			case <-c.stopCh:
				c.logger.Info("Worker consumer loop stopping...")
				return
			case msg, ok := <-msgs:
				if !ok {
					c.logger.Info("RabbitMQ delivery channel closed")
					return
				}
				c.processInvoiceTask(msg)
			}
		}
	}()

	return nil
}

func (c *InvoiceConsumer) processInvoiceTask(msg amqp.Delivery) {
	var task domain.InvoiceTask
	if err := json.Unmarshal(msg.Body, &task); err != nil {
		c.logger.Error("Failed to unmarshal InvoiceTask JSON", zap.Error(err))
		_ = msg.Nack(false, false)
		return
	}

	c.logger.Info("Started processing PDF invoice generation task",
		zap.String("order_id", task.OrderID),
		zap.String("trace_id", task.TraceID),
		zap.String("user_id", task.UserID),
		zap.Float64("total_amount", task.TotalAmount),
	)

	// Simulate PDF invoice generation & email sending
	time.Sleep(300 * time.Millisecond)

	c.logger.Info("Successfully generated PDF invoice and sent email",
		zap.String("order_id", task.OrderID),
		zap.String("trace_id", task.TraceID),
		zap.String("status", "INVOICE_SENT"),
	)

	if err := msg.Ack(false); err != nil {
		c.logger.Error("Failed to Acknowledge message to RabbitMQ", zap.Error(err))
	}
}

func (c *InvoiceConsumer) Stop(ctx context.Context) error {
	close(c.stopCh)
	c.wg.Wait()

	if err := c.channel.Close(); err != nil {
		c.logger.Warn("Error closing worker RabbitMQ channel", zap.Error(err))
	}
	if err := c.conn.Close(); err != nil {
		c.logger.Warn("Error closing worker RabbitMQ connection", zap.Error(err))
	}
	c.logger.Info("Worker Service shutdown complete")
	return nil
}
