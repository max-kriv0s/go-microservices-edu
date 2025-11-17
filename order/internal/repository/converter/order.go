package converter

import (
	"github.com/samber/lo"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
	repoModel "github.com/max-kriv0s/go-microservices-edu/order/internal/repository/model"
)

func OrderToRepoModel(order model.Order) repoModel.Order {
	return repoModel.Order{
		OrderUUID:     order.OrderUUID,
		UserUUID:      order.UserUUID,
		PartsUUIDs:    append([]string(nil), order.PartsUUIDs...),
		TotalPrice:    order.TotalPrice,
		PaymentMethod: PaymentMethodToRepoPaymentMethod(order.PaymentMethod),
		Status:        StatusToRepoStatus(order.Status),
	}
}

func OrderToModel(repoOrder repoModel.Order) model.Order {
	return model.Order{
		OrderUUID:     repoOrder.OrderUUID,
		UserUUID:      repoOrder.UserUUID,
		PartsUUIDs:    append([]string(nil), repoOrder.PartsUUIDs...),
		TotalPrice:    repoOrder.TotalPrice,
		PaymentMethod: repoPaymentMethodToPaymentMethod(repoOrder.PaymentMethod),
		Status:        repoStatusToStatus(repoOrder.Status),
	}
}

func PaymentMethodToRepoPaymentMethod(method *model.PaymentMethod) *repoModel.PaymentMethod {
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

func StatusToRepoStatus(status model.OrderStatus) repoModel.OrderStatus {
	switch status {
	case model.OrderStatusPaid:
		return repoModel.OrderStatusPaid
	case model.OrderStatusCancelled:
		return repoModel.OrderStatusCancelled
	case model.OrderStatusAssembled:
		return repoModel.OrderStatusAssembled
	default:
		return repoModel.OrderStatusPendingPayment
	}
}

func repoStatusToStatus(status repoModel.OrderStatus) model.OrderStatus {
	switch status {
	case repoModel.OrderStatusPaid:
		return model.OrderStatusPaid
	case repoModel.OrderStatusCancelled:
		return model.OrderStatusCancelled
	case repoModel.OrderStatusAssembled:
		return model.OrderStatusAssembled
	default:
		return model.OrderStatusPendingPayment
	}
}
