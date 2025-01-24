package service

import "component-master/infra/kafka"

type KafkaService interface {
}

type kafkaService struct {
	kafkaClient kafka.KafkaClientConfig
}

func NewKafkaService(kafkaClient kafka.KafkaClientConfig) KafkaService {
	return &kafkaService{kafkaClient: kafkaClient}
}
