package repositories

import (
	"backend/internal/models"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingRepository interface {
	CreateBooking(ctx context.Context, tx pgx.Tx, booking models.Booking) error
	GetBookingBySlotID(ctx context.Context, slotID int64) (models.Booking, error)
}

type BookingRepo struct {
	db *pgxpool.Pool
}

func NewBookingRepo(db *pgxpool.Pool) *BookingRepo {
	return &BookingRepo{
		db: db,
	}
}

// Метод создания бронирования (остался без изменений)
func (b *BookingRepo) CreateBooking(ctx context.Context, tx pgx.Tx, booking models.Booking) error {
    query := `INSERT INTO bookings 
        (slot_id, service_id, client_telegram_id, client_name, price_locked, status) 
        VALUES ($1, $2, $3, $4, $5, $6)`
    
    _, err := tx.Exec(ctx, query, 
        booking.SlotID, 
        booking.ServiceID, 
        booking.ClientTelegramID, 
        booking.ClientName,    
        booking.PriceLocked,  
        booking.Status,
    )
    return err
}

// 2. ИСПРАВЛЕНО: Добавлен недостающий метод, который требовал компилятор
func (b *BookingRepo) GetBookingBySlotID(ctx context.Context, slotID int64) (models.Booking, error) {
	var booking models.Booking

	query := `
		SELECT id, slot_id, service_id, client_telegram_id, status, created_at 
		FROM bookings 
		WHERE slot_id = $1
	`

	// Используем b.db (наш пул), так как этот запрос идет вне транзакции создания
	err := b.db.QueryRow(ctx, query, slotID).Scan(
		&booking.ID,
		&booking.SlotID,
		&booking.ServiceID,
		&booking.ClientTelegramID, // Проверьте, что в models.Booking поле называется именно так
		&booking.Status,
		&booking.CreatedAt,
	)

	return booking, err
}
