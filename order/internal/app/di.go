package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1Api "github.com/max-kriv0s/go-microservices-edu/order/api/order/v1"
	client "github.com/max-kriv0s/go-microservices-edu/order/internal/client/grpc"
	inventoryClient "github.com/max-kriv0s/go-microservices-edu/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/max-kriv0s/go-microservices-edu/order/internal/client/grpc/payment/v1"
	"github.com/max-kriv0s/go-microservices-edu/order/internal/config"
	kafkaConverter "github.com/max-kriv0s/go-microservices-edu/order/internal/converter/kafka"
	"github.com/max-kriv0s/go-microservices-edu/order/internal/converter/kafka/decoder"
	"github.com/max-kriv0s/go-microservices-edu/order/internal/repository"
	orderRepository "github.com/max-kriv0s/go-microservices-edu/order/internal/repository/order"
	"github.com/max-kriv0s/go-microservices-edu/order/internal/service"
	orderConsumer "github.com/max-kriv0s/go-microservices-edu/order/internal/service/consumer/order_consumer"
	orderService "github.com/max-kriv0s/go-microservices-edu/order/internal/service/order"
	orderProducer "github.com/max-kriv0s/go-microservices-edu/order/internal/service/producer/order_producer"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/closer"
	wrappedKafka "github.com/max-kriv0s/go-microservices-edu/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/max-kriv0s/go-microservices-edu/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/max-kriv0s/go-microservices-edu/platform/pkg/kafka/producer"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/logger"
	kafkaMiddleware "github.com/max-kriv0s/go-microservices-edu/platform/pkg/middleware/kafka"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/migrator"
	migratorPg "github.com/max-kriv0s/go-microservices-edu/platform/pkg/migrator/pg"
	orderV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	inventoryServiceClient client.InventoryServiceClient

	paymentServiceClient client.PaymentServiceClient

	orderV1Api orderV1.Handler

	orderService         service.OrderService
	orderProducerService service.OrderProducerService
	orderConsumerService service.ConsumerService

	orderRepository repository.OrderRepository

	postgresDBConn *pgx.Conn
	postgresDBPool *pgxpool.Pool

	migrator migrator.Migrator

	syncProducer      sarama.SyncProducer
	orderPaidProducer wrappedKafka.Producer

	orderAssembledConsumer wrappedKafka.Consumer
	consumerGroup          sarama.ConsumerGroup
	orderAssembledDecoder  kafkaConverter.OrderAssembledDecoder
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InventoryServiceClient(ctx context.Context) client.InventoryServiceClient {
	if d.inventoryServiceClient != nil {
		return d.inventoryServiceClient
	}

	inventoryConn, err := grpc.NewClient(
		config.AppConfig().InventoryGRPC.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			timeout.UnaryClientInterceptor(config.AppConfig().OrderHTTP.GRPCTimeout()),
		),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to inventory service: %s\n", err.Error()))
	}

	closer.AddNamed("Inventory service client", func(ctx context.Context) error {
		return inventoryConn.Close()
	})

	d.inventoryServiceClient = inventoryClient.NewInventoryServiceClient(inventoryV1.NewInventoryServiceClient(inventoryConn))

	return d.inventoryServiceClient
}

func (d *diContainer) PaymentServiceClient(ctx context.Context) client.PaymentServiceClient {
	if d.paymentServiceClient != nil {
		return d.paymentServiceClient
	}

	paymentConn, err := grpc.NewClient(
		config.AppConfig().PaymentGRPC.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			timeout.UnaryClientInterceptor(config.AppConfig().OrderHTTP.GRPCTimeout()),
		),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to payment service: %s\n", err.Error()))
	}

	closer.AddNamed("Payment service client", func(ctx context.Context) error {
		return paymentConn.Close()
	})

	d.paymentServiceClient = paymentClient.NewPaymentServiceClient(paymentV1.NewPaymentServiceClient(paymentConn))

	return d.paymentServiceClient
}

