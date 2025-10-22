package inventory

import (
	"context"

	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/model"
)

func (s *service) ListParts(ctx context.Context, filter *model.PartsFilter) ([]model.Part, error) {
	parts, err := s.inventoryRepository.FindAll(ctx, filter)
	if err != nil {
		return nil, model.ErrInternalServer
	}

	return parts, nil
}
