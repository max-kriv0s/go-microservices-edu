package app

import (
	"context"

	paymentV1API "github.com/max-kriv0s/go-microservices-edu/payment/internal/api/payment/v1"
	"github.com/max-kriv0s/go-microservices-edu/payment/internal/service"
	paymentService "github.com/max-kriv0s/go-microservices-edu/payment/internal/service/payment"
	paymentV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/payment/v1"
)

type diContaier struct {
	paymentV1API   paymentV1.PaymentServiceServer
	paymentService service.PaymentService
}

func NewDiContainer() *diContaier {
	return &diContaier{}
}

func (d *diContaier) PaymentV1API(ctx context.Context) paymentV1.PaymentServiceServer {
	if d.paymentV1API == nil {
		d.paymentV1API = paymentV1API.NewAPI(d.PaymentService(ctx))
	}

	return d.paymentV1API
}

func (d *diContaier) PaymentService(ctx context.Context) service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = paymentService.NewService()
	}

	return d.paymentService
}
