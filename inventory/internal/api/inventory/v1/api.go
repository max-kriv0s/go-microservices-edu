package v1

import (
	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/service"
	inventoryV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventoryV1.UnimplementedInventoryServiceServer

	inventoryService service.InventoryService
}

func NewAPI(inventoryService service.InventoryService) *api {
	return &api{
		inventoryService: inventoryService,
	}
}
