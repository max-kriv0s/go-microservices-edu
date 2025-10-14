package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	paymentV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/payment/v1"
)

const (
	grpcHost = "localhost"
	grpcPort = 50052
)

type paymentService struct {
	paymentV1.UnimplementedPaymentServiceServer
}

func (p *paymentService) PayOrder(_ context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	newUUID := uuid.NewString()
	log.Printf("Оплата прошла успешно, transaction_uuid: %s", newUUID)

	return &paymentV1.PayOrderResponse{
		TransactionUuid: newUUID,
	}, nil
}

func main() {
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
	service := &paymentService{}
	paymentV1.RegisterPaymentServiceServer(server, service)

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
