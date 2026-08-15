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

// 1. ИСПРАВЛЕНО: Теперь конструктор принимает пул (как NewSlotRepo и NewServiceRepo)
func NewBookingRepo(db *pgxpool.Pool) *BookingRepo {
	return &BookingRepo{
		db: db,
	}
}

// Метод создания бронирования (остался без изменений)
func (b *BookingRepo) CreateBooking(ctx context.Context, tx pgx.Tx, booking models.Booking) error {
	query := `INSERT INTO bookings (slot_id, service_id, client_tg_id, status) VALUES ($1, $2, $3, $4)`
	
	_, err := tx.Exec(ctx, query, booking.SlotID, booking.ServiceID, booking.ClientTelegramID, booking.Status)
	return err
}

// 2. ИСПРАВЛЕНО: Добавлен недостающий метод, который требовал компилятор
func (b *BookingRepo) GetBookingBySlotID(ctx context.Context, slotID int64) (models.Booking, error) {
	var booking models.Booking

	query := `
		SELECT id, slot_id, service_id, client_tg_id, status, created_at 
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
