package cifra_rabbit

import (
	"context"

	"github.com/streadway/amqp"
)

func (b *rabbitBroker) DeclareAndBindQueue(queueName, routingKey string, durable, autoDelete, exclusive bool, args amqp.Table) error {
	// 1. Declare queue
	_, err := b.channel.QueueDeclare(
		queueName,
		durable,
		autoDelete,
		exclusive,
		false,
		args,
	)
	if err != nil {
		return err
	}

	// 2. Bind queue
	err = b.channel.QueueBind(
		queueName,
		routingKey,
		b.exchange,
		false,
		nil,
	)
	return err
}

func (b *rabbitBroker) ConsumeQueue(ctx context.Context, queueName, consumerTag string, handler func(context.Context, []byte) error) error {
	msgs, err := b.channel.Consume(
		queueName,
		consumerTag, false, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	go b.handleMessages(ctx, msgs, handler)
	return nil
}

func (b *rabbitBroker) handleMessages(ctx context.Context, msgs <-chan amqp.Delivery, handler func(context.Context, []byte) error) {
	for msg := range msgs {
		if err := handler(ctx, msg.Body); err != nil {
			_ = msg.Nack(false, true) // requeue = true
		} else {
			_ = msg.Ack(false)
		}
	}
}
