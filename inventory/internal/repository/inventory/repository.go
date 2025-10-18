package inventory

import (
	"context"
	"sync"

	def "github.com/max-kriv0s/go-microservices-edu/inventory/internal/repository"
	repoModel "github.com/max-kriv0s/go-microservices-edu/inventory/internal/repository/model"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data map[string]repoModel.Part
}

func NewRepository() *repository {
	repo := &repository{
		data: make(map[string]repoModel.Part, 0),
	}

	repo.seed(context.Background(), 10)
	return repo
}
