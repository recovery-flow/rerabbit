package cifra_rabbit

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"
)

type PublishOptions struct {
	QueueName      string
	RoutingKey     string
	DeclareQueue   bool
	QueueDurable   bool
	QueueAutoDel   bool
	QueueExclusive bool
	QueueArgs      amqp.Table

	// Параметры самого Publish
	Exchange    string
	Mandatory   bool
	Immediate   bool
	ContentType string
}

// Publish отправляет сообщение в exchange с routingKey.
// Опционально может объявить и привязать очередь перед публикацией.
func (b *rabbitBroker) Publish(ctx context.Context, body []byte, opts PublishOptions) error {
	// Если надо объявить/привязать очередь — делаем
	if opts.DeclareQueue {
		if opts.QueueName == "" {
			return fmt.Errorf("queue name must not be empty when DeclareQueue is true")
		}

		err := b.DeclareAndBindQueue(
			opts.QueueName,
			opts.RoutingKey,
			opts.QueueDurable,
			opts.QueueAutoDel,
			opts.QueueExclusive,
			opts.QueueArgs,
		)
		if err != nil {
			return fmt.Errorf("failed to declare/bind queue: %w", err)
		}
	}

	// Укажем exchange по умолчанию, если не задан в opts
	exchangeName := b.exchange
	if opts.Exchange != "" {
		exchangeName = opts.Exchange
	}

	// Публикуем
	pub := amqp.Publishing{
		ContentType: opts.ContentType,
		Body:        body,
	}
	if pub.ContentType == "" {
		pub.ContentType = "application/json"
	}

	err := b.channel.Publish(
		exchangeName,
		opts.RoutingKey,
		opts.Mandatory,
		opts.Immediate,
		pub,
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (b *rabbitBroker) SetQos(prefetchCount, prefetchSize int, global bool) error {
	return b.channel.Qos(prefetchCount, prefetchSize, global)
}
