package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/max-kriv0s/go-microservices-edu/order/converter"
	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
	orderV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/openapi/order/v1"
)

func (a *api) APIV1OrdersOrderUUIDGet(ctx context.Context, params orderV1.APIV1OrdersOrderUUIDGetParams) (orderV1.APIV1OrdersOrderUUIDGetRes, error) {
	order, err := a.orderService.GetOrder(ctx, params.OrderUUID.String())
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}, nil
		}
		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, nil

	}

	res, err := converter.OrderToGetResponseDto(order)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, nil
	}
	return res, nil
}
