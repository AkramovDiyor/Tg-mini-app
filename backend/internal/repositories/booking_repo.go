package repositories

import (
	"backend/internal/models"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ClientBookingResponse описывает полный ответ для клиента
type ClientBookingResponse struct {
	BookingID      int64     `json:"booking_id"`
	ClientName     string    `json:"client_name"`
	ServiceName    string    `json:"service_name"`
	ServicePrice   int       `json:"service_price"`
	MasterName     string    `json:"master_name"`
	MasterAddress  string    `json:"master_address"`
	MasterInviteLink string  `json:"master_invite_link"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

type BookingRepository interface {
	CreateBooking(ctx context.Context, tx pgx.Tx, booking models.Booking) error
	GetBookingBySlotID(ctx context.Context, slotID int64) (models.Booking, error)
	GetBookingsByMasterAndDate(ctx context.Context, masterID int64, date time.Time) ([]models.Booking, error)
	GetBookingsByClientTgID(ctx context.Context, tgID int64) ([]ClientBookingResponse, error) // НОВОЕ
}

type BookingRepo struct {
	db *pgxpool.Pool
}

func NewBookingRepo(db *pgxpool.Pool) *BookingRepo {
	return &BookingRepo{db: db}
}

func (b *BookingRepo) CreateBooking(ctx context.Context, tx pgx.Tx, booking models.Booking) error {
	query := `INSERT INTO bookings 
        (slot_id, service_id, client_telegram_id, client_name, price_locked, status) 
        VALUES ($1, $2, $3, $4, $5, $6)`
	
	_, err := tx.Exec(ctx, query, 
		booking.SlotID, booking.ServiceID, booking.ClientTelegramID, 
		booking.ClientName, booking.PriceLocked, booking.Status,
	)
	return err
}

func (b *BookingRepo) GetBookingBySlotID(ctx context.Context, slotID int64) (models.Booking, error) {
	var booking models.Booking
	query := `SELECT id, slot_id, service_id, client_telegram_id, status, created_at FROM bookings WHERE slot_id = $1`
	err := b.db.QueryRow(ctx, query, slotID).Scan(
		&booking.ID, &booking.SlotID, &booking.ServiceID, &booking.ClientTelegramID, &booking.Status, &booking.CreatedAt,
	)
	return booking, err
}

func (b *BookingRepo) GetBookingsByMasterAndDate(ctx context.Context, masterID int64, date time.Time) ([]models.Booking, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := `
		SELECT s.start_time, s.end_time, b.status
		FROM bookings b
		JOIN slots s ON b.slot_id = s.id
		WHERE s.master_id = $1 AND s.start_time >= $2 AND s.start_time < $3
		AND b.status NOT IN ('cancelled_by_client', 'cancelled_no_show')
	`
	rows, err := b.db.Query(ctx, query, masterID, startOfDay, endOfDay)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var bk models.Booking
		err := rows.Scan(&bk.StartTime, &bk.EndTime, &bk.Status)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, bk)
	}
	return bookings, nil
}

// НОВЫЙ МЕТОД: Достаем все активные записи клиента с полной информацией
func (b *BookingRepo) GetBookingsByClientTgID(ctx context.Context, tgID int64) ([]ClientBookingResponse, error) {
	log.Printf("🔍 GetBookingsByClientTgID: ищем записи для tgID=%d", tgID)

	query := `
		SELECT 
			b.id,
			b.client_name,
			COALESCE(sv.name, 'Без названия') AS service_name,
			COALESCE(sv.price, 0) AS service_price,
			COALESCE(m.name, 'Без имени') AS master_name,
			COALESCE(m.address, 'Без адреса') AS master_address,
			COALESCE(m.invite_link, '') AS master_invite_link,
			s.start_time,
			s.end_time,
			b.status,
			b.created_at
		FROM bookings b
		LEFT JOIN slots s ON b.slot_id = s.id
		LEFT JOIN services sv ON b.service_id = sv.id
		LEFT JOIN masters m ON s.master_id = m.id
		WHERE b.client_telegram_id = $1
		ORDER BY s.start_time DESC
		LIMIT 50
	`

	rows, err := b.db.Query(ctx, query, tgID)
	if err != nil {
		log.Printf("❌ Database query error: %v", err)
		return nil, fmt.Errorf("database query error: %w", err)
	}
	defer rows.Close()

	bookings := []ClientBookingResponse{}

	for rows.Next() {
		var resp ClientBookingResponse
		err := rows.Scan(
			&resp.BookingID,
			&resp.ClientName,
			&resp.ServiceName,
			&resp.ServicePrice,
			&resp.MasterName,
			&resp.MasterAddress,
			&resp.MasterInviteLink,
			&resp.StartTime,
			&resp.EndTime,
			&resp.Status,
			&resp.CreatedAt,
		)
		if err != nil {
			log.Printf("❌ Scan error: %v", err)
			return nil, fmt.Errorf("scan error on row: %w", err)
		}
		bookings = append(bookings, resp)
	}

	log.Printf("✅ Найдено записей: %d", len(bookings))

	if err = rows.Err(); err != nil {
		log.Printf("❌ Rows iteration error: %v", err)
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return bookings, nil
}