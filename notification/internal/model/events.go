package model

// OrderPaidEvent - событие оплаты заказа
type OrderPaidEvent struct {
	EventUUID       string // Уникальный идентификатор события (для идемпотентности)
	OrderUUID       string // Идентификатор оплаченного заказа
	UserUUID        string // Идентификатор пользователя
	PaymentMethod   string // Способ оплаты (строкой, значение из `PaymentMethod`)
	TransactionUUID string // Идентификатор транзакции, сгенерированный в результате оплаты
}

// ShipAssembledEvent - событие сборки заказа
type ShipAssembledEvent struct {
	EventUUID    string // Уникальный идентификатор события (для идемпотентности)
	OrderUUID    string // Идентификатор собранного заказа
	UserUUID     string // Идентификатор пользователя
	BuildTimeSec int64  // Время (в секундах), потраченное на сборку корабля
}
