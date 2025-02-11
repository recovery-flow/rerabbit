package rerabbit

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// RabbitBroker is an interface for interacting with RabbitMQ
type RabbitBroker interface {
	// -- Exchange management --
	ExchangeDeclare(opts ExchangeOptions) error
	ExchangeDelete(name string, ifUnused, noWait bool) error

	// -- Queue management --
	QueueDeclare(opts QueueOptions) (amqp.Queue, error)
	QueueDelete(name string, ifUnused, ifEmpty, noWait bool) (int, error)

	// -- Bind / Unbind --
	QueueBind(opts BindOptions) error
	QueueUnbind(opts UnbindOptions) error
	DeclareQueueWithDLX(opts QueueOptions, dlxName string) (amqp.Queue, error)

	// -- Publish / Consume --
	Publish(ctx context.Context, opts PublishOptions) error
	PublishJSON(ctx context.Context, data interface{}, opts PublishOptions) error
	PublishWithRetry(ctx context.Context, opts PublishOptions, maxRetries int, retryDelay time.Duration) error
	PublishWithConfirm(ctx context.Context, opts PublishOptions, timeout time.Duration) error
	Consume(ctx context.Context, opts ConsumeOptions, handler func(context.Context, amqp.Delivery)) error
	ConsumeLimited(ctx context.Context, opts ConsumeOptions, maxMessages int, handler func(context.Context, amqp.Delivery)) error

	// -- Other --
	SetQos(prefetchCount, prefetchSize int, global bool) error
	Cancel(consumerTag string) error
	Close(log *logrus.Logger)
	Reconnect(url string) error
	GracefulShutdown(log *logrus.Logger)
}

// rabbitBroker is an implementation of RabbitBroker
type rabbitBroker struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewBroker creates a connection to RabbitMQ, but does not declare anything
func NewBroker(amqpURL string) (RabbitBroker, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	return &rabbitBroker{
		conn:    conn,
		channel: ch,
	}, nil
}
