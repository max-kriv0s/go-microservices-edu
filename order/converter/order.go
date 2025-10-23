package converter

import (
	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
	orderV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/openapi/order/v1"
)

func CreateOrderRequestToModel(req *orderV1.CreateOrderRequestDto) model.CreateOrderRequest {
	return model.CreateOrderRequest{
		UserUUID:  req.UserUUID,
		PartUuids: append([]string(nil), req.PartUuids...),
	}
}

func OrderToCreateOrderResponseDto(order model.Order) *orderV1.CreateOrderResponseDto {
	return &orderV1.CreateOrderResponseDto{
		OrderUUID:  order.OrderUUID,
		TotalPrice: float32(order.TotalPrice),
	}
}

func ApiPaymentMethodToPaymentMethod(method orderV1.PaymentMethod) model.PaymentMethod {
	switch method {
	case orderV1.PaymentMethodCARD:
		return model.PaymentMethodCard
	case orderV1.PaymentMethodSBP:
		return model.PaymentMethodSbp
	case orderV1.PaymentMethodCREDITCARD:
		return model.PaymentMethodCreditCard
	case orderV1.PaymentMethodINVESTORMONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnknown
	}
}

func OrderToGetResponseDto(order model.Order) (orderV1.APIV1OrdersOrderUUIDGetRes, error) {
	var transactionUUID orderV1.OptString
	if order.TransactionUUID != nil {
		transactionUUID = orderV1.NewOptString(*order.TransactionUUID)
	}

	var paymentMethod orderV1.OptPaymentMethod
	if order.PaymentMethod != nil {
		paymentMethod = orderV1.NewOptPaymentMethod(paymentMethodToApiPaymentMethod(*order.PaymentMethod))
	}

	status := statusToApiStatus(order.Status)

	return &orderV1.GetOrderResponseDto{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartsUUIDs,
		TotalPrice:      float32(order.TotalPrice),
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          status,
	}, nil
}

func paymentMethodToApiPaymentMethod(method model.PaymentMethod) orderV1.PaymentMethod {
	switch method {
	case model.PaymentMethodCard:
		return orderV1.PaymentMethodCARD
	case model.PaymentMethodSbp:
		return orderV1.PaymentMethodSBP
	case model.PaymentMethodCreditCard:
		return orderV1.PaymentMethodCREDITCARD
	case model.PaymentMethodInvestorMoney:
		return orderV1.PaymentMethodINVESTORMONEY
	default:
		return orderV1.PaymentMethodUNKNOWN
	}
}

func statusToApiStatus(status model.OrderStatus) orderV1.OrderStatus {
	switch status {
	case model.OrderStatusPaid:
		return orderV1.OrderStatusPAID
	case model.OrderStatusCancelled:
		return orderV1.OrderStatusCANCELLED
	default:
		return orderV1.OrderStatusPENDINGPAYMENT
	}
}
