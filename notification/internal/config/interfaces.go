package config

import "github.com/IBM/sarama"

type LoggerConfig interface {
	Level() string
	AsJson() bool
	EnableStdout() bool
	EnableOTLP() bool
	OTLPEndpoint() string
	ServiceName() string
	ServiceEnv() string
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
