package repositories

import (
	"backend/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WaitlistRepository interface {
	AddToWaitlist(ctx context.Context, tx pgx.Tx, masterID int64, clientTgID int64, desiredDate time.Time) error
	GetFirstInWaitlist(ctx context.Context, masterID int64, desiredDate time.Time) (models.Waitlist, error)
	UpdateWaitlistStatus(ctx context.Context, tx pgx.Tx, id int64, status string, slotID int64) error
	GetWaitlistByMaster(ctx context.Context, masterID int64) ([]models.Waitlist, error)
}

type WaitlistRepo struct {
	db *pgxpool.Pool
}

func NewWaitlistRepo(db *pgxpool.Pool) *WaitlistRepo {
	return &WaitlistRepo{db: db}
}

func (w *WaitlistRepo) AddToWaitlist(ctx context.Context, tx pgx.Tx, masterID int64, clientTgID int64, desiredDate time.Time) error {
	query := `INSERT INTO waitlist (master_id, client_telegram_id, desired_date, status) 
	VALUES ($1, $2, $3, 'waiting')`

	_, err := tx.Exec(ctx, query, masterID, clientTgID, desiredDate.Format("2006-01-02"))
	return err
}

func (w *WaitlistRepo) GetFirstInWaitlist(ctx context.Context, masterID int64, date time.Time) (models.Waitlist, error) {

	var req models.Waitlist

	query := `SELECT id, master_id, client_telegram_id, desired_date, status, offered_slot_id, created_at, updated_at 
FROM waitlist
WHERE master_id = $1 AND desired_date = $2 AND status = 'waiting'
ORDER BY created_at ASC 
LIMIT 1;
`

	err := w.db.QueryRow(ctx, query, masterID, date.Format("2006-01-02")).Scan(
		&req.ID,
		&req.MasterID,
		&req.ClientTelegramID,
		// &req.DesiredDate,
		&req.Status,
		&req.OfferedSlotID,
		&req.CreatedAt,
		&req.UpdatedAt,
	)

	return req, err

}

func (w *WaitlistRepo) UpdateWaitlistStatus(ctx context.Context, tx pgx.Tx, id int64, status string, slotID int64) error {

	query := `UPDATE waitlist 
	SET status = $1, offered_slot_id = $2, updated_at = NOW() 
	WHERE id = $3`
	

	_, err := tx.Exec(ctx, query, status, slotID, id)
	return err
}

// repositories/waitlist_repo.go
func (r *WaitlistRepo) GetWaitlistByMaster(ctx context.Context, masterID int64) ([]models.Waitlist, error) {
    query := `
        SELECT id, master_id, client_telegram_id, client_name, client_phone,
               service_name, desired_date, desired_time, status, 
               offered_slot_id, created_at, updated_at
        FROM waitlist
        WHERE master_id = $1
        ORDER BY created_at ASC
    `
    
    rows, err := r.db.Query(ctx, query, masterID)
    if err != nil {
        return nil, fmt.Errorf("failed to get waitlist: %w", err)
    }
    defer rows.Close()
    
    var waitlist []models.Waitlist
    for rows.Next() {
        var w models.Waitlist
        err := rows.Scan(
            &w.ID,
            &w.MasterID,
            &w.ClientTelegramID,
            &w.ClientName,
            &w.ClientPhone,
            &w.ServiceName,
            &w.DesiredDate,
            &w.DesiredTime,
            &w.Status,
            &w.OfferedSlotID,
            &w.CreatedAt,
            &w.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan waitlist: %w", err)
        }
        waitlist = append(waitlist, w)
    }
    
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("rows iteration error: %w", err)
    }
    
    return waitlist, nil
}