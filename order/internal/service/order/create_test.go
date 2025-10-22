package order

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
)

func fakePart() model.Part {
	part := model.Part{
		Uuid:  gofakeit.UUID(),
		Price: 100,
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
