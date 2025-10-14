package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/payment/v1"
)

const (
	inventoryUrl      = "localhost:50051"
	paymentUrl        = "localhost:50052"
	grpcTimeout       = 2 * time.Second
	serverTimeout     = 10 * time.Second
	serverHost        = "localhost"
	serverPort        = "8080"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

var ErrOrderNotFound = errors.New("order not found")

type Order struct {
	OrderUUID       string
	UserUUID        string
	PartsUUIDs      []string
	TotalPrice      float64
	TransactionUUID *string
	PaymentMethod   *orderV1.PaymentMethod
	Status          orderV1.OrderStatus
}

func NewInventoryServiceClient(conn *grpc.ClientConn) inventoryV1.InventoryServiceClient {
	return inventoryV1.NewInventoryServiceClient(conn)
}

func NewPaymentServiceClient(conn *grpc.ClientConn) paymentV1.PaymentServiceClient {
	return paymentV1.NewPaymentServiceClient(conn)
}

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*Order
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*Order),
	}
}

func (s *OrderStorage) AddOrder(order *Order) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[order.OrderUUID] = order
}

func (s *OrderStorage) GetOrder(uuid string) (*Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[uuid]
	if !ok {
		return nil, ErrOrderNotFound
	}
	return order, nil
}

func (s *OrderStorage) UpdateOrder(uuid string, order *Order) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[uuid] = order
}

type OrderHandler struct {
	storage *OrderStorage

	inventoryClient inventoryV1.InventoryServiceClient
	paymentClient   paymentV1.PaymentServiceClient
}

func NewOrderHandler(storage *OrderStorage, inventoryClient inventoryV1.InventoryServiceClient, paymentClient paymentV1.PaymentServiceClient) *OrderHandler {
	return &OrderHandler{
		storage:         storage,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequestDto) (orderV1.CreateOrderRes, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	inventoryReq := &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Uuids: req.PartUuids,
		},
	}
	listParts, err := h.inventoryClient.ListParts(ctxWithTimeout, inventoryReq)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, nil
	}
	if len(listParts.Parts) != len(req.PartUuids) {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "All the details were not found",
		}, nil
	}

	var totalPrice float64
	for _, part := range listParts.Parts {
		totalPrice += part.Price
	}

	createdOrder := &Order{
		OrderUUID:  uuid.NewString(),
		UserUUID:   req.UserUUID,
		PartsUUIDs: req.PartUuids,
		TotalPrice: totalPrice,
		Status:     orderV1.OrderStatusPENDINGPAYMENT,
	}
	h.storage.AddOrder(createdOrder)

	return &orderV1.CreateOrderResponseDto{
		OrderUUID:  createdOrder.OrderUUID,
		TotalPrice: float32(createdOrder.TotalPrice),
	}, nil
}

func (h *OrderHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequestDto, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	order, err := h.storage.GetOrder(params.OrderUUID.String())
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "order not found",
			}, nil
		} else {
			return &orderV1.InternalServerError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}, nil
		}
	}

	if order.Status != orderV1.OrderStatusPENDINGPAYMENT {
		return &orderV1.ConflictError{
			Code:    http.StatusConflict,
			Message: fmt.Sprintf("You can't pay an order. Order status %s", order.Status),
		}, nil
	}

	grpcPaymentMethod, err := ConvertPaymentMethodToGRPC(req.PaymentMethod)
	if err != nil {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "Not correct payment method",
		}, nil
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	paymentReq := &paymentV1.PayOrderRequest{
		OrderUuid:     order.OrderUUID,
		UserUuid:      order.UserUUID,
		PaymentMethod: grpcPaymentMethod,
	}
	res, err := h.paymentClient.PayOrder(ctxWithTimeout, paymentReq)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, nil
	}

	order.Status = orderV1.OrderStatusPAID
	order.PaymentMethod = &req.PaymentMethod
	order.TransactionUUID = &res.TransactionUuid

	h.storage.UpdateOrder(order.OrderUUID, order)

	return &orderV1.PayOrderResponseDto{
		TransactionUUID: res.TransactionUuid,
	}, nil
}

