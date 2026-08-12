package repositories

import (
	"backend/internal/models"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)


type SlotRepository interface {
	CreateSlot(ctx context.Context, slot models.Slot) error
	GetSlotsByMasterAndDate(ctx context.Context, masterID int64, date time.Time) ([]models.Slot, error)
	UpdateSlotStatus(ctx context.Context, slotID int64, status string) error
}


type SlotRepo struct {
	db *pgxpool.Pool
}

func NewSlotRepo(db *pgxpool.Pool) *SlotRepo{
	return &SlotRepo{
		db: db,
	}
}


func (s *SlotRepo)CreateSlot(ctx context.Context, slot models.Slot) error  {
	_, err := s.db.Exec(ctx, "INSERT INTO slots (master_id, start_time, end_time, status) VALUES ($1, $2, $3, $4)", slot.MasterID, slot.StartTime, slot.EndTime, slot.Status)
	return err
}

func (s *SlotRepo) GetSlotsByMasterAndDate(ctx context.Context, masterID int64, date time.Time) ([]models.Slot, error) {
    // Нам нужен диапазон: от начала переданного дня (00:00:00) до конца (23:59:59)
    startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
    endOfDay := startOfDay.Add(24 * time.Hour)

    rows, err := s.db.Query(ctx, 
        "SELECT id, master_id, start_time, end_time, status FROM slots WHERE master_id = $1 AND start_time >= $2 AND start_time < $3", 
        masterID, startOfDay, endOfDay,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var slots []models.Slot
    for rows.Next() {
        var s models.Slot
        if err := rows.Scan(&s.ID, &s.MasterID, &s.StartTime, &s.EndTime, &s.Status); err != nil {
            return nil, err
        }
        slots = append(slots, s)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return slots, nil
}