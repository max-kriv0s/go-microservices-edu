package client

import (
	"context"

	"google.golang.org/grpc"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
	authV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/auth/v1"
)

type InventoryServiceClient interface {
	ListParts(ctx context.Context, partsUUIDs []string) ([]model.Part, error)
}

type PaymentServiceClient interface {
	PayOrder(ctx context.Context, order model.Order, paymentMethod model.PaymentMethod) (string, error)
}

type IamServiceClient interface {
	Login(ctx context.Context, req *authV1.LoginRequest, opts ...grpc.CallOption) (*authV1.LoginResponse, error)
	Whoami(ctx context.Context, req *authV1.WhoamiRequest, opts ...grpc.CallOption) (*authV1.WhoamiResponse, error)
}
