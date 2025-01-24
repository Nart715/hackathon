package kafka

import (
	"component-master/util"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

type MessageHandler interface {
	Handle(context.Context, *sarama.ConsumerMessage) error
}

type Consumer struct {
	group   sarama.ConsumerGroup
	client  *KafkaClientConfig
	handler MessageHandler
	topics  []string
	wg      sync.WaitGroup
}

func (c *KafkaClientConfig) NewConsumer(ctx context.Context, groupID string, topics []string, handler MessageHandler) (*Consumer, error) {
	if groupID == "" {
		return nil, fmt.Errorf("consumer group ID is required")
	}
	if len(topics) == 0 {
		return nil, fmt.Errorf("at least one topic is required")
	}
	if handler == nil {
		return nil, fmt.Errorf("message handler is required")
	}

	group, err := sarama.NewConsumerGroup(c.cfg.Brokers, groupID, &c.sarama)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	consumer := &Consumer{
		group:   group,
		client:  c,
		handler: handler,
		topics:  topics,
	}

	return consumer, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			if err := c.group.Consume(ctx, c.topics, c); err != nil {
				c.client.logger.Error(util.String("Error from consumer: %s", err.Error()))

				select {
				case <-ctx.Done():
					return
				case <-time.After(c.client.cfg.Consumer.RetryBackoff):
					continue
				}
			}

			if ctx.Err() != nil {
				return
			}
		}
	}()

	return nil
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		start := time.Now()
		err := c.handler.Handle(session.Context(), message)
		duration := time.Since(start)

		c.client.logger.Info(util.String("Topic: %s", message.Topic))
		c.client.logger.Info(util.String("Partition: %d", message.Partition))
		c.client.logger.Info(util.String("Offset: %d", message.Offset))
		c.client.logger.Info(util.String("Processing time: %v", duration))
		c.client.logger.Info(util.String("Error: %v", err))
		if err != nil {
			c.client.logger.Error("Failed to process message")
			continue
		}
		session.MarkMessage(message, "")
	}
	return nil
}

func (c *Consumer) Close() error {
	err := c.group.Close()
	c.wg.Wait()
	return err
}
