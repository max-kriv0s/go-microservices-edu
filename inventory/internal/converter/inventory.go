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

func PartsToProto(parts []model.Part) []*inventoryV1.Part {
	protoParts := make([]*inventoryV1.Part, len(parts))
	for i, part := range parts {
		protoParts[i] = PartToProto(part)
	}
	return protoParts
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

func metadataToProto(metadata map[string]any) map[string]*inventoryV1.Value {
	if metadata == nil {
		return nil
	}

	protoMetadata := make(map[string]*inventoryV1.Value, len(metadata))
	for key, value := range metadata {
		switch val := value.(type) {
		case string:
			protoMetadata[key] = &inventoryV1.Value{Value: &inventoryV1.Value_StringValue{StringValue: val}}
		case int64:
			protoMetadata[key] = &inventoryV1.Value{Value: &inventoryV1.Value_Int64Value{Int64Value: val}}
		case float64:
			protoMetadata[key] = &inventoryV1.Value{Value: &inventoryV1.Value_DoubleValue{DoubleValue: val}}
		case bool:
			protoMetadata[key] = &inventoryV1.Value{Value: &inventoryV1.Value_BoolValue{BoolValue: val}}
		default:
			return nil
		}
	}
	return protoMetadata
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
