package models

import "time"

// Константы для статусов слота (вместо сырых строк)
const (
	SlotStatusFree    = "free"    // Свободен для записи
	SlotStatusBooked  = "booked"  // Забронирован и подтвержден
	SlotStatusPending = "pending" // Клиент выбрал время, но еще не подтвердил (бронь держится Х минут)
	SlotStatusBlocked = "blocked" // Перерыв, обед или ручная блокировка мастером
)

// Slot описывает конкретный отрезок времени в расписании мастера
type Slot struct {
	ID        int64     `json:"id"`
	MasterID  int64     `json:"master_id"`  // Внешний ключ на Master
	StartTime time.Time `json:"start_time"` // Когда начинается слот (например, 2026-11-15 15:00:00 +03:00)
	EndTime   time.Time `json:"end_time"`   // Когда заканчивается (StartTime + DurationMin)
	Status    string    `json:"status"`     // free, booked, pending, blocked
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Тестовая функция для демонстрации сетки слотов
func GetTestSlots(masterID int64) []Slot {
	baseTime := time.Date(2026, 11, 15, 15, 0, 0, 0, time.Local) // 15 ноября 15:00

	return []Slot{
		{
			ID:        1,
			MasterID:  masterID,
			StartTime: baseTime,
			EndTime:   baseTime.Add(45 * time.Minute), // 15:45
			Status:    SlotStatusFree,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			MasterID:  masterID,
			StartTime: baseTime.Add(45 * time.Minute), // 15:45
			EndTime:   baseTime.Add(65 * time.Minute), // 16:05 (например, под быструю стрижку 20 мин)
			Status:    SlotStatusPending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}
