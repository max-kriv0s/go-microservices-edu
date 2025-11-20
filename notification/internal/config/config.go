package config

import (
	"github.com/joho/godotenv"

	"github.com/max-kriv0s/go-microservices-edu/notification/internal/config/env"
)

var appConfig *config

type config struct {
	Logger                 LoggerConfig
	Kafka                  KafkaConfig
	OrderAssembledConsumer OrderAssembledConsumerConfig
	OrderPaidConsumer      OrderPaidConsumerConfig
	TelegramBot            TelegramBotConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	orderAssembledConsumerCfg, err := env.NewOrderAssembledConsumerConfig()
	if err != nil {
		return err
	}

	orderPaidConsumerCfg, err := env.NewOrderPaidConsumerConfig()
	if err != nil {
		return err
	}

	telegramBotCfg, err := env.NewTelegramBotConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:                 loggerCfg,
		Kafka:                  kafkaCfg,
		OrderAssembledConsumer: orderAssembledConsumerCfg,
		OrderPaidConsumer:      orderPaidConsumerCfg,
		TelegramBot:            telegramBotCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
