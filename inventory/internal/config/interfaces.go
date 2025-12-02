package config

import "time"

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

type InventoryGRPCConfig interface {
	Address() string
	GRPCTimeout() time.Duration
}

type MongoConfig interface {
	URI() string
	DatabaseName() string
}

type IamGRPCConfig interface {
	Address() string
}
