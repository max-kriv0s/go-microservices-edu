package v1

import (
	"time"

	def "github.com/max-kriv0s/go-microservices-edu/order/internal/client/grpc"
	inventoryV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
)

var _ def.InventoryServiceClient = (*inventoryServiceClient)(nil)

type inventoryServiceClient struct {
	grpcTimeout time.Duration
	client      inventoryV1.InventoryServiceClient
}

func NewInventoryServiceClient(conn *grpc.ClientConn) *inventoryServiceClient {
	return &inventoryServiceClient{
		grpcTimeout: 2 * time.Second,
		client:      inventoryV1.NewInventoryServiceClient(conn),
	}
}
