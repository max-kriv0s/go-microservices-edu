package converter

import (
	"time"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
	inventoryV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/inventory/v1"
	"github.com/samber/lo"
)

func ClientPartToModel(part *inventoryV1.Part) *model.Part {
	if part == nil {
		return nil
	}

	var createdAt *time.Time
	if part.CreatedAt != nil {
		createdAt = lo.ToPtr(part.CreatedAt.AsTime())
	}

	var updatedAt *time.Time
	if part.UpdatedAt != nil {
		updatedAt = lo.ToPtr(part.UpdatedAt.AsTime())
	}

	return &model.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      fromClientToCategory(part.Category),
		Dimensions:    fromClientToDimensions(part.Dimensions),
		Manufacturer:  fromClientToManufacturer(part.Manufacturer),
		Tags:          append([]string(nil), part.Tags...),
		Metadata:      fromClientToMetadata(part.Metadata),
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func fromClientToCategory(category inventoryV1.Category) model.Category {
	switch category {
	case inventoryV1.Category_ENGINE:
		return model.CategoryEngine
	case inventoryV1.Category_FUEL:
		return model.CategoryFuel
	case inventoryV1.Category_PORTHOLE:
		return model.CategoryPorthole
	case inventoryV1.Category_WING:
		return model.CategoryWing
	default:
		return model.CategoryUnknown
	}
}

func fromClientToDimensions(dimensions *inventoryV1.Dimensions) *model.Dimensions {
	if dimensions == nil {
		return nil
	}

	model := model.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
	return lo.ToPtr(model)
}

func fromClientToManufacturer(manufacturer *inventoryV1.Manufacturer) *model.Manufacturer {
	if manufacturer == nil {
		return nil
	}

	model := model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
	return lo.ToPtr(model)
}

func fromClientToMetadata(metadata map[string]*inventoryV1.Value) map[string]*model.Value {
	if metadata == nil {
		return nil
	}

	repoMetadata := make(map[string]*model.Value, len(metadata))

	for key, value := range metadata {
		repoMetadata[key] = fromClientToValue(value)
	}
	return repoMetadata
}

func fromClientToValue(pb *inventoryV1.Value) *model.Value {
	if pb == nil {
		return nil
	}

	v := &model.Value{}

	switch val := pb.Value.(type) {
	case *inventoryV1.Value_StringValue:
		v.String = &val.StringValue
	case *inventoryV1.Value_Int64Value:
		v.Int64 = &val.Int64Value
	case *inventoryV1.Value_DoubleValue:
		v.Double = &val.DoubleValue
	case *inventoryV1.Value_BoolValue:
		v.Bool = &val.BoolValue
	default:
	}

	return v
}
