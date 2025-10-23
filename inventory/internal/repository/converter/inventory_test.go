package converter

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"

	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/model"
	repoModel "github.com/max-kriv0s/go-microservices-edu/inventory/internal/repository/model"
)

func (s *ConserterSuite) TestPartToRepoModelSuccess() {
	var (
		uuid          = gofakeit.UUID()
		name          = gofakeit.ProductName()
		description   = gofakeit.Sentence()
		price         = gofakeit.Price(1, 1000)
		stockQuantity = int64(gofakeit.IntRange(0, 100))
		category      = repoModel.CategoryEngine
		dimensions    = &model.Dimensions{
			Length: gofakeit.Float64Range(1, 100),
			Width:  gofakeit.Float64Range(1, 100),
			Height: gofakeit.Float64Range(1, 100),
			Weight: gofakeit.Float64Range(0.1, 50),
		}

		manufacturer = &model.Manufacturer{
			Name:    gofakeit.Company(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		}

		expectedCategory = model.CategoryEngine
	)

	tags := make([]string, 1)
	tags[0] = gofakeit.Word()

	metadata := make(map[string]any)
	metadata["string"] = lo.ToPtr(gofakeit.Word())

	now := gofakeit.Date()

	part := model.Part{
		Uuid:          uuid,
		Name:          name,
		Description:   description,
		Price:         price,
		StockQuantity: stockQuantity,
		Category:      expectedCategory,
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          append([]string(nil), tags...),
		Metadata:      metadata,
		CreatedAt:     lo.ToPtr(now),
		UpdatedAt:     lo.ToPtr(now),
	}

	expectedPart := repoModel.Part{
		Uuid:          uuid,
		Name:          name,
		Description:   description,
		Price:         price,
		StockQuantity: stockQuantity,
		Category:      category,
		Dimensions: &repoModel.Dimensions{
			Length: dimensions.Length,
			Width:  dimensions.Width,
			Height: dimensions.Height,
			Weight: dimensions.Weight,
		},
		Manufacturer: &repoModel.Manufacturer{
			Name:    manufacturer.Name,
			Country: manufacturer.Country,
			Website: manufacturer.Website,
		},
		Tags:      tags,
		Metadata:  metadata,
		CreatedAt: lo.ToPtr(now),
		UpdatedAt: lo.ToPtr(now),
	}

	res := PartToRepoModel(part)
	s.Require().Equal(expectedPart, res)
}

func (s *ConserterSuite) TestPartToModelSuccess() {
	var (
		uuid          = gofakeit.UUID()
		name          = gofakeit.ProductName()
		description   = gofakeit.Sentence()
		price         = gofakeit.Price(1, 1000)
		stockQuantity = int64(gofakeit.IntRange(0, 100))
		category      = repoModel.CategoryEngine
		dimensions    = &repoModel.Dimensions{
			Length: gofakeit.Float64Range(1, 100),
			Width:  gofakeit.Float64Range(1, 100),
			Height: gofakeit.Float64Range(1, 100),
			Weight: gofakeit.Float64Range(0.1, 50),
		}

		manufacturer = &repoModel.Manufacturer{
			Name:    gofakeit.Company(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		}

		expectedCategory = model.CategoryEngine
	)

	tags := make([]string, 1)
	tags[0] = gofakeit.Word()

	metadata := make(map[string]any)
	metadata["string"] = lo.ToPtr(gofakeit.Word())

	now := gofakeit.Date()

	repoPart := repoModel.Part{
		Uuid:          uuid,
		Name:          name,
		Description:   description,
		Price:         price,
		StockQuantity: stockQuantity,
		Category:      category,
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          tags,
		Metadata:      metadata,
		CreatedAt:     lo.ToPtr(now),
		UpdatedAt:     lo.ToPtr(now),
	}

	expectedModel := model.Part{
		Uuid:          uuid,
		Name:          name,
		Description:   description,
		Price:         price,
		StockQuantity: stockQuantity,
		Category:      expectedCategory,
		Dimensions: &model.Dimensions{
			Length: dimensions.Length,
			Width:  dimensions.Width,
			Height: dimensions.Height,
			Weight: dimensions.Weight,
		},
		Manufacturer: &model.Manufacturer{
			Name:    manufacturer.Name,
			Country: manufacturer.Country,
			Website: manufacturer.Website,
		},
		Tags:      append([]string(nil), tags...),
		Metadata:  metadata,
		CreatedAt: lo.ToPtr(now),
		UpdatedAt: lo.ToPtr(now),
	}

	res := PartToModel(repoPart)
	s.Require().Equal(expectedModel, res)
}

func (s *ConserterSuite) TestCategoryToRepoCategoryDefault() {
	var category model.Category = 10
	expected := repoModel.CategoryUnknown

	res := categoryToRepoCategory(category)
	s.Require().Equal(expected, res)
}

func (s *ConserterSuite) TestCategoryToRepoCategoryReturnsEngine() {
	category := model.CategoryEngine
	expected := repoModel.CategoryEngine

	res := categoryToRepoCategory(category)
	s.Require().Equal(expected, res)
}

func (s *ConserterSuite) TestCategoryToRepoCategoryReturnsFuel() {
	category := model.CategoryFuel
	expected := repoModel.CategoryFuel

	res := categoryToRepoCategory(category)
	s.Require().Equal(expected, res)
}

func (s *ConserterSuite) TestCategoryToRepoCategoryReturnsPorthole() {
	category := model.CategoryPorthole
	expected := repoModel.CategoryPorthole

	res := categoryToRepoCategory(category)
	s.Require().Equal(expected, res)
}

func (s *ConserterSuite) TestCategoryToRepoCategoryReturnsWing() {
	category := model.CategoryWing
	expected := repoModel.CategoryWing

	res := categoryToRepoCategory(category)
	s.Require().Equal(expected, res)
}

func (s *ConserterSuite) TestRepoCategoryToCategoryReturnsDefault() {
	var category repoModel.Category = 10
	expected := model.CategoryUnknown

	res := repoCategoryToCategory(category)
	s.Require().Equal(expected, res)
}

func (s *ConserterSuite) TestRepoCategoryToCategoryReturnsEngine() {
	category := repoModel.CategoryEngine
	expected := model.CategoryEngine

	res := repoCategoryToCategory(category)
	s.Require().Equal(expected, res)
}

func (s *ConserterSuite) TestRepoCategoryToCategoryReturnsFuel() {
	category := repoModel.CategoryFuel
	expected := model.CategoryFuel

	res := repoCategoryToCategory(category)
	s.Require().Equal(expected, res)
}

func (s *ConserterSuite) TestRepoCategoryToCategoryReturnsPorthole() {
	category := repoModel.CategoryPorthole
	expected := model.CategoryPorthole

	res := repoCategoryToCategory(category)
	s.Require().Equal(expected, res)
}

func (s *ConserterSuite) TestRepoCategoryToCategoryReturnsWing() {
	category := repoModel.CategoryWing
	expected := model.CategoryWing

	res := repoCategoryToCategory(category)
	s.Require().Equal(expected, res)
}

func (s *ConserterSuite) TestDimensionsToRepoDimensionsReturnsNil() {
	res := dimensionsToRepoDimensions(nil)
	s.Require().Nil(res)
}

func (s *ConserterSuite) TestDimensionsToRepoDimensionsSuccess() {
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
	s.Require().NotNil(res)
	s.Require().Equal(expected, res)
}

func (s *ConserterSuite) TestRepoDimensionsToDimensionsReturnsNil() {
	res := repoDimensionsToDimensions(nil)
	s.Require().Nil(res)
}

func (s *ConserterSuite) TestRepoDimensionsToDimensionsSuccess() {
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
	s.Require().NotNil(res)
	s.Require().Equal(expected, res)
}

func (s *ConserterSuite) TestManufacturerToRepoManufacturerReturnsNil() {
	res := manufacturerToRepoManufacturer(nil)
	s.Require().Nil(res)
}

func (s *ConserterSuite) TestManufacturerToRepoManufacturerSuccess() {
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
	s.Require().NotNil(res)
	s.Require().Equal(expected, res)
}

func (s *ConserterSuite) TestRepoManufacturerToManufacturerReturnsNil() {
	res := repoManufacturerToManufacturer(nil)
	s.Require().Nil(res)
}

func (s *ConserterSuite) TestRepoManufacturerToManufacturerSuccess() {
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
	s.Require().NotNil(res)
	s.Require().Equal(expected, res)
}

func (s *ConserterSuite) TestMetadataToRepoMetadataReturnsNil() {
	var m map[string]any
	res := metadataToRepoMetadata(m)
	s.Require().Nil(res)
}

func (s *ConserterSuite) TestMetadataToRepoMetadataSuccess() {
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
	s.Require().NotNil(res)
	s.Require().Equal(expected, m)
}

func (s *ConserterSuite) TestRepoMetadataToMetadataReturnsNil() {
	var m map[string]any
	res := repoMetadataToMetadata(m)
	s.Require().Nil(res)
}

func (s *ConserterSuite) TestRepoMetadataToMetadataSuccess() {
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
	s.Require().NotNil(res)
	s.Require().Equal(expected, m)
}
