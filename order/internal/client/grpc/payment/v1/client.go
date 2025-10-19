package v1

import (
	"time"

	def "github.com/max-kriv0s/go-microservices-edu/order/internal/client/grpc"
	paymentV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
)

var _ def.PaymentServiceClient = (*paymentServiceClient)(nil)

type paymentServiceClient struct {
	grpcTimeout time.Duration
	client      paymentV1.PaymentServiceClient
}

func NewPaymentServiceClient(conn *grpc.ClientConn) *paymentServiceClient {
	return &paymentServiceClient{
		grpcTimeout: 2 * time.Second,
		client:      paymentV1.NewPaymentServiceClient(conn),
	}
}
