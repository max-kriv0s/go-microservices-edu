package inventory

import (
	"context"
	"errors"
	"log"

	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/model"
)

func (s *service) GetPart(ctx context.Context, uuid string) (model.Part, error) {
	part, err := s.inventoryRepository.Get(ctx, uuid)
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return model.Part{}, model.ErrPartNotFound
		}

		log.Printf("[service.GetPart] internal error getting part (uuid=%s): %v", uuid, err)

		return model.Part{}, model.ErrInternalServer
	}

	return part, nil
}
