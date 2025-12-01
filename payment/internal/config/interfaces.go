package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
	EnableStdout() bool
	EnableOTLP() bool
	OTLPEndpoint() string
	ServiceName() string
	ServiceEnv() string
}

type PaymentGRPCConfig interface {
	Address() string
}
