package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJson() bool
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
