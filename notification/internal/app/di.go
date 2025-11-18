package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/go-telegram/bot"

	httpClient "github.com/max-kriv0s/go-microservices-edu/notification/internal/client/http"
	telegramClient "github.com/max-kriv0s/go-microservices-edu/notification/internal/client/http/telegram"
	"github.com/max-kriv0s/go-microservices-edu/notification/internal/config"
	kafkaConverter "github.com/max-kriv0s/go-microservices-edu/notification/internal/converter/kafka"
	"github.com/max-kriv0s/go-microservices-edu/notification/internal/converter/kafka/decoder"
	"github.com/max-kriv0s/go-microservices-edu/notification/internal/service"
	orderAssembledConsumer "github.com/max-kriv0s/go-microservices-edu/notification/internal/service/consumer/order_assembled_consumer"
	orderPaidConsumer "github.com/max-kriv0s/go-microservices-edu/notification/internal/service/consumer/order_paid_consumer"
	telegramService "github.com/max-kriv0s/go-microservices-edu/notification/internal/service/telegram"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/closer"
	wrappedKafka "github.com/max-kriv0s/go-microservices-edu/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/max-kriv0s/go-microservices-edu/platform/pkg/kafka/consumer"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/logger"
	kafkaMiddleware "github.com/max-kriv0s/go-microservices-edu/platform/pkg/middleware/kafka"
)

type diContainer struct {
	telegramService               service.TelegramService
	orderAssembledConsumerService service.ConsumerService
	orderPaidConsumerService      service.ConsumerService

	orderAssembledConsumerGroup sarama.ConsumerGroup
	orderPaidConsumerGroup      sarama.ConsumerGroup

	orderAssembledRecordedConsumer wrappedKafka.Consumer
	orderPaidRecordedConsumer      wrappedKafka.Consumer

	orderAssembledDecoder kafkaConverter.OrderAssembledDecoder
	orderPaidDecoder      kafkaConverter.OrderPaidDecoder

	telegramClient httpClient.TelegramClient
	telegramBot    *bot.Bot
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderAssembledConsumerService(ctx context.Context) service.ConsumerService {
	if d.orderAssembledConsumerService == nil {
		d.orderAssembledConsumerService = orderAssembledConsumer.NewService(
			d.OrderAssembledRecordedConsumer(), d.OrderAssembledDecoder(), d.TelegramService(ctx),
		)
	}

	return d.orderAssembledConsumerService
}

func (d *diContainer) OrderAssembledRecordedConsumer() wrappedKafka.Consumer {
	if d.orderAssembledRecordedConsumer == nil {
		d.orderAssembledRecordedConsumer = wrappedKafkaConsumer.NewConsumer(
			d.OrderAssembledConsumerGroup(),
			[]string{
				config.AppConfig().OrderAssembledConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderAssembledRecordedConsumer
}

func (d *diContainer) OrderAssembledConsumerGroup() sarama.ConsumerGroup {
	if d.orderAssembledConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledConsumer.GroupID(),
			config.AppConfig().OrderAssembledConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create orderAssembledConsumerGroup group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka orderAssembledConsumerGroup", func(ctx context.Context) error {
			return d.orderAssembledConsumerGroup.Close()
		})

		d.orderAssembledConsumerGroup = consumerGroup
	}

	return d.orderAssembledConsumerGroup
}

func (d *diContainer) OrderAssembledDecoder() kafkaConverter.OrderAssembledDecoder {
	if d.orderAssembledDecoder == nil {
		d.orderAssembledDecoder = decoder.NewOrderAssembledDecoder()
	}

	return d.orderAssembledDecoder
}

func (d *diContainer) OrderPaidConsumerService(ctx context.Context) service.ConsumerService {
	if d.orderPaidConsumerService == nil {
		d.orderPaidConsumerService = orderPaidConsumer.NewService(
			d.OrderPaidRecordedConsumer(), d.OrderPaidDecoder(), d.TelegramService(ctx),
		)
	}

	return d.orderPaidConsumerService
}

func (d *diContainer) OrderPaidRecordedConsumer() wrappedKafka.Consumer {
	if d.orderPaidRecordedConsumer == nil {
		d.orderPaidRecordedConsumer = wrappedKafkaConsumer.NewConsumer(
			d.OrderPaidConsumerGroup(),
			[]string{
				config.AppConfig().OrderPaidConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderPaidRecordedConsumer
}

func (d *diContainer) OrderPaidConsumerGroup() sarama.ConsumerGroup {
	if d.orderPaidConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			config.AppConfig().OrderPaidConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create orderPaidConsumerGroup: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka orderAssembledConsumerGroup", func(ctx context.Context) error {
			return d.orderPaidConsumerGroup.Close()
		})

		d.orderPaidConsumerGroup = consumerGroup
	}

	return d.orderPaidConsumerGroup
}

func (d *diContainer) OrderPaidDecoder() kafkaConverter.OrderPaidDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}

	return d.orderPaidDecoder
}

func (d *diContainer) TelegramService(ctx context.Context) service.TelegramService {
	if d.telegramService == nil {
		d.telegramService = telegramService.NewService(
			d.TelegramClient(ctx),
		)
	}

	return d.telegramService
}

func (d *diContainer) TelegramClient(ctx context.Context) httpClient.TelegramClient {
	if d.telegramClient == nil {
		d.telegramClient = telegramClient.NewClient(d.TelegramBot(ctx))
	}

	return d.telegramClient
}

func (d *diContainer) TelegramBot(ctx context.Context) *bot.Bot {
	if d.telegramBot == nil {
		b, err := bot.New(config.AppConfig().TelegramBot.Token())
		if err != nil {
			panic(fmt.Sprintf("failed to create telegram bot: %s\n", err.Error()))
		}

		d.telegramBot = b
	}

	return d.telegramBot
}
