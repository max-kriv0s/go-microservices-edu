package converter

import (
	"github.com/samber/lo"

	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/model"
	repoModel "github.com/max-kriv0s/go-microservices-edu/inventory/internal/repository/model"
)

func PartToRepoModel(part model.Part) repoModel.Part {
	return repoModel.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      categoryToRepoCategory(part.Category),
		Dimensions:    dimensionsToRepoDimensions(part.Dimensions),
		Manufacturer:  manufacturerToRepoManufacturer(part.Manufacturer),
		Tags:          append([]string(nil), part.Tags...),
		Metadata:      metadataToRepoMetadata(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

func PartToModel(part repoModel.Part) model.Part {
	return model.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      repoCategoryToCategory(part.Category),
		Dimensions:    repoDimensionsToDimensions(part.Dimensions),
		Manufacturer:  repoManufacturerToManufacturer(part.Manufacturer),
		Tags:          append([]string(nil), part.Tags...),
		Metadata:      repoMetadataToMetadata(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

func categoryToRepoCategory(category model.Category) repoModel.Category {
	switch category {
	case model.CategoryEngine:
		return repoModel.CategoryEngine
	case model.CategoryFuel:
		return repoModel.CategoryFuel
	case model.CategoryPorthole:
		return repoModel.CategoryPorthole
	case model.CategoryWing:
		return repoModel.CategoryWing
	default:
		return repoModel.CategoryUnknown
	}
}

func repoCategoryToCategory(category repoModel.Category) model.Category {
	switch category {
	case repoModel.CategoryEngine:
		return model.CategoryEngine
	case repoModel.CategoryFuel:
		return model.CategoryFuel
	case repoModel.CategoryPorthole:
		return model.CategoryPorthole
	case repoModel.CategoryWing:
		return model.CategoryWing
	default:
		return model.CategoryUnknown
	}
}

func dimensionsToRepoDimensions(dimensions *model.Dimensions) *repoModel.Dimensions {
	if dimensions == nil {
		return nil
	}

	model := repoModel.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
	return lo.ToPtr(model)
}

func repoDimensionsToDimensions(dimensions *repoModel.Dimensions) *model.Dimensions {
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

func manufacturerToRepoManufacturer(manufacturer *model.Manufacturer) *repoModel.Manufacturer {
	if manufacturer == nil {
		return nil
	}

	model := repoModel.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
	return lo.ToPtr(model)
}

func repoManufacturerToManufacturer(manufacturer *repoModel.Manufacturer) *model.Manufacturer {
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

func metadataToRepoMetadata(metadata map[string]any) map[string]any {
	if metadata == nil {
		return nil
	}

	repoMetadata := make(map[string]any, len(metadata))
	for key, value := range metadata {
		repoMetadata[key] = value
	}

	return repoMetadata
}

func repoMetadataToMetadata(metadata map[string]any) map[string]any {
	if metadata == nil {
		return nil
	}

	repoMetadata := make(map[string]any, len(metadata))
	for key, value := range metadata {
		repoMetadata[key] = value
	}

	return repoMetadata
}
