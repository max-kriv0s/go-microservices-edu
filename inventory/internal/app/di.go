package app

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	inventoryV1API "github.com/max-kriv0s/go-microservices-edu/inventory/internal/api/inventory/v1"
	client "github.com/max-kriv0s/go-microservices-edu/inventory/internal/client/grpc"
	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/config"
	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/repository"
	inventoryRepository "github.com/max-kriv0s/go-microservices-edu/inventory/internal/repository/inventory"
	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/service"
	inventoryService "github.com/max-kriv0s/go-microservices-edu/inventory/internal/service/inventory"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/closer"
	inventoryV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/inventory/v1"

	iamClient "github.com/max-kriv0s/go-microservices-edu/inventory/internal/client/grpc/iam/v1"
	grpcMiddlevare "github.com/max-kriv0s/go-microservices-edu/platform/pkg/middleware/grpc"
	authV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/auth/v1"
)

type diContainer struct {
	inventoryV1API      inventoryV1.InventoryServiceServer
	inventoryService    service.InventoryService
	inventoryRepository repository.InventoryRepository

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database

	iamServiceClient client.IamServiceClient
	authInterceptor  *grpcMiddlevare.AuthInterceptor
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InventoryV1API(ctx context.Context) inventoryV1.InventoryServiceServer {
	if d.inventoryV1API == nil {
		d.inventoryV1API = inventoryV1API.NewAPI(d.InventoryService(ctx))
	}

	return d.inventoryV1API
}

func (d *diContainer) InventoryService(ctx context.Context) service.InventoryService {
	if d.inventoryService == nil {
		d.inventoryService = inventoryService.NewService(d.InventoryRepository(ctx))
	}

	return d.inventoryService
}

func (d *diContainer) InventoryRepository(ctx context.Context) repository.InventoryRepository {
	if d.inventoryRepository == nil {
		d.inventoryRepository = inventoryRepository.NewRepository(d.MongoDBHandle(ctx))
	}

	return d.inventoryRepository
}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to MongoDB: %s\n", err.Error()))
		}

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			panic(fmt.Sprintf("failed to ping MongoDB: %v\n", err))
		}

		closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		d.mongoDBClient = client
	}

	return d.mongoDBClient
}

func (d *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.MongoDBClient(ctx).Database(config.AppConfig().Mongo.DatabaseName())
	}

	return d.mongoDBHandle
}

func (d *diContainer) IamServiceClient() client.IamServiceClient {
	if d.iamServiceClient != nil {
		return d.iamServiceClient
	}

	iamConn, err := grpc.NewClient(
		config.AppConfig().IamGRPC.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			timeout.UnaryClientInterceptor(config.AppConfig().InventoryGRPC.GRPCTimeout()),
		),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to iam service: %s\n", err.Error()))
	}

	closer.AddNamed("Inventory service client", func(ctx context.Context) error {
		return iamConn.Close()
	})

	d.iamServiceClient = iamClient.NewIamServiceClient(authV1.NewAuthServiceClient(iamConn))

	return d.iamServiceClient
}

func (d *diContainer) AuthInterceptor() *grpcMiddlevare.AuthInterceptor {
	if d.authInterceptor == nil {
		d.authInterceptor = grpcMiddlevare.NewAuthInterceptor(d.IamServiceClient())
	}

	return d.authInterceptor
}
