package v1

import (
	"context"
	"log"

	clientConverter "github.com/max-kriv0s/go-microservices-edu/order/internal/client/converter"
	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
	grpcAuth "github.com/max-kriv0s/go-microservices-edu/platform/pkg/middleware/grpc"
	inventoryV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/inventory/v1"
)

func (c *inventoryServiceClient) ListParts(ctx context.Context, partsUUIDs []string) ([]model.Part, error) {
	// добавляем session UUID в gRPC metadata для передачи в Inventory Service
	ctx = grpcAuth.ForwardSessionUUIDToGRPC(ctx)

	inventoryReq := &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Uuids: append([]string(nil), partsUUIDs...),
		},
	}
	listParts, err := c.client.ListParts(ctx, inventoryReq)
	if err != nil {
		log.Printf("failed client request: %v", err)
		return nil, err
	}

	parts := make([]model.Part, len(listParts.Parts))
	for i, part := range listParts.Parts {
		parts[i] = *clientConverter.ClientPartToModel(part)
	}
	return parts, nil
}
