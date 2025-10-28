package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Part struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	// Уникальный идентификатор детали
	Uuid string `bson:"uuid"`
	// Название детали
	Name string `bson:"name"`
	// Описание детали
	Description string `bson:"description"`
	// Цена за единицу
	Price float64 `bson:"price"`
	// Количество на складе
	StockQuantity int64 `bson:"stock_quantity"`
	// Категория
	Category Category `bson:"category"`
	// Размеры детали
	Dimensions *Dimensions `bson:"dimensions,omitempty"`
	// Информация о производителе
	Manufacturer *Manufacturer `bson:"manufacturer,omitempty"`
	// Теги для быстрого поиска
	Tags []string `bson:"tags,omitempty"`
	// Гибкие метаданные
	Metadata map[string]any `bson:"metadate,omitempty"`
	// Дата создания
	CreatedAt *time.Time `bson:"created_at,omitempty"`
	// Дата обновления
	UpdatedAt *time.Time `bson:"updated_at,omitempty"`
}

type Category int32

const (
	// Неизвестная категория
	CategoryUnknown Category = iota
	// Двигатель
	CategoryEngine
	// Топливо
	CategoryFuel
	// Иллюминатор
	CategoryPorthole
	// Крыло
	CategoryWing
)

type Dimensions struct {
	// Длина в см
	Length float64 `bson:"length"`
	// Ширина в см
	Width float64 `bson:"width"`
	// Высота в см
	Height float64 `bson:"height"`
	// Вес в кг
	Weight float64 `bson:"weight"`
}

type Manufacturer struct {
	// Название
	Name string `bson:"name"`
	// Страна производства
	Country string `bson:"country"`
	// Сайт производителя
	Website string `bson:"website"`
}
