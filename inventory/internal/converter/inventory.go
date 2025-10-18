package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/model"
	inventoryV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/inventory/v1"
)

func PartToProto(part model.Part) *inventoryV1.Part {
	var createdAt *timestamppb.Timestamp
	if part.CreatedAt != nil {
		createdAt = timestamppb.New(*part.CreatedAt)
	}

	var updatedAt *timestamppb.Timestamp
	if part.UpdatedAt != nil {
		updatedAt = timestamppb.New(*part.UpdatedAt)
	}

	return &inventoryV1.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      categoryToProto(part.Category),
		Dimensions:    dimensionsToProto(part.Dimensions),
		Manufacturer:  manufacturerToProto(part.Manufacturer),
		Tags:          append([]string(nil), part.Tags...),
		Metadata:      metadataToProto(part.Metadata),
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func categoryToProto(category model.Category) inventoryV1.Category {
	switch category {
	case model.CategoryEngine:
		return inventoryV1.Category_ENGINE
	case model.CategoryFuel:
		return inventoryV1.Category_FUEL
	case model.CategoryPorthole:
		return inventoryV1.Category_PORTHOLE
	case model.CategoryWing:
		return inventoryV1.Category_WING
	default:
		return inventoryV1.Category_UNKNOWN
	}
}

func dimensionsToProto(dimensions *model.Dimensions) *inventoryV1.Dimensions {
	if dimensions == nil {
		return nil
	}
	return &inventoryV1.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func manufacturerToProto(manufacturer *model.Manufacturer) *inventoryV1.Manufacturer {
	if manufacturer == nil {
		return nil
	}

	return &inventoryV1.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func metadataToProto(metadata map[string]*model.Value) map[string]*inventoryV1.Value {
	if metadata == nil {
		return nil
	}

	protoMetadata := make(map[string]*inventoryV1.Value, len(metadata))
	for key, value := range metadata {
		protoMetadata[key] = valueToProto(value)
	}
	return protoMetadata
}

func valueToProto(value *model.Value) *inventoryV1.Value {
	if value == nil {
		return nil
	}

	switch {
	case value.String != nil:
		return &inventoryV1.Value{Value: &inventoryV1.Value_StringValue{StringValue: *value.String}}
	case value.Int64 != nil:
		return &inventoryV1.Value{Value: &inventoryV1.Value_Int64Value{Int64Value: *value.Int64}}
	case value.Double != nil:
		return &inventoryV1.Value{Value: &inventoryV1.Value_DoubleValue{DoubleValue: *value.Double}}
	case value.Bool != nil:
		return &inventoryV1.Value{Value: &inventoryV1.Value_BoolValue{BoolValue: *value.Bool}}
	default:
		return nil
	}
}

func PartsFilterToModel(filter *inventoryV1.PartsFilter) *model.PartsFilter {
	if filter == nil {
		return nil
	}

	return &model.PartsFilter{
		Uuids:                 append([]string(nil), filter.Uuids...),
		Names:                 append([]string(nil), filter.Names...),
		Categories:            categoriesToModel(filter.Categories),
		ManufacturerCountries: append([]string(nil), filter.ManufacturerCountries...),
		Tags:                  append([]string(nil), filter.Tags...),
	}
}

func categoriesToModel(categories []inventoryV1.Category) []model.Category {
	if categories == nil {
		return nil
	}

	modelCategories := make([]model.Category, len(categories))
	for i, category := range categories {
		modelCategories[i] = categoryToModel(category)
	}

	return modelCategories
}

func categoryToModel(category inventoryV1.Category) model.Category {
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
