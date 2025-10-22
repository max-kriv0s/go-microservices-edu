package inventory

import (
	"context"
	"slices"

	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/model"
	repoConverter "github.com/max-kriv0s/go-microservices-edu/inventory/internal/repository/converter"
)

func (r *repository) FindAll(ctx context.Context, filter *model.PartsFilter) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	parts := make([]model.Part, 0, len(r.data))
	for _, repoPart := range r.data {
		part := repoConverter.PartToModel(repoPart)
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
