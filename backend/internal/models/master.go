package models

import "time"

// WorkHours описывает график работы мастера из UI
type WorkHours struct {
	WorkingDays []string `json:"working_days"` // ["Пн", "Вт", "Ср", "Чт", "Пт"]
	StartDay    string   `json:"start_day"`    // "09:00"
	EndDay      string   `json:"end_day"`      // "20:00"
	BreakStart  string   `json:"break_start"`  // "13:00"
	BreakEnd    string   `json:"break_end"`    // "14:00"
}

// MasterSettings описывает правила записи и автоматизации
type MasterSettings struct {
	AutoCancel        bool `json:"auto_cancel"`         // Авто-отмена без подтверждения
	CancelBeforeHours int  `json:"cancel_before_hours"` // 1, 2, 4 часа
	UseWaitlist       bool `json:"use_waitlist"`        // Предлагать окно в лист ожидания
	SlotStepMin       int  `json:"slot_step_min"`
}

// Master — основное ядро системы
type Master struct {
	ID         int64          `json:"id"`
	TelegramID int64          `json:"telegram_id"` // Уникальный индекс в БД
	Name       string         `json:"name"`        // "Педро Барбер"
	Bio        string         `json:"bio"`         // "Барбер · 6 лет опыта"
	Address    string         `json:"address"`     // "ул. Центральная, 1"
	InviteLink string         `json:"invite_link"` // "2hq0og3k"
	WorkHours  WorkHours      `json:"work_hours"`  // В БД будет сохраняться как JSONB
	Settings   MasterSettings `json:"settings"`    // В БД будет сохраняться как JSONB
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

// Тестовая функция для проверки создания объекта
func Get() Master {
	return Master{
		ID:         1,
		TelegramID: 12345678910,
		Name:       "Педро Барбер",
		Bio:        "Барбер · 6 лет опыта",
		Address:    "ул. Центральная, 1",
		InviteLink: "2hq0og3k",
		WorkHours: WorkHours{
			WorkingDays: []string{"Пн", "Вт", "Ср", "Чт", "Пт"},
			StartDay:    "09:00",
			EndDay:      "20:00",
			BreakStart:  "13:00",
			BreakEnd:    "14:00",
		},
		Settings: MasterSettings{
			AutoCancel:        true,
			CancelBeforeHours: 2,
			UseWaitlist:       true,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
