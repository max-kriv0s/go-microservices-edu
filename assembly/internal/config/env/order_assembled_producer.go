package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type orderAssembleProducerCfgConfig struct {
	TopicName string `env:"ORDER_ASSEMBLED_TOPIC_NAME,required"`
}

type orderAssembleProducerConfig struct {
	raw orderAssembleProducerCfgConfig
}

func NewOrderAssembleProducerConfig() (*orderAssembleProducerConfig, error) {
	var raw orderAssembleProducerCfgConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderAssembleProducerConfig{raw: raw}, nil
}

func (cfg *orderAssembleProducerConfig) TopicName() string {
	return cfg.raw.TopicName
}

// Config возвращает конфигурацию для sarama consumer
func (cfg *orderAssembleProducerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Producer.Return.Successes = true

	return config
}
