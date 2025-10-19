package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1Api "github.com/max-kriv0s/go-microservices-edu/order/api/order/v1"
	inventoryClient "github.com/max-kriv0s/go-microservices-edu/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/max-kriv0s/go-microservices-edu/order/internal/client/grpc/payment/v1"
	orderRepository "github.com/max-kriv0s/go-microservices-edu/order/internal/repository/order"
	orderService "github.com/max-kriv0s/go-microservices-edu/order/internal/service/order"
	orderV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/openapi/order/v1"
)

const (
	inventoryUrl      = "localhost:50051"
	paymentUrl        = "localhost:50052"
	serverTimeout     = 10 * time.Second
	serverHost        = "localhost"
	serverPort        = "8080"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	inventoryConn, err := grpc.NewClient(
		inventoryUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect to inventory service: %v", err)
		return
	}
	defer func() {
		if err := inventoryConn.Close(); err != nil {
			log.Printf("failed to close inventory connection: %v", err)
		}
	}()

	paymentConn, err := grpc.NewClient(
		paymentUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect to payment service: %v", err)
		return
	}
	defer func() {
		if err := paymentConn.Close(); err != nil {
			log.Printf("failed to close payment connection: %v", err)
		}
	}()

	inventoryClient := inventoryClient.NewInventoryServiceClient(inventoryConn)
	paymentClient := paymentClient.NewPaymentServiceClient(paymentConn)

	orderRepository := orderRepository.NewRepository()
	orderService := orderService.NewService(inventoryClient, paymentClient, orderRepository)

	orderApi := orderV1Api.NewAPI(orderService)

	orderServer, err := orderV1.NewServer(orderApi)
	if err != nil {
		log.Printf("ошибка создания сервера OpenAPI: %v", err)
		return
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(serverTimeout))

	r.Mount("/", orderServer)

	serverAddr := net.JoinHostPort(serverHost, serverPort)
	server := &http.Server{
		Addr:              serverAddr,
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на %s\n", serverAddr)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
