package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/converter"
	inventoryV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	parts, err := a.inventoryService.ListParts(ctx, converter.PartsFilterToModel(req.GetFilter()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unknown error: %v", err)
	}

	protoParts := converter.PartsToProto(parts)

	return &inventoryV1.ListPartsResponse{
		Parts: protoParts,
	}, nil
}
