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

type OrderAssembledConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}

type OrderPaidConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}

type TelegramBotConfig interface {
	Token() string
	ChatID() int64
}
