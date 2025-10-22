package model

import "time"

type Part struct {
	// Уникальный идентификатор детали
	Uuid string
	// Название детали
	Name string
	// Описание детали
	Description string
	// Цена за единицу
	Price float64
	// Количество на складе
	StockQuantity int64
	// Категория
	Category Category
	// Размеры детали
	Dimensions *Dimensions
	// Информация о производителе
	Manufacturer *Manufacturer
	// Теги для быстрого поиска
	Tags []string
	// Гибкие метаданные
	Metadata map[string]any
	// Дата создания
	CreatedAt *time.Time
	// Дата обновления
	UpdatedAt *time.Time
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
	Length float64
	// Ширина в см
	Width float64
	// Высота в см
	Height float64
	// Вес в кг
	Weight float64
}

type Manufacturer struct {
	// Название
	Name string
	// Страна производства
	Country string
	// Сайт производителя
	Website string
}

type PartsFilter struct {
	// Список UUID'ов. Пусто — не фильтруем по UUID
	Uuids []string
	// Список имён. Пусто — не фильтруем по имени
	Names []string
	// Список категорий. Пусто — не фильтруем по категории
	Categories []Category
	// Список стран производителей. Пусто — не фильтруем по стране
	ManufacturerCountries []string
	// Список тегов. Пусто — не фильтруем по тегам
	Tags []string
}
