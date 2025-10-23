package v1

import (
	def "github.com/max-kriv0s/go-microservices-edu/order/internal/client/grpc"
	inventoryV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/inventory/v1"
)

var _ def.InventoryServiceClient = (*inventoryServiceClient)(nil)

type inventoryServiceClient struct {
	client inventoryV1.InventoryServiceClient
}

func NewInventoryServiceClient(client inventoryV1.InventoryServiceClient) *inventoryServiceClient {
	return &inventoryServiceClient{
		client: client,
	}
}
