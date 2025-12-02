package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type otelCollectorEnvConfig struct {
	OTLPEndpoint string `env:"OTEL_COLLECTOR_ENDPOINT"`
	ServiceName  string `env:"SERVICE_NAME,required"`
	ServiceEnv   string `env:"SERVICE_ENV,required"`
}

type otelCollectorConfig struct {
	raw otelCollectorEnvConfig
}

func NewOtelCollectroConfig() (*otelCollectorConfig, error) {
	var raw otelCollectorEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &otelCollectorConfig{raw: raw}, nil
}

func (cfg *otelCollectorConfig) CollectorEndpoint() string {
	return cfg.raw.OTLPEndpoint
}

func (cfg *otelCollectorConfig) ServiceName() string {
	return cfg.raw.ServiceName
}

func (cfg *otelCollectorConfig) ServiceEnv() string {
	return cfg.raw.ServiceEnv
}

func (cfg *otelCollectorConfig) CollectorInterval() time.Duration {
	return 10 * time.Second
}
