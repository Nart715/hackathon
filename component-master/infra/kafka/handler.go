package kafka

import (
	"context"

	"github.com/IBM/sarama"
)

type MessageHandler interface {
	Handle(context.Context, *sarama.ConsumerMessage) error
}
