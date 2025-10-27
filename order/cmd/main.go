package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1Api "github.com/max-kriv0s/go-microservices-edu/order/api/order/v1"
	inventoryClient "github.com/max-kriv0s/go-microservices-edu/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/max-kriv0s/go-microservices-edu/order/internal/client/grpc/payment/v1"
	"github.com/max-kriv0s/go-microservices-edu/order/internal/migrator"
	orderRepository "github.com/max-kriv0s/go-microservices-edu/order/internal/repository/order"
	orderService "github.com/max-kriv0s/go-microservices-edu/order/internal/service/order"
	orderV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/payment/v1"
)

const (
	inventoryUrl      = "localhost:50051"
	paymentUrl        = "localhost:50052"
	serverTimeout     = 10 * time.Second
	serverHost        = "localhost"
	serverPort        = "8080"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
	grpcTimeout       = 5 * time.Second
)

func main() {
	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("failed to load .env file: %v\n", err)
		return
	}

	inventoryConn, err := grpc.NewClient(
		inventoryUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			timeout.UnaryClientInterceptor(grpcTimeout),
		),
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
		grpc.WithUnaryInterceptor(
			timeout.UnaryClientInterceptor(grpcTimeout),
		),
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

	inventoryServiceClient := inventoryClient.NewInventoryServiceClient(inventoryV1.NewInventoryServiceClient(inventoryConn))

	paymentServiceClient := paymentClient.NewPaymentServiceClient(paymentV1.NewPaymentServiceClient(paymentConn))

	dbName := os.Getenv("ORDER_DB_DATABASE")
	dbUser := os.Getenv("ORDER_DB_USER")
	dbPassword := os.Getenv("ORDER_DB_PASSWORD")
	dbHost := os.Getenv("ORDER_DB_HOST")
	dbPort := os.Getenv("ORDER_DB_PORT")

	dbURI := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Создаем соединение с базой данных для миграции
	con, err := pgx.Connect(ctx, dbURI)
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return
	}
	defer func() {
		cerr := con.Close(ctx)
		if cerr != nil {
			log.Printf("failed to close connection: %v\n", cerr)
		}
	}()

	// Проверяем, что соединение с базой установлено
	err = con.Ping(ctx)
	if err != nil {
		log.Printf("База данных недоступна: %v\n", err)
		return
	}

	// Инициализируем мигратор
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	migratorRunner := migrator.NewMigrator(stdlib.OpenDB(*con.Config().Copy()), migrationsDir)

	err = migratorRunner.Up()
	if err != nil {
		log.Printf("Ошибка миграции базы данных: %v\n", err)
		return
	}

	// Создаем пул соединений с базой данных
	dbPool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return
	}
	defer dbPool.Close()

	// Проверяем, что соединение с базой установлено
	err = dbPool.Ping(ctx)
	if err != nil {
		log.Printf("База данных недоступна: %v\n", err)
		return
	}

	orderRepository := orderRepository.NewRepository(dbPool)
	orderService := orderService.NewService(inventoryServiceClient, paymentServiceClient, orderRepository)

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
