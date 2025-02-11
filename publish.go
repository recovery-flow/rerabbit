package rerabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

type ConsumeOptions struct {
	QueueName   string
	ConsumerTag string
	AutoAck     bool
	Exclusive   bool
	NoLocal     bool
	NoWait      bool
	Args        amqp.Table
}

type PublishOptions struct {
	Exchange   string
	RoutingKey string
	Mandatory  bool
	Immediate  bool

	ContentType   string
	DeliveryMode  uint8 // 1=Transient, 2=Persistent
	Headers       amqp.Table
	CorrelationID string
	ReplyTo       string

	Body []byte
}

// Publish publishes a message to the exchange with the specified routing key.
func (r *rabbitBroker) Publish(ctx context.Context, opts PublishOptions) error {
	if opts.ContentType == "" {
		opts.ContentType = "application/json"
	}
	pub := amqp.Publishing{
		ContentType:   opts.ContentType,
		Body:          opts.Body,
		DeliveryMode:  opts.DeliveryMode,
		Headers:       opts.Headers,
		CorrelationId: opts.CorrelationID,
		ReplyTo:       opts.ReplyTo,
	}

	resultCh := make(chan error, 1)

	go func() {
		err := r.channel.Publish(
			opts.Exchange,
			opts.RoutingKey,
			opts.Mandatory,
			opts.Immediate,
			pub,
		)
		resultCh <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-resultCh:
		return err
	}
}

// PublishJSON json-encodes the data and publishes it to the exchange with the specified routing key.
func (r *rabbitBroker) PublishJSON(ctx context.Context, data interface{}, opts PublishOptions) error {
	// маршалим в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	opts.Body = jsonData
	if opts.ContentType == "" {
		opts.ContentType = "application/json"
	}
	return r.Publish(ctx, opts)
}

// Consume consumes messages from the queue.
func (r *rabbitBroker) Consume(ctx context.Context, opts ConsumeOptions, handler func(context.Context, amqp.Delivery)) error {
	msgs, err := r.channel.Consume(
		opts.QueueName,
		opts.ConsumerTag,
		opts.AutoAck,
		opts.Exclusive,
		opts.NoLocal,
		opts.NoWait,
		opts.Args,
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				// Останавливаем consumer
				_ = r.channel.Cancel(opts.ConsumerTag, false)
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}
				handler(ctx, msg)
			}
		}
	}()

	return nil
}

// ConsumeLimited consumes messages from the queue, but stops after maxMessages.
func (r *rabbitBroker) ConsumeLimited(ctx context.Context, opts ConsumeOptions, maxMessages int, handler func(context.Context, amqp.Delivery)) error {
	msgs, err := r.channel.Consume(
		opts.QueueName, opts.ConsumerTag, opts.AutoAck, opts.Exclusive, opts.NoLocal, opts.NoWait, opts.Args,
	)
	if err != nil {
		return err
	}

	for i := 0; i < maxMessages; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("channel closed")
			}
			handler(ctx, msg)
		}
	}

	return nil
}

// PublishWithRetry publishes a message to the exchange with the specified routing key.
func (r *rabbitBroker) PublishWithRetry(ctx context.Context, opts PublishOptions, maxRetries int, retryDelay time.Duration) error {
	var lastErr error

	for i := 0; i <= maxRetries; i++ {
		err := r.Publish(ctx, opts)
		if err == nil {
			return nil
		}

		lastErr = err
		if i < maxRetries {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(retryDelay):
			}
		}
	}

	return fmt.Errorf("publish failed after %d retries: %w", maxRetries, lastErr)
}

// PublishWithConfirm publishes a message to the exchange with the specified routing key and waits for a confirmation.
func (r *rabbitBroker) PublishWithConfirm(ctx context.Context, opts PublishOptions, timeout time.Duration) error {
	if err := r.channel.Confirm(false); err != nil {
		return fmt.Errorf("failed to put channel in confirm mode: %w", err)
	}

	confirmCh := r.channel.NotifyPublish(make(chan amqp.Confirmation, 1))

	err := r.Publish(ctx, opts)
	if err != nil {
		return err
	}

	select {
	case confirm := <-confirmCh:
		if confirm.Ack {
			return nil // Сообщение подтверждено
		}
		return fmt.Errorf("message was not confirmed")
	case <-time.After(timeout):
		return fmt.Errorf("timeout waiting for publish confirmation")
	}
}
