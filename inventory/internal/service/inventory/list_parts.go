package inventory

import (
	"context"
	"slices"

	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/model"
)

func (s *service) ListParts(ctx context.Context, filter *model.PartsFilter) ([]model.Part, error) {
	allParts, err := s.inventoryRepository.FindAll(ctx)
	if err != nil {
		return nil, model.ErrInternalServer
	}

	parts := make([]model.Part, 0)
	for _, part := range allParts {
		if !matchFilters(filter, part) {
			continue
		}
		parts = append(parts, part)
	}

	return parts, nil
}

func matchFilters(filter *model.PartsFilter, part model.Part) bool {
	if filter == nil {
		return true
	}

	return hasString(filter.Uuids, part.Uuid) &&
		hasString(filter.Names, part.Name) &&
		hasCategory(filter.Categories, part.Category) &&
		hasString(filter.ManufacturerCountries, part.Manufacturer.Country) &&
		hasAnyTags(filter.Tags, part.Tags)
}

func hasString(filter []string, value string) bool {
	if len(filter) == 0 {
		return true
	}
	return slices.Contains(filter, value)
}

func hasCategory(filterCategories []model.Category, category model.Category) bool {
	if len(filterCategories) == 0 {
		return true
	}
	return slices.Contains(filterCategories, category)
}

func hasAnyTags(filterTags, tags []string) bool {
	if len(filterTags) == 0 {
		return true
	}

	tagSet := make(map[string]struct{}, len(tags))
	for _, tag := range tags {
		tagSet[tag] = struct{}{}
	}

	for _, filterTag := range filterTags {
		if _, ok := tagSet[filterTag]; ok {
			return true
		}
	}
	return false
}
