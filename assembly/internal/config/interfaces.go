package config

import (
	"time"

	"github.com/IBM/sarama"
)

type OtelCollectorConfig interface {
	CollectorEndpoint() string
	ServiceName() string
	ServiceEnv() string
	CollectorInterval() time.Duration
}

type LoggerConfig interface {
	Level() string
	AsJson() bool
	EnableStdout() bool
	EnableOTLP() bool
}

type KafkaConfig interface {
	Brokers() []string
}

type OrderPaidConsumerConfig interface {
	TopicName() string
	ConsumerGroupId() string
	Config() *sarama.Config
}

type OrderAssembleProducerConfig interface {
	TopicName() string
	Config() *sarama.Config
}
