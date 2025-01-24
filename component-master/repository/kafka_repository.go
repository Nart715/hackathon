package repository

import (
	"component-master/infra/kafka"
	"component-master/util"
	"context"
	"log/slog"

	"github.com/IBM/sarama"
)

type MessageHandlerImpl struct {
}

func NewMessageHandler() kafka.MessageHandler {
	return &MessageHandlerImpl{}
}

func (r *MessageHandlerImpl) Handle(ctx context.Context, message *sarama.ConsumerMessage) error {
	slog.Info("Message received", slog.String("Topic", message.Topic), slog.String("Partition", util.String("%d", message.Partition)), slog.String("Offset", util.String("%d", message.Offset)))
	slog.Info(string(message.Value))
	return nil
}
