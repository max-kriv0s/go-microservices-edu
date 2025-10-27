package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryV1API "github.com/max-kriv0s/go-microservices-edu/inventory/internal/api/inventory/v1"
	inventoryRepository "github.com/max-kriv0s/go-microservices-edu/inventory/internal/repository/inventory"
	inventoryService "github.com/max-kriv0s/go-microservices-edu/inventory/internal/service/inventory"
	inventoryV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/inventory/v1"
)

const (
	grpcHost = "localhost"
	grpcPort = 50051
)

func main() {
	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("failed to load .env file: %v\n", err)
		return
	}

	dbURI := os.Getenv("INVENTORY_MONGO_URI")

	// Создаем клиент MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return
	}
	defer func() {
		cerr := client.Disconnect(ctx)
		if cerr != nil {
			log.Printf("failed to disconnect: %v\n", cerr)
		}
	}()

	// Проверяем соединение с базой данных
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("failed to ping database: %v\n", err)
		return
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", grpcHost, grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	server := grpc.NewServer()

	dbName := os.Getenv("INVENTORY_DATABASE")
	db := client.Database(dbName)

	repository := inventoryRepository.NewRepository(db)

	service := inventoryService.NewService(repository)
	api := inventoryV1API.NewAPI(service)

	inventoryV1.RegisterInventoryServiceServer(server, api)

	// Включаем рефлексию для отладки
	reflection.Register(server)

	go func() {
		log.Printf("🚀 gRPC server listening on %s:%d\n", grpcHost, grpcPort)
		err := server.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	server.GracefulStop()
	log.Println("✅ Server stopped")
}
