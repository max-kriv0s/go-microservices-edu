package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"

	"github.com/max-kriv0s/go-microservices-edu/assembly/internal/config"
	kafkaConverter "github.com/max-kriv0s/go-microservices-edu/assembly/internal/converter/kafka"
	"github.com/max-kriv0s/go-microservices-edu/assembly/internal/converter/kafka/decoder"
	"github.com/max-kriv0s/go-microservices-edu/assembly/internal/service"
	orderConsumer "github.com/max-kriv0s/go-microservices-edu/assembly/internal/service/consumer/order_consumer"
	orderProducer "github.com/max-kriv0s/go-microservices-edu/assembly/internal/service/producer/order_producer"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/closer"
	wrappedKafka "github.com/max-kriv0s/go-microservices-edu/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/max-kriv0s/go-microservices-edu/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/max-kriv0s/go-microservices-edu/platform/pkg/kafka/producer"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/logger"
	kafkaMiddleware "github.com/max-kriv0s/go-microservices-edu/platform/pkg/middleware/kafka"
)

type diContainer struct {
	orderProducerService service.OrderProducerService
	orderConsumerService service.ConsumerService

	consumerGroup         sarama.ConsumerGroup
	orderRecordedConsumer wrappedKafka.Consumer

	orderPaidDecoder      kafkaConverter.OrderPaidDecoder
	syncProducer          sarama.SyncProducer
	orderRecordedProducer wrappedKafka.Producer
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderProducerService() service.OrderProducerService {
	if d.orderProducerService == nil {
		d.orderProducerService = orderProducer.NewService(d.OrderRecordedProducer())
	}

	return d.orderProducerService
}

func (d *diContainer) OrderConsumerService() service.ConsumerService {
	if d.orderConsumerService == nil {
		d.orderConsumerService = orderConsumer.NewService(d.OrderRecordedConsumer(), d.OrderPaidRecordedDecoder(), d.OrderProducerService())
	}

	return d.orderConsumerService
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.ConsumerGroupId(),
			config.AppConfig().OrderPaidConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.consumerGroup.Close()
		})

		d.consumerGroup = consumerGroup
	}

	return d.consumerGroup
}

func (d *diContainer) OrderRecordedConsumer() wrappedKafka.Consumer {
	if d.orderRecordedConsumer == nil {
		d.orderRecordedConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().OrderPaidConsumer.TopicName(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderRecordedConsumer
}

func (d *diContainer) OrderPaidRecordedDecoder() kafkaConverter.OrderPaidDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}

	return d.orderPaidDecoder
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembleProducer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}

		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return p.Close()
		})

		d.syncProducer = p
	}

	return d.syncProducer
}

func (d *diContainer) OrderRecordedProducer() wrappedKafka.Producer {
	if d.orderRecordedProducer == nil {
		d.orderRecordedProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderAssembleProducer.TopicName(),
			logger.Logger(),
		)
	}

	return d.orderRecordedProducer
}