func (h *OrderHandler) APIV1OrdersOrderUUIDGet(_ context.Context, params orderV1.APIV1OrdersOrderUUIDGetParams) (orderV1.APIV1OrdersOrderUUIDGetRes, error) {
	order, err := h.storage.GetOrder(params.OrderUUID.String())
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "order not found",
			}, nil
		} else {
			return &orderV1.InternalServerError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}, nil
		}
	}

	var transactionUUID orderV1.OptString
	if order.TransactionUUID != nil {
		transactionUUID = orderV1.NewOptString(*order.TransactionUUID)
	}

	var paymentMethod orderV1.OptPaymentMethod
	if order.PaymentMethod != nil {
		paymentMethod = orderV1.NewOptPaymentMethod(*order.PaymentMethod)
	}

	return &orderV1.GetOrderResponseDto{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartsUUIDs,
		TotalPrice:      float32(order.TotalPrice),
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          order.Status,
	}, nil
}

func (h *OrderHandler) APIV1OrdersOrderUUIDCancelPost(ctx context.Context, params orderV1.APIV1OrdersOrderUUIDCancelPostParams) (orderV1.APIV1OrdersOrderUUIDCancelPostRes, error) {
	order, err := h.storage.GetOrder(params.OrderUUID.String())
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "order not found",
			}, nil
		} else {
			return &orderV1.InternalServerError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}, nil
		}
	}

	if order.Status != orderV1.OrderStatusPENDINGPAYMENT {
		return &orderV1.ConflictError{
			Code:    http.StatusConflict,
			Message: fmt.Sprintf("You can't cancel an order. Order status %s", order.Status),
		}, nil
	}

	order.Status = orderV1.OrderStatusCANCELLED

	h.storage.UpdateOrder(order.OrderUUID, order)

	return &orderV1.APIV1OrdersOrderUUIDCancelPostNoContent{}, nil
}

func (h *OrderHandler) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}

func ConvertPaymentMethodToGRPC(method orderV1.PaymentMethod) (paymentV1.PaymentMethod, error) {
	switch method {
	case orderV1.PaymentMethodCARD:
		return paymentV1.PaymentMethod_CARD, nil
	case orderV1.PaymentMethodSBP:
		return paymentV1.PaymentMethod_SBP, nil
	case orderV1.PaymentMethodCREDITCARD:
		return paymentV1.PaymentMethod_CREDIT_CARD, nil
	case orderV1.PaymentMethodINVESTORMONEY:
		return paymentV1.PaymentMethod_INVESTOR_MONEY, nil
	default:
		return paymentV1.PaymentMethod_UNKNOWN, fmt.Errorf("unknown payment method: %s", method)
	}
}

func main() {
	storage := NewOrderStorage()

	inventoryConn, err := grpc.NewClient(
		inventoryUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect to inventory service: %v", err)
		return
	}
	defer func() {
		if err := inventoryConn.Close(); err != nil {
			log.Printf("failed to close inventory connection: %v", err)
		}
	}()

	paymentConn, err := grpc.NewClient(
		paymentUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect to payment service: %v", err)
		return
	}
	defer func() {
		if err := paymentConn.Close(); err != nil {
			log.Printf("failed to close payment connection: %v", err)
		}
	}()

	inventoryClient := NewInventoryServiceClient(inventoryConn)
	paymentClient := NewPaymentServiceClient(paymentConn)

	orderHandler := NewOrderHandler(storage, inventoryClient, paymentClient)

	orderServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		log.Printf("ошибка создания сервера OpenAPI: %v", err)
		return
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(serverTimeout))

	r.Mount("/", orderServer)

	serverAddr := net.JoinHostPort(serverHost, serverPort)
	server := &http.Server{
		Addr:              serverAddr,
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на %s\n", serverAddr)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
