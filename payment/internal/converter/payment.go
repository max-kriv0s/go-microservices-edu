package converter

import (
	"github.com/max-kriv0s/go-microservices-edu/payment/internal/model"
	paymentV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/payment/v1"
)

func PayDtoToModel(dto *paymentV1.PayOrderRequest) model.PayOrder {
	return model.PayOrder{
		OrderUuid:     dto.OrderUuid,
		UserUuid:      dto.UserUuid,
		PaymentMethod: toPaymentMethod(dto.PaymentMethod),
	}
}

func toPaymentMethod(pm paymentV1.PaymentMethod) model.PaymentMethod {
	switch pm {
	case paymentV1.PaymentMethod_CARD:
		return model.PaymentCard
	case paymentV1.PaymentMethod_SBP:
		return model.PaymentSBP
	case paymentV1.PaymentMethod_CREDIT_CARD:
		return model.PaymentCreditCard
	case paymentV1.PaymentMethod_INVESTOR_MONEY:
		return model.PaymentInvestorMoney
	default:
		return model.PaymentUnknown
	}
}
