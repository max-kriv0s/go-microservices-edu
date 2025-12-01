package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

const (
	LogOutputStdout = "stdout"
	LogOutputOTLP   = "otlp"
)

type loggerEnvConfig struct {
	Level        string   `env:"LOGGER_LEVEL,required"`
	AsJson       bool     `env:"LOGGER_AS_JSON"`
	LogOutputs   []string `env:"LOG_OUTPUTS,required" envSeparator:","`
	OTLPEndpoint string   `env:"OTEL_COLLECTOR_ENDPOINT"`
	ServiceName  string   `env:"SERVICE_NAME,required"`
	ServiceEnv   string   `env:"SERVICE_ENV,required"`
}

type loggerConfig struct {
	raw loggerEnvConfig
}

func NewLoggerConfig() (*loggerConfig, error) {
	var raw loggerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	loggerCfg := &loggerConfig{raw: raw}

	// Валидация: если включён OTLP, должны быть все необходимые поля
	if loggerCfg.hasOutput(LogOutputOTLP) {
		if loggerCfg.OTLPEndpoint() == "" {
			return nil, fmt.Errorf("OTEL_COLLECTOR_ENDPOINT is required when OTLP output is enabled")
		}
	}

	return loggerCfg, nil
}

func (cfg *loggerConfig) Level() string {
	return cfg.raw.Level
}

func (cfg *loggerConfig) AsJson() bool {
	return cfg.raw.AsJson
}

func (cfg *loggerConfig) EnableStdout() bool {
	return cfg.hasOutput(LogOutputStdout)
}

func (cfg *loggerConfig) EnableOTLP() bool {
	return cfg.hasOutput(LogOutputOTLP)
}

func (cfg *loggerConfig) OTLPEndpoint() string {
	return cfg.raw.OTLPEndpoint
}

func (cfg *loggerConfig) ServiceName() string {
	return cfg.raw.ServiceName
}

func (cfg *loggerConfig) ServiceEnv() string {
	return cfg.raw.ServiceEnv
}

func (cfg *loggerConfig) hasOutput(target string) bool {
	for _, out := range cfg.raw.LogOutputs {
		if out == target {
			return true
		}
	}
	return false
}
