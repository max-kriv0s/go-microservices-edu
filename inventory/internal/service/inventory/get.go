package inventory

import (
	"context"
	"errors"

	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/model"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) GetPart(ctx context.Context, uuid string) (model.Part, error) {
	part, err := s.inventoryRepository.Get(ctx, uuid)
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return model.Part{}, model.ErrPartNotFound
		}

		logger.Error(ctx, "internal error getting part", zap.String("func", "GetPart"), zap.String("uuid", uuid), zap.Error(err))

		return model.Part{}, model.ErrInternalServer
	}

	return part, nil
}
