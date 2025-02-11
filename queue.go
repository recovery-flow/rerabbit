package rerabbit

import (
	"github.com/streadway/amqp"
)

type QueueOptions struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

type BindOptions struct {
	QueueName  string
	RoutingKey string
	Exchange   string
	NoWait     bool
	Args       amqp.Table
}

type UnbindOptions struct {
	QueueName  string
	RoutingKey string
	Exchange   string
	Args       amqp.Table
}

// QueueDeclare creates a queue with parameters
func (r *rabbitBroker) QueueDeclare(opts QueueOptions) (amqp.Queue, error) {
	return r.channel.QueueDeclare(
		opts.Name,
		opts.Durable,
		opts.AutoDelete,
		opts.Exclusive,
		opts.NoWait,
		opts.Args,
	)
}

// QueueDelete deletes a queue
func (r *rabbitBroker) QueueDelete(name string, ifUnused, ifEmpty, noWait bool) (int, error) {
	return r.channel.QueueDelete(name, ifUnused, ifEmpty, noWait)
}

// QueueBind binds a queue to an exchange
func (r *rabbitBroker) QueueBind(opts BindOptions) error {
	return r.channel.QueueBind(
		opts.QueueName,
		opts.RoutingKey,
		opts.Exchange,
		opts.NoWait,
		opts.Args,
	)
}

// QueueUnbind unbinds a queue from an exchange
func (r *rabbitBroker) QueueUnbind(opts UnbindOptions) error {
	return r.channel.QueueUnbind(
		opts.QueueName,
		opts.RoutingKey,
		opts.Exchange,
		opts.Args,
	)
}
