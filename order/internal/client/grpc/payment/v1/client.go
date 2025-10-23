package v1

import (
	def "github.com/max-kriv0s/go-microservices-edu/order/internal/client/grpc"
	paymentV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/payment/v1"
)

var _ def.PaymentServiceClient = (*paymentServiceClient)(nil)

type paymentServiceClient struct {
	client paymentV1.PaymentServiceClient
}

func NewPaymentServiceClient(client paymentV1.PaymentServiceClient) *paymentServiceClient {
	return &paymentServiceClient{
		client: client,
	}
}
