package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/max-kriv0s/go-microservices-edu/order/converter"
	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
	orderV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequestDto) (orderV1.CreateOrderRes, error) {
	order, err := a.orderService.CreateOrder(ctx, converter.CreateOrderRequestToModel(req))
	if err != nil {
		if errors.Is(err, model.ErrBadRequest) {
			return &orderV1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}, nil
		}
		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, nil
	}
	return converter.OrderToCreateOrderResponseDto(order), nil
}
