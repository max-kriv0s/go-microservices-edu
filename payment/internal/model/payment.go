package model

type PayOrder struct {
	OrderUuid     string
	UserUuid      string
	PaymentMethod PaymentMethod
}

type PaymentMethod int

const (
	PaymentUnknown PaymentMethod = iota
	PaymentCard
	PaymentSBP
	PaymentCreditCard
	PaymentInvestorMoney
)
