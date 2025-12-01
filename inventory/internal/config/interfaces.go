package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJson() bool
	EnableStdout() bool
	EnableOTLP() bool
	OTLPEndpoint() string
	ServiceName() string
	ServiceEnv() string
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
