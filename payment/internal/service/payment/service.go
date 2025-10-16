package payment

import def "github.com/max-kriv0s/go-microservices-edu/payment/internal/service"

var _ def.PaymentService = (*service)(nil)

type service struct{}

func NewService() *service {
	return &service{}
}
