package rerabbit

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// SetQos Configures preloading (how many messages can be “in flight”).
func (r *rabbitBroker) SetQos(prefetchCount, prefetchSize int, global bool) error {
	return r.channel.Qos(prefetchCount, prefetchSize, global)
}

// Cancel Stops consuming messages.
func (r *rabbitBroker) Cancel(consumerTag string) error {
	return r.channel.Cancel(consumerTag, false)
}

// Close Closes the connection and channel.
func (r *rabbitBroker) Close(log *logrus.Logger) {
	if err := r.channel.Close(); err != nil {
		log.Errorf("Failed to close channel: %v", err)
	}
	if err := r.conn.Close(); err != nil {
		log.Errorf("Failed to close connection: %v", err)
	}
}

// Reconnect Reconnects to the broker.
func (r *rabbitBroker) Reconnect(url string) error {
	for {
		conn, err := amqp.Dial(url)
		if err == nil {
			r.conn = conn
			r.channel, err = conn.Channel()
			if err == nil {
				return nil
			}
		}

		time.Sleep(5 * time.Second) // Подождать перед ретраем
	}
}

// GracefulShutdown Gracefully shuts down the broker.
func (r *rabbitBroker) GracefulShutdown(log *logrus.Logger) {
	log.Info("Shutting down RabbitMQ broker...")

	if err := r.channel.Cancel("", false); err != nil {
		log.Errorf("Failed to cancel consumers: %v", err)
	}

	if err := r.channel.Close(); err != nil {
		log.Errorf("Failed to close channel: %v", err)
	}

	if err := r.conn.Close(); err != nil {
		log.Errorf("Failed to close connection: %v", err)
	}

	log.Info("RabbitMQ broker shutdown complete.")
}
