package converter

import (
	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
	inventoryV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/inventory/v1"
)

func ClientPartToModel(part *inventoryV1.Part) *model.Part {
	if part == nil {
		return nil
	}

	return &model.Part{
		Uuid:  part.Uuid,
		Price: part.Price,
	}
}
