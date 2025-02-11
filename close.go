package rerabbit

import "github.com/sirupsen/logrus"

func (b *rabbitBroker) Close(log *logrus.Logger) {
	if err := b.channel.Close(); err != nil {
		log.Errorf("Failed to close channel: %v", err)
	}

	if err := b.connection.Close(); err != nil {
		log.Errorf("Failed to close connection: %v", err)
	}
}

func (b *rabbitBroker) Cancel(consumerTag string) error {
	return b.channel.Cancel(consumerTag, false)
}