func (d *diContainer) OrderV1Api(ctx context.Context) orderV1.Handler {
	if d.orderV1Api == nil {
		d.orderV1Api = orderV1Api.NewAPI(d.OrderService(ctx))
	}

	return d.orderV1Api
}

func (d *diContainer) OrderService(ctx context.Context) service.OrderService {
	if d.orderService == nil {
		d.orderService = orderService.NewService(
			d.InventoryServiceClient(ctx),
			d.PaymentServiceClient(ctx),
			d.OrderRepository(ctx),
			d.OrderProducerService(),
		)
	}

	return d.orderService
}

func (d *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {
	if d.orderRepository == nil {
		d.orderRepository = orderRepository.NewRepository(d.PostgresDBPool(ctx))
	}

	return d.orderRepository
}

func (d *diContainer) PostgresDBConn(ctx context.Context) *pgx.Conn {
	if d.postgresDBConn != nil {
		return d.postgresDBConn
	}

	dbURI := config.AppConfig().Postgres.URI()

	dbConn, err := pgx.Connect(ctx, dbURI)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v\n", err.Error()))
	}

	closer.AddNamed("Postgres database connect", func(ctx context.Context) error {
		return dbConn.Close(ctx)
	})

	err = dbConn.Ping(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to ping postgresDB: %v\n", err.Error()))
	}

	d.postgresDBConn = dbConn

	return d.postgresDBConn
}

func (d *diContainer) PostgresDBPool(ctx context.Context) *pgxpool.Pool {
	if d.postgresDBPool != nil {
		return d.postgresDBPool
	}

	dbURI := config.AppConfig().Postgres.URI()

	// Создаем пул соединений с базой данных
	dbPool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v\n", err.Error()))
	}

	closer.AddNamed("", func(ctx context.Context) error {
		dbPool.Close()
		return nil
	})

	// Проверяем, что соединение с базой установлено
	err = dbPool.Ping(ctx)
	if err != nil {
		panic(fmt.Sprintf("База данных недоступна: %v\n", err.Error()))
	}

	d.postgresDBPool = dbPool

	return d.postgresDBPool
}

func (d *diContainer) Migrator(ctx context.Context) migrator.Migrator {
	if d.migrator != nil {
		return d.migrator
	}
	pgxConfig := d.PostgresDBConn(ctx).Config().Copy()

	migrationsDir := config.AppConfig().Postgres.MigrationDirectory()

	d.migrator = migratorPg.NewMigrator(stdlib.OpenDB(*pgxConfig), migrationsDir)

	return d.migrator
}

func (d *diContainer) OrderProducerService() service.OrderProducerService {
	if d.orderProducerService == nil {
		d.orderProducerService = orderProducer.NewService(d.OrderPaidProducer())
	}

	return d.orderProducerService
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidProducer.Config(),
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

func (d *diContainer) OrderPaidProducer() wrappedKafka.Producer {
	if d.orderPaidProducer == nil {
		d.orderPaidProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderPaidProducer.Topic(),
			logger.Logger(),
		)
	}

	return d.orderPaidProducer
}

func (d *diContainer) OrderConsumerService(ctx context.Context) service.ConsumerService {
	if d.orderConsumerService == nil {
		d.orderConsumerService = orderConsumer.NewService(d.OrderAssembledConsumer(), d.OrderAssembledDecoder(), d.OrderService(ctx))
	}

	return d.orderConsumerService
}

func (d *diContainer) OrderAssembledConsumer() wrappedKafka.Consumer {
	if d.orderAssembledConsumer == nil {
		d.orderAssembledConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().OrderAssembledConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderAssembledConsumer
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledConsumer.GroupID(),
			config.AppConfig().OrderAssembledConsumer.Config(),
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

func (d *diContainer) OrderAssembledDecoder() kafkaConverter.OrderAssembledDecoder {
	if d.orderAssembledDecoder == nil {
		d.orderAssembledDecoder = decoder.NewOrderAssembledDecoder()
	}

	return d.orderAssembledDecoder
}
