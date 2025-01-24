package kafka

import (
	"component-master/config"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
)

type KafkaClientConfig struct {
	cfg    *config.KafkaConfig
	sarama sarama.Config
	logger *slog.Logger
}

type Admin struct {
	admin  sarama.ClusterAdmin
	client *KafkaClientConfig
}

func NewKafkaClient(cfg *config.KafkaConfig) (*KafkaClientConfig, error) {

	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Version = sarama.V2_8_0_0
	kafkaConfig.ClientID = cfg.ClientId

	// Producer configurations
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.Compression = sarama.CompressionSnappy
	kafkaConfig.Producer.Retry.Max = cfg.Producer.RetryMax
	kafkaConfig.Producer.Retry.Backoff = cfg.Producer.RetryBackoff
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Producer.Return.Errors = true
	kafkaConfig.Producer.MaxMessageBytes = cfg.Producer.MaxMessageBytes

	// Consumer configurations
	kafkaConfig.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	kafkaConfig.Consumer.Offsets.AutoCommit.Enable = true
	kafkaConfig.Consumer.Offsets.AutoCommit.Interval = time.Second
	kafkaConfig.Consumer.Group.Session.Timeout = time.Second * 30
	kafkaConfig.Consumer.MaxProcessingTime = cfg.Consumer.MaxProcessingTime
	kafkaConfig.Consumer.Fetch.Min = cfg.Consumer.FetchMin
	kafkaConfig.Consumer.Fetch.Max = cfg.Consumer.FetchMax

	// Network configurations
	kafkaConfig.Net.DialTimeout = cfg.DialTimeout
	kafkaConfig.Net.ReadTimeout = cfg.ReadTimeout
	kafkaConfig.Net.WriteTimeout = cfg.WriteTimeout
	kafkaConfig.Net.SASL.Enable = true
	kafkaConfig.Net.TLS.Enable = true

	// Metadata configurations
	kafkaConfig.Metadata.Retry.Max = cfg.MaxRetry
	kafkaConfig.Metadata.Retry.Backoff = cfg.RetryBackoff
	kafkaConfig.Metadata.RefreshFrequency = cfg.RefreshFrequency
	kafkaConfig.Metadata.Full = true

	return &KafkaClientConfig{
		cfg:    cfg,
		sarama: *kafkaConfig,
		logger: slog.Default(),
	}, nil
}

func (c *KafkaClientConfig) NewAdmin() (*Admin, error) {
	admin, err := sarama.NewClusterAdmin(c.cfg.Brokers, &c.sarama)
	if err != nil {
		return nil, fmt.Errorf("failed to create cluster admin: %w", err)
	}

	return &Admin{
		admin:  admin,
		client: c,
	}, nil
}

func (a *Admin) CreateTopic(ctx context.Context, topic string, numPartitions int32, replicationFactor int16) error {
	topicDetail := &sarama.TopicDetail{
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
		ConfigEntries: map[string]*string{
			"cleanup.policy":      stringPtr("delete"),
			"retention.ms":        stringPtr("604800000"), // 7 days
			"max.message.bytes":   stringPtr("1048576"),   // 1MB
			"min.insync.replicas": stringPtr("2"),
		},
	}

	err := a.admin.CreateTopic(topic, topicDetail, false)
	if err != nil {
		return fmt.Errorf("failed to create topic: %w", err)
	}

	return nil
}

func stringPtr(s string) *string {
	return &s
}
