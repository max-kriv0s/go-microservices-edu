package v1

import (
	"context"
	"log"

	clientConverter "github.com/max-kriv0s/go-microservices-edu/order/internal/client/converter"
	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
	inventoryV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/inventory/v1"
)

func (c *inventoryServiceClient) ListParts(ctx context.Context, partsUUIDs []string) ([]model.Part, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, c.grpcTimeout)
	defer cancel()

	inventoryReq := &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Uuids: append([]string(nil), partsUUIDs...),
		},
	}
	listParts, err := c.client.ListParts(ctxWithTimeout, inventoryReq)
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
