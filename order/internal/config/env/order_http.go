package env

import (
	"net"
	"time"

	"github.com/caarlos0/env/v11"
)

type orderHTTPEnvConfig struct {
	Host              string        `env:"HTTP_HOST,required"`
	Port              string        `env:"HTTP_PORT,required"`
	ReadTimeout       time.Duration `env:"HTTP_READ_TIMEOUT,required"`
	ReadHeaderTimeout time.Duration `env:"HTTP_READ_HEADER_TIMEOUT" envDefault:"5s"`
}

type orderHTTPConfig struct {
	raw             orderHTTPEnvConfig
	grpcTimeout     time.Duration
	serverTimeout   time.Duration
	shutdownTimeout time.Duration
}

func NewOrderHTTPConfig() (*orderHTTPConfig, error) {
	var raw orderHTTPEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	grpcTimeout := 5 * time.Second
	serverTimeout := 10 * time.Second
	shutdownTimeout := 10 * time.Second

	return &orderHTTPConfig{
		raw:             raw,
		grpcTimeout:     grpcTimeout,
		serverTimeout:   serverTimeout,
		shutdownTimeout: shutdownTimeout,
	}, nil
}

func (cfg *orderHTTPConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}

func (cfg *orderHTTPConfig) ReadTimeout() time.Duration {
	return cfg.raw.ReadTimeout
}

func (cfg *orderHTTPConfig) ReadHeaderTimeout() time.Duration {
	return cfg.raw.ReadHeaderTimeout
}

func (cfg *orderHTTPConfig) GRPCTimeout() time.Duration {
	return cfg.grpcTimeout
}

func (cfg *orderHTTPConfig) ServerTimeout() time.Duration {
	return cfg.serverTimeout
}

func (cfg *orderHTTPConfig) ShutdownTimeout() time.Duration {
	return cfg.shutdownTimeout
}
