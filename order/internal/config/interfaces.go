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

type OrderHTTPConfig interface {
	Address() string
	ReadTimeout() time.Duration
	ReadHeaderTimeout() time.Duration
	GRPCTimeout() time.Duration
	ServerTimeout() time.Duration
	ShutdownTimeout() time.Duration
}

type PostgresConfig interface {
	URI() string
	MigrationDirectory() string
}

type InventoryGRPCConfig interface {
	Address() string
}

type PaymentGRPCConfig interface {
	Address() string
}

type KafkaConfig interface {
	Brokers() []string
}

type OrderPaidProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}

type OrderAssembledConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}

type IamGRPCConfig interface {
	Address() string
}
