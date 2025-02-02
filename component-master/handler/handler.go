package handler

import (
	"component-master/infra/kafka"
	"context"
	"log/slog"

	"github.com/IBM/sarama"
)

type MessageProcessor func(context.Context, []byte) error

type KafkaHandler interface {
	RegisterProcessor(topic string, processor MessageProcessor)
	Handle(ctx context.Context, msg *sarama.ConsumerMessage) error
}

type kafkaHandler struct {
	kafkaClient *kafka.KafkaClientConfig
	processors  map[string]MessageProcessor
}

func NewKafkaHandler(kafkaClient *kafka.KafkaClientConfig) KafkaHandler {
	return &kafkaHandler{kafkaClient: kafkaClient, processors: make(map[string]MessageProcessor)}
}

func (s *kafkaHandler) RegisterProcessor(topic string, processor MessageProcessor) {
	s.processors[topic] = processor
}

func (s *kafkaHandler) Handle(ctx context.Context, msg *sarama.ConsumerMessage) error {
	processor, ok := s.processors[msg.Topic]
	if !ok {
		return nil
	}
	if err := processor(ctx, msg.Value); err != nil {
		slog.Error("Failed to process message",
			"topic", msg.Topic,
			"partition", msg.Partition,
			"offset", msg.Offset,
			"error", err)
		return err
	}
	return nil
}
