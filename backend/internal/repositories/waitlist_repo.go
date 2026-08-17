package repositories

import (
	"backend/internal/models"
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WaitlistRepository interface {
	AddToWaitlist(ctx context.Context, tx pgx.Tx, masterID int64, clientTgID int64, desiredDate time.Time) error
	GetFirstInWaitlist(ctx context.Context, masterID int64, desiredDate time.Time) (models.WaitlistRequest, error)
	UpdateWaitlistStatus(ctx context.Context, tx pgx.Tx, id int64, status string, slotID int64) error
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

func (w *WaitlistRepo) GetFirstInWaitlist(ctx context.Context, masterID int64, date time.Time) (models.WaitlistRequest, error) {

	var req models.WaitlistRequest

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
		&req.DesiredDate,
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
