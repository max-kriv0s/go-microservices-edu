package model

type Order struct {
	OrderUUID       string `db:"id"`
	UserUUID        string `db:"user_id"`
	PartsUUIDs      []string
	TotalPrice      float64        `db:"total_price"`
	TransactionUUID *string        `db:"transaction_uuid"`
	PaymentMethod   *PaymentMethod `db:"payment_method"`
	Status          OrderStatus    `db:"status"`
}

type PaymentMethod string

const (
	PaymentMethodUnknown       PaymentMethod = "UNKNOWN"
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSbp           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)
