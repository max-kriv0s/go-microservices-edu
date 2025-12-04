package env

import (
	"github.com/caarlos0/env/v11"
)

const (
	LogOutputStdout = "stdout"
	LogOutputOTLP   = "otlp"
)

type loggerEnvConfig struct {
	Level      string   `env:"LOGGER_LEVEL,required"`
	AsJson     bool     `env:"LOGGER_AS_JSON"`
	LogOutputs []string `env:"LOG_OUTPUTS,required" envSeparator:","`
}

type loggerConfig struct {
	raw loggerEnvConfig
}

func NewLoggerConfig() (*loggerConfig, error) {
	var raw loggerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &loggerConfig{raw: raw}, nil
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

func (cfg *loggerConfig) hasOutput(target string) bool {
	for _, out := range cfg.raw.LogOutputs {
		if out == target {
			return true
		}
	}
	return false
}
