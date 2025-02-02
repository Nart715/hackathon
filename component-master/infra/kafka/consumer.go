package kafka

import (
	"context"
	"log/slog"
	"sync"

	"github.com/IBM/sarama"
)

type ConsumerHandler struct {
	ready          chan bool
	messageHandler func(context.Context, *sarama.ConsumerMessage) error
	logger         *slog.Logger
	wg             *sync.WaitGroup
}

func (h *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	close(h.ready)
	return nil
}

func (h *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		h.logger.Info("Received message",
			"topic", message.Topic,
			"partition", message.Partition,
			"offset", message.Offset,
			"key", string(message.Key),
			"value", string(message.Value))

		if err := h.messageHandler(context.Background(), message); err != nil {
			h.logger.Error("Failed to process message",
				"topic", message.Topic,
				"partition", message.Partition,
				"offset", message.Offset,
				"error", err)
			continue
		}

		session.MarkMessage(message, "")
	}
	return nil
}

func (k *KafkaClientConfig) Consume(ctx context.Context, topics []string, handler func(context.Context, *sarama.ConsumerMessage) error) error {
	k.wg.Add(1)

	consumerHandler := &ConsumerHandler{
		ready:          make(chan bool),
		messageHandler: handler,
		logger:         k.logger,
		wg:             &k.wg,
	}

	// Start error handler
	go func() {
		for err := range k.consumer.Errors() {
			k.logger.Error("Consumer group error", "error", err)
		}
	}()

	go func() {
		defer k.wg.Done()
		for {
			k.logger.Info("Starting consumer group",
				"topics", topics,
				"group", k.cfg.Kafka.Consumer.GroupId)

			err := k.consumer.Consume(ctx, topics, consumerHandler)
			if err != nil {
				k.logger.Error("Consumer group error", "error", err)
			}

			if ctx.Err() != nil {
				k.logger.Info("Context cancelled, stopping consumer")
				return
			}

			consumerHandler.ready = make(chan bool)
		}
	}()

	<-consumerHandler.ready
	k.logger.Info("Consumer started successfully",
		"topics", topics,
		"group", k.cfg.Kafka.Consumer.GroupId)

	return nil
}
