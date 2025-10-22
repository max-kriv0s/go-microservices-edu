package converter

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/model"
	repoModel "github.com/max-kriv0s/go-microservices-edu/inventory/internal/repository/model"
	"github.com/stretchr/testify/require"
)

func TestCategoryToRepoCategory(t *testing.T) {
	t.Run("CategoryToRepoCategory - returns default value for unknown category", func(t *testing.T) {

		var category model.Category = 10
		expected := repoModel.CategoryUnknown

		res := categoryToRepoCategory(category)
		require.Equal(t, expected, res)
	})

	t.Run("CategoryToRepoCategory - returns Engine value", func(t *testing.T) {

		category := model.CategoryEngine
		expected := repoModel.CategoryEngine

		res := categoryToRepoCategory(category)
		require.Equal(t, expected, res)
	})

	t.Run("CategoryToRepoCategory - returns Fuel value", func(t *testing.T) {

		category := model.CategoryFuel
		expected := repoModel.CategoryFuel

		res := categoryToRepoCategory(category)
		require.Equal(t, expected, res)
	})

	t.Run("CategoryToRepoCategory - returns Porthole value", func(t *testing.T) {

		category := model.CategoryPorthole
		expected := repoModel.CategoryPorthole

		res := categoryToRepoCategory(category)
		require.Equal(t, expected, res)
	})

	t.Run("CategoryToRepoCategory - returns Wing value", func(t *testing.T) {

		category := model.CategoryWing
		expected := repoModel.CategoryWing

		res := categoryToRepoCategory(category)
		require.Equal(t, expected, res)
	})
}

func TestRepoCategoryToCategory(t *testing.T) {
	t.Run("RepoCategoryToCategory - returns default value for unknown category", func(t *testing.T) {

		var category repoModel.Category = 10
		expected := model.CategoryUnknown

		res := repoCategoryToCategory(category)
		require.Equal(t, expected, res)
	})

	t.Run("RepoCategoryToCategory - returns Engine value", func(t *testing.T) {

		category := repoModel.CategoryEngine
		expected := model.CategoryEngine

		res := repoCategoryToCategory(category)
		require.Equal(t, expected, res)
	})

	t.Run("RepoCategoryToCategory - returns Fuel value", func(t *testing.T) {

		category := repoModel.CategoryFuel
		expected := model.CategoryFuel

		res := repoCategoryToCategory(category)
		require.Equal(t, expected, res)
	})

	t.Run("RepoCategoryToCategory - returns Porthole value", func(t *testing.T) {

		category := repoModel.CategoryPorthole
		expected := model.CategoryPorthole

		res := repoCategoryToCategory(category)
		require.Equal(t, expected, res)
	})

	t.Run("RepoCategoryToCategory - returns Wing value", func(t *testing.T) {

		category := repoModel.CategoryWing
		expected := model.CategoryWing

		res := repoCategoryToCategory(category)
		require.Equal(t, expected, res)
	})
}

func TestDimensionsToRepoDimensions(t *testing.T) {
	t.Run("DimensionsToRepoDimensions - nil returns nil", func(t *testing.T) {
		res := dimensionsToRepoDimensions(nil)
		require.Nil(t, res)
	})

	t.Run("DimensionsToRepoDimensions - convert success", func(t *testing.T) {
		var (
			length = gofakeit.Float64Range(1, 100)
			width  = gofakeit.Float64Range(1, 100)
			height = gofakeit.Float64Range(1, 100)
			weight = gofakeit.Float64Range(0.1, 50)

			repo = &model.Dimensions{
				Length: length,
				Width:  width,
				Height: height,
				Weight: weight,
			}

			expected = &repoModel.Dimensions{
				Length: length,
				Width:  width,
				Height: height,
				Weight: weight,
			}
		)

		res := dimensionsToRepoDimensions(repo)
		require.NotNil(t, res)
		require.Equal(t, expected, res)
	})
}

