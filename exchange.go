package rerabbit

import "github.com/streadway/amqp"

type ExchangeOptions struct {
	Name       string
	Kind       string // "direct", "fanout", "topic", "headers"
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

// ExchangeDeclare create an exchange with parameters
func (r *rabbitBroker) ExchangeDeclare(opts ExchangeOptions) error {
	return r.channel.ExchangeDeclare(
		opts.Name,
		opts.Kind,
		opts.Durable,
		opts.AutoDelete,
		opts.Internal,
		opts.NoWait,
		opts.Args,
	)
}

// ExchangeDelete deletes an exchange
func (r *rabbitBroker) ExchangeDelete(name string, ifUnused, noWait bool) error {
	return r.channel.ExchangeDelete(name, ifUnused, noWait)
}
