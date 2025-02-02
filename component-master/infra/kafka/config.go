package kafka

import (
	mconfig "component-master/config"
	"component-master/util"
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

type KafkaClientConfig struct {
	cfg      *mconfig.Config
	sarama   sarama.Config
	producer sarama.SyncProducer
	consumer sarama.ConsumerGroup
	logger   *slog.Logger
	wg       sync.WaitGroup
}

func (r *KafkaClientConfig) GetConfig() *mconfig.Config {
	return r.cfg
}

func NewKafkaClient(cfg *mconfig.Config) (*KafkaClientConfig, error) {
	config := sarama.NewConfig()

	config.Version = sarama.V2_8_0_0 // Use appropriate version for your Kafka cluster

	config.Producer.RequiredAcks = 1
	config.Producer.Retry.Max = cfg.Kafka.MaxRetry
	config.Producer.Retry.Backoff = cfg.Kafka.RetryBackoff
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Group.Session.Timeout = time.Second * 20
	config.Consumer.Group.Heartbeat.Interval = time.Second * 6
	config.Consumer.Group.Rebalance.Timeout = time.Second * 60
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true

	// Client settings
	config.Net.MaxOpenRequests = 5
	config.Net.KeepAlive = time.Second * 30
	config.Net.ReadTimeout = cfg.Kafka.ReadTimeout
	config.Net.WriteTimeout = cfg.Kafka.WriteTimeout
	config.Net.SASL.Handshake = true

	// Debug
	config.ClientID = "your-client-id" // Set a meaningful client ID

	if cfg.Kafka.EnableTLS {
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	if cfg.Kafka.EnableSASL {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = "admin"
		config.Net.SASL.Password = "admin"
		config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	}

	// Validate the config
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// Create producer
	producer, err := sarama.NewSyncProducer(cfg.Kafka.Brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	// Create consumer group
	consumer, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers, cfg.Kafka.Consumer.GroupId, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return &KafkaClientConfig{
		producer: producer,
		consumer: consumer,
		cfg:      cfg,
		sarama:   *config,
		logger:   slog.Default(),
	}, nil
}

func (k *KafkaClientConfig) Produce(ctx context.Context, topic string, key, value []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	partition, offset, err := k.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	k.logger.Info(util.String("Message sent to partition %d at offset %d\n", partition, offset))
	return nil
}

func (k *KafkaClientConfig) Close() error {
	k.wg.Wait()
	return k.consumer.Close()
}
