package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
	orderV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/openapi/order/v1"
)

func (a *api) APIV1OrdersOrderUUIDCancelPost(ctx context.Context, params orderV1.APIV1OrdersOrderUUIDCancelPostParams) (orderV1.APIV1OrdersOrderUUIDCancelPostRes, error) {
	err := a.orderService.CancelOrder(ctx, params.OrderUUID.String())
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}, nil
		case errors.Is(err, model.ErrConflict):
			return &orderV1.ConflictError{
				Code:    http.StatusConflict,
				Message: err.Error(),
			}, nil
		default:
			return &orderV1.InternalServerError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}, nil
		}
	}
	return &orderV1.APIV1OrdersOrderUUIDCancelPostNoContent{}, nil
}
