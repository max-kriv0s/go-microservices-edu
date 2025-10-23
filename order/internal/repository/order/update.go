package order

import (
	"context"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
	repoConverter "github.com/max-kriv0s/go-microservices-edu/order/internal/repository/converter"
)

func (r *repository) Update(ctx context.Context, uuid string, updateOrder model.UpdateOrder) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, found := r.data[uuid]
	if !found {
		return model.ErrOrderNotFound
	}

	if updateOrder.UserUUID != nil {
		order.UserUUID = *updateOrder.UserUUID
	}

	if updateOrder.PartsUUIDs != nil {
		order.PartsUUIDs = *updateOrder.PartsUUIDs
	}

	if updateOrder.TotalPrice != nil {
		order.TotalPrice = *updateOrder.TotalPrice
	}

	if updateOrder.TransactionUUID != nil {
		order.TransactionUUID = updateOrder.TransactionUUID
	}

	if updateOrder.PaymentMethod != nil {
		order.PaymentMethod = repoConverter.PaymentMethodToRepoPaymentMethod(updateOrder.PaymentMethod)
	}

	if updateOrder.Status != nil {
		order.Status = repoConverter.StatusToRepoStatus(*updateOrder.Status)
	}

	r.data[order.OrderUUID] = order

	return nil
}
