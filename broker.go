package cifra_rabbit

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type RabbitBroker interface {
	DeclareAndBindQueue(queueName, routingKey string, durable, autoDelete, exclusive bool, args amqp.Table) error
	ConsumeQueue(ctx context.Context, queueName, consumerTag string, handler func(context.Context, []byte) error) error

	Publish(ctx context.Context, body []byte, opts PublishOptions) error

	Close(log *logrus.Logger)
	Cancel(consumerTag string) error
}

type rabbitBroker struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	exchange   string
}

func NewBroker(
	url string,
	exchangeName string, // name of the exchange
	kind string, // type of the exchange (direct, topic, headers)
	durable bool, // durable exchange survives exchange restart
	autoDel bool, // auto-deleted exchange is deleted when last queue is unbound
	internal bool, // internal exchange is not listed in exchange.declare-ok
	noWait bool, // no-wait exchange.declare-ok
) (RabbitBroker, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		err = conn.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	err = ch.ExchangeDeclare(
		exchangeName,
		kind,
		durable,
		autoDel,
		internal,
		noWait,
		nil,
	)
	if err != nil {
		err = ch.Close()
		if err != nil {
			return nil, err
		}
		err = conn.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	return &rabbitBroker{
		connection: conn,
		channel:    ch,
		exchange:   exchangeName,
	}, nil
}
