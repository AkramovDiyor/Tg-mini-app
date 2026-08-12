package models

import "time"

// Service описывает услугу, которую предоставляет мастер
type Service struct {
	ID          int64     `json:"id"`
	MasterID    int64     `json:"master_id"`    // Внешний ключ на Master
	Name        string    `json:"name"`         // "Мужская стрижка"
	DurationMin int       `json:"duration_min"` // Длительность в минутах (например, 45)
	Price       int       `json:"price"`        // Цена в целых числах (например, 1500)
	Description string    `json:"description"`  // Опциональное описание услуги
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Тестовая функция для создания списка услуг 
func GetTestServices(masterID int64) []Service {
	return []Service{
		{
			ID:          1,
			MasterID:    masterID,
			Name:        "Мужская стрижка",
			DurationMin: 45,
			Price:       1500,
			Description: "Классическая стрижка ножницами и машинкой",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          2,
			MasterID:    masterID,
			Name:        "Стрижка машинкой",
			DurationMin: 20,
			Price:       800,
			Description: "Быстрая стрижка под одну или несколько насадок",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
}