func TestRepoDimensionsToDimensions(t *testing.T) {
	t.Run("RepoDimensionsToDimensions - nil returns nil", func(t *testing.T) {
		res := repoDimensionsToDimensions(nil)
		require.Nil(t, res)
	})

	t.Run("RepoDimensionsToDimensions - convert success", func(t *testing.T) {
		var (
			length = gofakeit.Float64Range(1, 100)
			width  = gofakeit.Float64Range(1, 100)
			height = gofakeit.Float64Range(1, 100)
			weight = gofakeit.Float64Range(0.1, 50)

			repo = &repoModel.Dimensions{
				Length: length,
				Width:  width,
				Height: height,
				Weight: weight,
			}

			expected = &model.Dimensions{
				Length: length,
				Width:  width,
				Height: height,
				Weight: weight,
			}
		)

		res := repoDimensionsToDimensions(repo)
		require.NotNil(t, res)
		require.Equal(t, expected, res)
	})
}

func TestManufacturerToRepoManufacturer(t *testing.T) {
	t.Run("ManufacturerToRepoManufacturer - nil returns nil", func(t *testing.T) {
		res := manufacturerToRepoManufacturer(nil)
		require.Nil(t, res)
	})

	t.Run("ManufacturerToRepoManufacturer - convert success", func(t *testing.T) {
		var (
			name    = gofakeit.Name()
			country = gofakeit.Country()
			website = gofakeit.URL()

			repo = &model.Manufacturer{
				Name:    name,
				Country: country,
				Website: website,
			}

			expected = &repoModel.Manufacturer{
				Name:    name,
				Country: country,
				Website: website,
			}
		)

		res := manufacturerToRepoManufacturer(repo)
		require.NotNil(t, res)
		require.Equal(t, expected, res)
	})
}

func TestRepoManufacturerToManufacturer(t *testing.T) {
	t.Run("RepoManufacturerToManufacturer - nil returns nil", func(t *testing.T) {
		res := repoManufacturerToManufacturer(nil)
		require.Nil(t, res)
	})

	t.Run("RepoManufacturerToManufacturer - convert success", func(t *testing.T) {
		var (
			name    = gofakeit.Name()
			country = gofakeit.Country()
			website = gofakeit.URL()

			repo = &repoModel.Manufacturer{
				Name:    name,
				Country: country,
				Website: website,
			}

			expected = &model.Manufacturer{
				Name:    name,
				Country: country,
				Website: website,
			}
		)

		res := repoManufacturerToManufacturer(repo)
		require.NotNil(t, res)
		require.Equal(t, expected, res)
	})
}

func TestMetadataToRepoMetadata(t *testing.T) {
	t.Run("MetadataToRepoMetadata - nil returns nil", func(t *testing.T) {
		var m map[string]any
		res := metadataToRepoMetadata(m)
		require.Nil(t, res)
	})

	t.Run("MetadataToRepoMetadata - convert success", func(t *testing.T) {
		m := map[string]any{
			"string": "hello",
			"int64":  int64(42),
			"double": 3.14,
			"bool":   true,
		}

		expected := map[string]any{
			"string": "hello",
			"int64":  int64(42),
			"double": 3.14,
			"bool":   true,
		}

		res := metadataToRepoMetadata(m)
		require.NotNil(t, res)
		require.Equal(t, expected, m)
	})
}

func TestRepoMetadataToMetadata(t *testing.T) {
	t.Run("RepoMetadataToMetadata - nil returns nil", func(t *testing.T) {
		var m map[string]any
		res := repoMetadataToMetadata(m)
		require.Nil(t, res)
	})

	t.Run("RepoMetadataToMetadata - convert success", func(t *testing.T) {
		m := map[string]any{
			"string": "hello",
			"int64":  int64(42),
			"double": 3.14,
			"bool":   true,
		}

		expected := map[string]any{
			"string": "hello",
			"int64":  int64(42),
			"double": 3.14,
			"bool":   true,
		}

		res := repoMetadataToMetadata(m)
		require.NotNil(t, res)
		require.Equal(t, expected, m)
	})
}
