package v1

import (
	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/model"
	inventoryV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/inventory/v1"
)

func (s *APISuite) TestGetPartSuccess() {
	var (
		uuid = gofakeit.UUID()

		part = fakePart()

		req = &inventoryV1.GetPartRequest{
			Uuid: uuid,
		}
	)

	val, ok := part.Metadata["string"].(string)
	s.Require().True(ok, "expected metadata['string'] to be string, got %T", part.Metadata["string"])

	expectedMetadata := make(map[string]*inventoryV1.Value, len(part.Metadata))
	expectedMetadata["string"] = &inventoryV1.Value{
		Value: &inventoryV1.Value_StringValue{
			StringValue: val,
		},
	}

	var createdAt *timestamppb.Timestamp
	if part.CreatedAt != nil {
		createdAt = timestamppb.New(*part.CreatedAt)
	}

	var updatedAt *timestamppb.Timestamp
	if part.UpdatedAt != nil {
		updatedAt = timestamppb.New(*part.UpdatedAt)
	}

	expected := &inventoryV1.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      inventoryV1.Category_ENGINE,
		Dimensions: &inventoryV1.Dimensions{
			Length: part.Dimensions.Length,
			Width:  part.Dimensions.Width,
			Height: part.Dimensions.Height,
			Weight: part.Dimensions.Weight,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    part.Manufacturer.Name,
			Country: part.Manufacturer.Country,
			Website: part.Manufacturer.Website,
		},
		Tags:      append([]string(nil), part.Tags...),
		Metadata:  expectedMetadata,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	s.inventoryService.On("GetPart", s.Ctx(), uuid).Return(part, nil)

	res, err := s.api.GetPart(s.Ctx(), req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expected, res.GetPart())
}

func (s *APISuite) TestGetPartNotFoundError() {
	var (
		serviceErr = model.ErrPartNotFound
		uuid       = gofakeit.UUID()

		req = &inventoryV1.GetPartRequest{
			Uuid: uuid,
		}
	)

	s.inventoryService.On("GetPart", s.Ctx(), uuid).Return(model.Part{}, serviceErr)

	res, err := s.api.GetPart(s.Ctx(), req)
	s.Require().Nil(res)
	s.Require().Error(err)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.NotFound, st.Code())
}

func (s *APISuite) TestGetPartInternalError() {
	var (
		serviceErr = gofakeit.Error()
		uuid       = gofakeit.UUID()

		req = &inventoryV1.GetPartRequest{
			Uuid: uuid,
		}
	)

	s.inventoryService.On("GetPart", s.Ctx(), uuid).Return(model.Part{}, serviceErr)

	res, err := s.api.GetPart(s.Ctx(), req)
	s.Require().Nil(res)
	s.Require().Error(err)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.Internal, st.Code())
}
