package env

import (
	"net"
	"time"

	"github.com/caarlos0/env/v11"
)

type inventoryGRPCEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type inventoryGRPCConfig struct {
	raw         inventoryGRPCEnvConfig
	grpcTimeout time.Duration
}

func NewInventoryGRPCConfig() (*inventoryGRPCConfig, error) {
	var raw inventoryGRPCEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	grpcTimeout := 5 * time.Second

	return &inventoryGRPCConfig{
		raw:         raw,
		grpcTimeout: grpcTimeout,
	}, nil
}

func (cfg *inventoryGRPCConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}

func (cfg *inventoryGRPCConfig) GRPCTimeout() time.Duration {
	return cfg.grpcTimeout
}
