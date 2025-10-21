package order

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
)

func fakePart() model.Part {
	now := time.Now()

	category := model.CategoryEngine

	metadata := make(map[string]*model.Value)

	key := gofakeit.Word()
	metadata[key] = &model.Value{String: lo.ToPtr(gofakeit.Word())}

	tagsCount := gofakeit.Number(1, 5)
	tags := make([]string, tagsCount)
	for i := 0; i < tagsCount; i++ {
		tags[i] = gofakeit.Word()
	}

	part := model.Part{
		Uuid:          gofakeit.UUID(),
		Name:          gofakeit.ProductName(),
		Description:   gofakeit.Sentence(),
		Price:         100,
		StockQuantity: int64(gofakeit.IntRange(0, 100)),
		Category:      category,
		Dimensions: &model.Dimensions{
			Length: gofakeit.Float64Range(1, 100),
			Width:  gofakeit.Float64Range(1, 100),
			Height: gofakeit.Float64Range(1, 100),
			Weight: gofakeit.Float64Range(0.1, 50),
		},
		Manufacturer: &model.Manufacturer{
			Name:    gofakeit.Company(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		},
		Tags:      tags,
		Metadata:  metadata,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	return part
}

func (s *ServiceSuite) TestCreateOrderSuccess() {
	var (
		part1 = fakePart()
		part2 = fakePart()

		partUUIDs = []string{part1.Uuid, part2.Uuid}
		userUUID  = gofakeit.UUID()
		orderUUID = gofakeit.UUID()

		listParts = []model.Part{part1, part2}
		order     = model.Order{
			OrderUUID:  "",
			UserUUID:   userUUID,
			PartsUUIDs: partUUIDs,
			TotalPrice: 200,
			Status:     model.OrderStatusPendingPayment,
		}
		req = model.CreateOrderRequest{
			UserUUID:  userUUID,
			PartUuids: partUUIDs,
		}
	)

	s.inventoryServiceClient.On("ListParts", s.Ctx(), partUUIDs).Return(listParts, nil)

	createdOrder := order
	createdOrder.OrderUUID = orderUUID
	s.orderRepository.On("Create", s.Ctx(), order).Return(orderUUID, nil)

	res, err := s.service.CreateOrder(s.Ctx(), req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(createdOrder, res)
}

func (s *ServiceSuite) TestCreateOrderClientError() {
	var (
		clietnErr   = gofakeit.Error()
		expectedErr = model.ErrInternalServer
		part1       = fakePart()
		part2       = fakePart()

		partUUIDs = []string{part1.Uuid, part2.Uuid}
		userUUID  = gofakeit.UUID()
		req       = model.CreateOrderRequest{
			UserUUID:  userUUID,
			PartUuids: partUUIDs,
		}
	)

	s.inventoryServiceClient.On("ListParts", s.Ctx(), partUUIDs).Return(nil, clietnErr)

	res, err := s.service.CreateOrder(s.Ctx(), req)
	s.Require().Empty(res)
	s.Require().Error(err)
	s.Require().ErrorIs(err, expectedErr)
}

func (s *ServiceSuite) TestCreateOrderClientBadRequestError() {
	var (
		expectedErr = model.ErrBadRequest
		part1       = fakePart()
		part2       = fakePart()

		partUUIDs = []string{part1.Uuid, part2.Uuid}
		userUUID  = gofakeit.UUID()
		req       = model.CreateOrderRequest{
			UserUUID:  userUUID,
			PartUuids: partUUIDs,
		}

		listParts = []model.Part{part1}
	)

	s.inventoryServiceClient.On("ListParts", s.Ctx(), partUUIDs).Return(listParts, nil)

	res, err := s.service.CreateOrder(s.Ctx(), req)
	s.Require().Empty(res)
	s.Require().Error(err)
	s.Require().ErrorIs(err, expectedErr)
}

func (s *ServiceSuite) TestCreateOrderInternalError() {
	var (
		repoErr     = gofakeit.Error()
		expectedErr = model.ErrInternalServer

		part1 = fakePart()
		part2 = fakePart()

		partUUIDs = []string{part1.Uuid, part2.Uuid}
		userUUID  = gofakeit.UUID()
		orderUUID = gofakeit.UUID()

		listParts = []model.Part{part1, part2}
		order     = model.Order{
			OrderUUID:  "",
			UserUUID:   userUUID,
			PartsUUIDs: partUUIDs,
			TotalPrice: 200,
			Status:     model.OrderStatusPendingPayment,
		}
		req = model.CreateOrderRequest{
			UserUUID:  userUUID,
			PartUuids: partUUIDs,
		}
	)

	s.inventoryServiceClient.On("ListParts", s.Ctx(), partUUIDs).Return(listParts, nil)

	createdOrder := order
	createdOrder.OrderUUID = orderUUID
	s.orderRepository.On("Create", s.Ctx(), order).Return("", repoErr)

	res, err := s.service.CreateOrder(s.Ctx(), req)
	s.Require().Empty(res)
	s.Require().Error(err)
	s.Require().ErrorIs(err, expectedErr)
}
