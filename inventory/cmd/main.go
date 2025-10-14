package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	inventoryV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/inventory/v1"
)

const (
	grpcHost = "localhost"
	grpcPort = 50051
)

type InventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer

	mu    sync.RWMutex
	parts map[string]*inventoryV1.Part
}

func NewInventoryService() *InventoryService {
	return &InventoryService{
		parts: make(map[string]*inventoryV1.Part),
	}
}

func (i *InventoryService) AddPart(part *inventoryV1.Part) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.parts[part.Uuid] = part
}

func (i *InventoryService) GetPart(_ context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	part, ok := i.parts[req.Uuid]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.GetUuid())
	}

	return &inventoryV1.GetPartResponse{
		Part: part,
	}, nil
}

func (i *InventoryService) ListParts(_ context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	parts := make([]*inventoryV1.Part, 0)
	for _, part := range i.parts {
		if !matchFilters(req.Filter, part) {
			continue
		}
		parts = append(parts, part)
	}

	return &inventoryV1.ListPartsResponse{
		Parts: parts,
	}, nil
}

func matchFilters(filter *inventoryV1.PartsFilter, part *inventoryV1.Part) bool {
	if filter == nil {
		return true
	}

	return hasString(filter.Uuids, part.Uuid) &&
		hasString(filter.Names, part.Name) &&
		hasCategory(filter.Categories, part.Category) &&
		hasString(filter.ManufacturerCountries, part.Manufacturer.Country) &&
		hasAnyTags(filter.Tags, part.Tags)
}

func hasString(filter []string, value string) bool {
	if len(filter) == 0 {
		return true
	}
	return slices.Contains(filter, value)
}

func hasCategory(filterCategories []inventoryV1.Category, category inventoryV1.Category) bool {
	if len(filterCategories) == 0 {
		return true
	}
	return slices.Contains(filterCategories, category)
}

func hasAnyTags(filterTags, tags []string) bool {
	if len(filterTags) == 0 {
		return true
	}

	tagSet := make(map[string]struct{}, len(tags))
	for _, tag := range tags {
		tagSet[tag] = struct{}{}
	}

	for _, filterTag := range filterTags {
		if _, ok := tagSet[filterTag]; ok {
			return true
		}
	}
	return false
}

// helper для заполнения поля Value (oneof)
func randomValue() *inventoryV1.Value {
	//nolint:gosec // Используем math/rand для некритичных целей, OK
	switch rand.Intn(4) {
	case 0:
		return &inventoryV1.Value{Value: &inventoryV1.Value_StringValue{StringValue: gofakeit.Word()}}
	case 1:
		return &inventoryV1.Value{Value: &inventoryV1.Value_Int64Value{Int64Value: gofakeit.Int64()}}
	case 2:
		return &inventoryV1.Value{Value: &inventoryV1.Value_DoubleValue{DoubleValue: gofakeit.Float64()}}
	default:
		return &inventoryV1.Value{Value: &inventoryV1.Value_BoolValue{BoolValue: gofakeit.Bool()}}
	}
}

func generatePart() *inventoryV1.Part {
	now := time.Now()

	// Случайный Category (не UNKNOWN)
	//nolint:gosec // Используем math/rand для некритичных целей, OK
	category := inventoryV1.Category(rand.Int31n(int32(inventoryV1.Category_WING)) + 1)

	// Заполняем map metadata с несколькими парами
	metadata := make(map[string]*inventoryV1.Value)
	for i := 0; i < 3; i++ {
		key := gofakeit.Word()
		metadata[key] = randomValue()
	}

	// Заполняем теги
	tagsCount := gofakeit.Number(1, 5)
	tags := make([]string, tagsCount)
	for i := 0; i < tagsCount; i++ {
		tags[i] = gofakeit.Word()
	}

	part := &inventoryV1.Part{
		Uuid:          gofakeit.UUID(),
		Name:          gofakeit.ProductName(),
		Description:   gofakeit.Sentence(10),
		Price:         gofakeit.Price(1, 1000),
		StockQuantity: int64(gofakeit.IntRange(0, 100)),
		Category:      category,
		Dimensions: &inventoryV1.Dimensions{
			Length: gofakeit.Float64Range(1, 100),
			Width:  gofakeit.Float64Range(1, 100),
			Height: gofakeit.Float64Range(1, 100),
			Weight: gofakeit.Float64Range(0.1, 50),
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    gofakeit.Company(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		},
		Tags:      tags,
		Metadata:  metadata,
		CreatedAt: timestamppb.New(now.Add(-time.Duration(gofakeit.Number(0, 1000)) * time.Hour)),
		UpdatedAt: timestamppb.New(now),
	}

	return part
}

func seed(service *InventoryService, count int) {
	for i := 0; i < count; i++ {
		part := generatePart()
		service.AddPart(part)
	}
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
	service := NewInventoryService()
	inventoryV1.RegisterInventoryServiceServer(server, service)

	seed(service, 10)

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
