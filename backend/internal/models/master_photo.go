package models

import "time"

// MasterPhoto представляет запись о фото работы мастера
type MasterPhoto struct {
    ID        int64     `json:"id" db:"id"`
    MasterID  int64     `json:"master_id" db:"master_id"`
    URL       string    `json:"url" db:"url"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}