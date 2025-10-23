package v1

import (
	"context"
	"fmt"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
	paymentV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/payment/v1"
)

func (c *paymentServiceClient) PayOrder(ctx context.Context, order model.Order, paymentMethod model.PaymentMethod) (string, error) {
	grpcPaymentMethod, err := ConvertPaymentMethodToGRPC(paymentMethod)
	if err != nil {
		return "", err
	}

	paymentReq := &paymentV1.PayOrderRequest{
		OrderUuid:     order.OrderUUID,
		UserUuid:      order.UserUUID,
		PaymentMethod: grpcPaymentMethod,
	}

	res, err := c.client.PayOrder(ctx, paymentReq)
	if err != nil {
		return "", err
	}
	return res.TransactionUuid, nil
}

func ConvertPaymentMethodToGRPC(method model.PaymentMethod) (paymentV1.PaymentMethod, error) {
	switch method {
	case model.PaymentMethodCard:
		return paymentV1.PaymentMethod_CARD, nil
	case model.PaymentMethodSbp:
		return paymentV1.PaymentMethod_SBP, nil
	case model.PaymentMethodCreditCard:
		return paymentV1.PaymentMethod_CREDIT_CARD, nil
	case model.PaymentMethodInvestorMoney:
		return paymentV1.PaymentMethod_INVESTOR_MONEY, nil
	default:
		return paymentV1.PaymentMethod_UNKNOWN, fmt.Errorf("unknown payment method: %s", method)
	}
}
