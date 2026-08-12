package models

import "time"

// Client описывает пользователя (клиента), который записывается на услуги через Telegram
type Client struct {
	TelegramID   int64     `json:"telegram_id"` // Идентификатор Telegram (PRIMARY KEY)
	Username     string    `json:"username"`    // @никнейм клиента
	FirstName    string    `json:"first_name"`   // Имя из Telegram
	LastName     string    `json:"last_name"`    // Фамилия из Telegram
	StrikesCount int       `json:"strikes_count"` // Количество нарушений (не подтвердил запись вовремя)
	IsBanned     bool      `json:"is_banned"`    // Флаг блокировки (при StrikesCount >= 3)
	CreatedAt    time.Time `json:"created_at"`
}

// Тестовая функция для демонстрации структуры клиента
func GetTestClient() Client {
	return Client{
		TelegramID:   987654321,
		Username:     "alex_barber_client",
		FirstName:    "Алексей",
		LastName:     "Петров",
		StrikesCount: 0,
		IsBanned:     false,
		CreatedAt:    time.Now(),
	}
}
