package converter

import (
	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
	"github.com/samber/lo"

	repoModel "github.com/max-kriv0s/go-microservices-edu/order/internal/repository/model"
)

func OrderToRepoModel(order model.Order) (repoModel.Order, error) {
	status, err := statusToRepoStatus(order.Status)
	if err != nil {
		return repoModel.Order{}, err
	}

	return repoModel.Order{
		OrderUUID:     order.OrderUUID,
		UserUUID:      order.UserUUID,
		PartsUUIDs:    append([]string(nil), order.PartsUUIDs...),
		TotalPrice:    order.TotalPrice,
		PaymentMethod: paymentMethodToRepoPaymentMethod(order.PaymentMethod),
		Status:        status,
	}, nil
}

func OrderToModel(repoOrder repoModel.Order) (model.Order, error) {
	status, err := repoStatusToStatus(repoOrder.Status)
	if err != nil {
		return model.Order{}, err
	}

	return model.Order{
		OrderUUID:     repoOrder.OrderUUID,
		UserUUID:      repoOrder.UserUUID,
		PartsUUIDs:    append([]string(nil), repoOrder.PartsUUIDs...),
		TotalPrice:    repoOrder.TotalPrice,
		PaymentMethod: repoPaymentMethodToPaymentMethod(repoOrder.PaymentMethod),
		Status:        status,
	}, nil
}

func paymentMethodToRepoPaymentMethod(method *model.PaymentMethod) *repoModel.PaymentMethod {
	if method == nil {
		return nil
	}

	switch *method {
	case model.PaymentMethodCard:
		return lo.ToPtr(repoModel.PaymentMethodCard)
	case model.PaymentMethodSbp:
		return lo.ToPtr(repoModel.PaymentMethodSbp)
	case model.PaymentMethodCreditCard:
		return lo.ToPtr(repoModel.PaymentMethodCreditCard)
	case model.PaymentMethodInvestorMoney:
		return lo.ToPtr(repoModel.PaymentMethodInvestorMoney)
	default:
		return lo.ToPtr(repoModel.PaymentMethodUnknown)
	}
}

func repoPaymentMethodToPaymentMethod(method *repoModel.PaymentMethod) *model.PaymentMethod {
	if method == nil {
		return nil
	}

	switch *method {
	case repoModel.PaymentMethodCard:
		return lo.ToPtr(model.PaymentMethodCard)
	case repoModel.PaymentMethodSbp:
		return lo.ToPtr(model.PaymentMethodSbp)
	case repoModel.PaymentMethodCreditCard:
		return lo.ToPtr(model.PaymentMethodCreditCard)
	case repoModel.PaymentMethodInvestorMoney:
		return lo.ToPtr(model.PaymentMethodInvestorMoney)
	default:
		return lo.ToPtr(model.PaymentMethodUnknown)
	}
}

func statusToRepoStatus(status model.OrderStatus) (repoModel.OrderStatus, error) {
	switch status {
	case model.OrderStatusPendingPayment:
		return repoModel.OrderStatusPendingPayment, nil
	case model.OrderStatusPaid:
		return repoModel.OrderStatusPaid, nil
	case model.OrderStatusCancelled:
		return repoModel.OrderStatusCancelled, nil
	default:
		return repoModel.OrderStatusPendingPayment, model.ErrUnknownRepoOrderStatus
	}
}

func repoStatusToStatus(status repoModel.OrderStatus) (model.OrderStatus, error) {
	switch status {
	case repoModel.OrderStatusPendingPayment:
		return model.OrderStatusPendingPayment, nil
	case repoModel.OrderStatusPaid:
		return model.OrderStatusPaid, nil
	case repoModel.OrderStatusCancelled:
		return model.OrderStatusCancelled, nil
	default:
		return model.OrderStatusPendingPayment, model.ErrUnknownRepoOrderStatus
	}
}
