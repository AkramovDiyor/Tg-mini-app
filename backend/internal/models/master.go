package models

import "time"

// WorkHours описывает график работы мастера из UI
type WorkHours struct {
    WorkDays   []bool `json:"work_days"`   // [Пн, Вт, Ср, Чт, Пт, Сб, Вс]
    StartTime  string `json:"start_time"`
    EndTime    string `json:"end_time"`
    LunchStart string `json:"lunch_start"`
    LunchEnd   string `json:"lunch_end"`
}

// MasterSettings описывает правила записи и автоматизации
type MasterSettings struct {
    AutoCancel    bool   `json:"auto_cancel"`
    CancelHours   string `json:"cancel_hours"`
    OfferWaitlist bool   `json:"offer_waitlist"`
}

// Master — основное ядро системы
type Master struct {
    ID         int64          `json:"id"`
    TelegramID int64          `json:"telegram_id"`
    Name       string         `json:"name"`
    Bio        string         `json:"bio"`
    Address    string         `json:"address"`
    InviteLink string         `json:"invite_link"`
    WorkHours  WorkHours      `json:"work_hours"`
    Settings   MasterSettings `json:"settings"`
    CreatedAt  time.Time      `json:"created_at"`
    UpdatedAt  time.Time      `json:"updated_at"`
}

// Тестовая функция для проверки создания объекта
// func Get() Master {
// 	return Master{
// 		ID:         1,
// 		TelegramID: 12345678910,
// 		Name:       "Педро Барбер",
// 		Bio:        "Барбер · 6 лет опыта",
// 		Address:    "ул. Центральная, 1",
// 		InviteLink: "2hq0og3k",
// 		WorkHours: WorkHours{
// 			WorkingDays: []string{"Пн", "Вт", "Ср", "Чт", "Пт"},
// 			StartDay:    "09:00",
// 			EndDay:      "20:00",
// 			BreakStart:  "13:00",
// 			BreakEnd:    "14:00",
// 		},
// 		Settings: MasterSettings{
// 			AutoCancel:        true,
// 			CancelBeforeHours: 2,
// 			UseWaitlist:       true,
// 		},
// 		CreatedAt: time.Now(),
// 		UpdatedAt: time.Now(),
// 	}
// }
