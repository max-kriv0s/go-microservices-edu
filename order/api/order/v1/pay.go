package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/converter"
	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
	orderV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderV1.PayOrderRequestDto, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	transactionUUID, err := a.orderService.PayOrder(ctx, params.OrderUUID.String(), converter.ApiPaymentMethodToPaymentMethod(req.PaymentMethod))
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}, nil
		} else if errors.Is(err, model.ErrConflict) {
			return &orderV1.ConflictError{
				Code:    http.StatusConflict,
				Message: err.Error(),
			}, nil
		}
		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, nil
	}

	return &orderV1.PayOrderResponseDto{
		TransactionUUID: transactionUUID,
	}, nil
}
