package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/converter"
	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/model"
	inventoryV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part, err := a.inventoryService.GetPart(ctx, req.GetUuid())
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return nil, status.Errorf(codes.NotFound, "part not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "unknown error: %v", err)
	}

	return &inventoryV1.GetPartResponse{
		Part: converter.PartToProto(part),
	}, nil
}
