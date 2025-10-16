package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/max-kriv0s/go-microservices-edu/payment/internal/converter"
	paymentV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	transactionUuid, err := a.paymentService.PayOrder(ctx, converter.PayDtoToModel(req))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}

	return &paymentV1.PayOrderResponse{
		TransactionUuid: transactionUuid,
	}, nil
}
