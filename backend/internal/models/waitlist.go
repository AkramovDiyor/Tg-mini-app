// models/waitlist.go
package models

import "time"

// Константы для статусов
const (
    WaitlistStatusWaiting   = "waiting"
    WaitlistStatusOffered   = "offered"
    WaitlistStatusFulfilled = "fulfilled"
    WaitlistStatusExpired   = "expired"
)


type Waitlist struct {
    ID               int64      `json:"id"`
    MasterID         int64      `json:"master_id"`
    ClientTelegramID int64      `json:"client_telegram_id"`
    ClientName       string     `json:"client_name"`
    ClientPhone      string     `json:"client_phone"`
    ServiceName      string     `json:"service_name"`
    DesiredDate      time.Time  `json:"desired_date"`
    DesiredTime      *string    `json:"desired_time"`      // может быть null
    Status           string     `json:"status"`
    OfferedSlotID    *int64     `json:"offered_slot_id,omitempty"`
    CreatedAt        time.Time  `json:"created_at"`
    UpdatedAt        time.Time  `json:"updated_at"`
}