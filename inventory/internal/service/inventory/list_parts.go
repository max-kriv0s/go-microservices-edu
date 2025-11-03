package inventory

import (
	"context"

	"go.uber.org/zap"

	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/model"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/logger"
)

func (s *service) ListParts(ctx context.Context, filter *model.PartsFilter) ([]model.Part, error) {
	parts, err := s.inventoryRepository.ListParts(ctx, filter)
	if err != nil {
		logger.Error(ctx, "internal error list parts", zap.String("func", "ListParts"), zap.Any("filter", filter), zap.Error(err))
		return nil, model.ErrInternalServer
	}

	return parts, nil
}
