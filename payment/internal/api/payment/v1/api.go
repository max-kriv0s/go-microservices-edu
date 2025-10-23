package v1

import (
	"github.com/max-kriv0s/go-microservices-edu/payment/internal/service"
	paymentV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/payment/v1"
)

type api struct {
	paymentV1.UnimplementedPaymentServiceServer

	paymentService service.PaymentService
}

func NewAPI(paymentService service.PaymentService) *api {
	return &api{
		paymentService: paymentService,
	}
}
