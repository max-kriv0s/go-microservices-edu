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

type IamGRPCConfig interface {
	Address() string
}

type PostgresConfig interface {
	URI() string
	MigrationDirectory() string
}

type RedisConfig interface {
	Address() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
}

type SessionConfig interface {
	SessionTTL() time.Duration
}
