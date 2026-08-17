package models

import "time"

// Константы для статусов листа ожидания
const (
	WaitlistStatusWaiting   = "waiting"   // Клиент просто ждет в очереди (по принципу FIFO)
	WaitlistStatusOffered   = "offered"   // Бот нашел окно и предложил его клиенту (таймаут 10 минут)
	WaitlistStatusFulfilled = "fulfilled" // Клиент успешно согласился и записался (создан Booking)
	WaitlistStatusExpired   = "expired"   // Клиент не ответил за 10 минут, очередь перешла к следующему
)

// WaitlistRequest описывает заявку клиента на поиск свободного окна
type WaitlistRequest struct {
	ID               int64     `json:"id"`
	MasterID         int64     `json:"master_id"`          // Внешний ключ на Master
	ClientTelegramID int64     `json:"client_telegram_id"` // Кто ждет
	DesiredDate      time.Time `json:"desired_date"`       // Желаемая дата (время обнулено до 00:00:00)
	Status           string    `json:"status"`             // waiting, offered, fulfilled, expired
	OfferedSlotID    *int64    `json:"offered_slot_id"`    // Указатель на ID слота (может быть nil/null)
	CreatedAt        time.Time `json:"created_at"`         // Приоритет очереди: кто раньше встал, тот первый в списке
	UpdatedAt        time.Time `json:"updated_at"`
}

