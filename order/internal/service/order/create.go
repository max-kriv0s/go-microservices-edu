package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
)

func (s *service) CreateOrder(ctx context.Context, data model.CreateOrderRequest) (model.Order, error) {
	listParts, err := s.inventoryClient.ListParts(ctx, data.PartUuids)
	if err != nil {
		return model.Order{}, model.ErrInternalServer
	}
	if len(listParts) != len(data.PartUuids) {
		return model.Order{}, model.NewBadRequestError("All the details were not found")
	}

	var totalPrice float64
	for _, part := range listParts {
		totalPrice += part.Price
	}

	newOrder := model.Order{
		OrderUUID:  uuid.NewString(),
		UserUUID:   data.UserUUID,
		PartsUUIDs: data.PartUuids,
		TotalPrice: totalPrice,
		Status:     model.OrderStatusPendingPayment,
	}
	_, err = s.orderRepository.Create(ctx, newOrder)
	if err != nil {
		return model.Order{}, model.ErrInternalServer
	}
	return newOrder, nil
}
