package models

import "time"

// Константы для статусов записи
const (
	BookingStatusActive            = "active"              // Создана, ждет времени сеанса
	BookingStatusConfirmed         = "confirmed"           // Подтверждена клиентом (например, за 2 часа до начала)
	BookingStatusCancelledNoShow   = "cancelled_no_show"   // Авто-отмена (клиент не подтвердил вовремя или не пришел)
	BookingStatusCancelledByClient = "cancelled_by_client" // Клиент отменил запись вручную
	BookingStatusCompleted         = "completed"           // Сеанс успешно завершен, мастер выполнил работу
)

// Booking описывает факт записи клиента на конкретную услугу и время
type Booking struct {
	ID               int64     `json:"id"`
	SlotID           int64     `json:"slot_id"`            // Внешний ключ на Slot (Связь 1-к-1)
	ServiceID        int64     `json:"service_id"`         // Внешний ключ на Service
	ClientTelegramID int64     `json:"client_telegram_id"` // ID клиента в Telegram
	ClientName       string    `json:"client_name"`        // Имя клиента из Telegram
	PriceLocked      int       `json:"price_locked"`       // Фиксация цены на момент записи (в целых числах)
	Status           string    `json:"status"`             // active, confirmed, cancelled_..., completed
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Тестовая функция для демонстрации созданной записи
func GetTestBooking(slotID, serviceID int64) Booking {
	return Booking{
		ID:               1,
		SlotID:           slotID,
		ServiceID:        serviceID,
		ClientTelegramID: 987654321,
		ClientName:       "Алексей",
		PriceLocked:      1500, 
		Status:           BookingStatusActive,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}
