package rerabbit

import "github.com/sirupsen/logrus"

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
