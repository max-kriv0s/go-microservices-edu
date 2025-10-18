package inventory

import (
	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/repository"
	def "github.com/max-kriv0s/go-microservices-edu/inventory/internal/service"
)

var _ def.InventoryService = (*service)(nil)

type service struct {
	inventoryRepository repository.InventoryRepository
}

func NewService(inventoryRepository repository.InventoryRepository) *service {
	return &service{
		inventoryRepository: inventoryRepository,
	}
}
