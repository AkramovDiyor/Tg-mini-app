package repositories

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ClientBookingResponse описывает полный ответ для клиента
type ClientBookingResponse struct {
	BookingID        int64     `json:"booking_id"`
	ClientName       string    `json:"client_name"`
	ServiceName      string    `json:"service_name"`
	ServicePrice     int       `json:"service_price"`
	MasterName       string    `json:"master_name"`
	MasterAddress    string    `json:"master_address"`
	MasterInviteLink string    `json:"master_invite_link"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
}

type BookingRepository interface {
	CreateBooking(ctx context.Context, tx pgx.Tx, booking models.Booking) error
	GetBookingBySlotID(ctx context.Context, slotID int64) (models.Booking, error)
	GetBookingsByMasterAndDate(ctx context.Context, masterID int64, date time.Time) ([]models.Booking, error)
	GetBookingsByClientTgID(ctx context.Context, tgID int64) ([]ClientBookingResponse, error)
	
	// НОВЫЕ методы для отмены
	GetBookingByID(ctx context.Context, tx pgx.Tx, bookingID int64) (models.Booking, error)
	CancelBooking(ctx context.Context, tx pgx.Tx, bookingID int64) error
}

type BookingRepo struct {
	db *pgxpool.Pool
}

func NewBookingRepo(db *pgxpool.Pool) *BookingRepo {
	return &BookingRepo{db: db}
}

func (b *BookingRepo) CreateBooking(ctx context.Context, tx pgx.Tx, booking models.Booking) error {
    query := `INSERT INTO bookings 
        (slot_id, service_id, master_id, client_telegram_id, client_name, price_locked, status) 
        VALUES ($1, $2, $3, $4, $5, $6, $7)`
    
    _, err := tx.Exec(ctx, query, 
        booking.SlotID, 
        booking.ServiceID, 
        booking.MasterID,  // 🔥 ДОБАВЛЕНО
        booking.ClientTelegramID, 
        booking.ClientName, 
        booking.PriceLocked, 
        booking.Status,
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

    // 🔥 ИСПРАВЛЕНО: достаём ВСЕ поля из bookings
    query := `
        SELECT 
            b.id,
            b.slot_id,
            b.service_id,
            b.master_id,
            b.client_telegram_id,
            b.client_name,
            b.price_locked,
            b.status,
            s.start_time,
            s.end_time,
            b.created_at,
            b.updated_at
        FROM bookings b
        JOIN slots s ON b.slot_id = s.id
        WHERE s.master_id = $1 
          AND s.start_time >= $2 
          AND s.start_time < $3
          AND b.status NOT IN ('cancelled_by_client', 'cancelled_no_show')
        ORDER BY s.start_time ASC
    `
    
    rows, err := b.db.Query(ctx, query, masterID, startOfDay, endOfDay)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var bookings []models.Booking
    for rows.Next() {
        var bk models.Booking
        err := rows.Scan(
            &bk.ID,
            &bk.SlotID,
            &bk.ServiceID,
            &bk.MasterID,
            &bk.ClientTelegramID,
            &bk.ClientName,
            &bk.PriceLocked,
            &bk.Status,
            &bk.StartTime,
            &bk.EndTime,
            &bk.CreatedAt,
            &bk.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        bookings = append(bookings, bk)
    }
    return bookings, nil
}

func (b *BookingRepo) GetBookingsByClientTgID(ctx context.Context, tgID int64) ([]ClientBookingResponse, error) {
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
		  AND b.status NOT IN ('cancelled_by_client', 'cancelled_no_show')
		ORDER BY s.start_time ASC
		LIMIT 50
	`

	rows, err := b.db.Query(ctx, query, tgID)
	if err != nil {
		return nil, fmt.Errorf("database query error: %w", err)
	}
	defer rows.Close()

	bookings := []ClientBookingResponse{}

	for rows.Next() {
		var resp ClientBookingResponse
		err := rows.Scan(
			&resp.BookingID, &resp.ClientName, &resp.ServiceName, &resp.ServicePrice,
			&resp.MasterName, &resp.MasterAddress, &resp.MasterInviteLink,
			&resp.StartTime, &resp.EndTime, &resp.Status, &resp.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error on row: %w", err)
		}
		bookings = append(bookings, resp)
	}

	return bookings, nil
}

// 🔥 НОВЫЙ: Получаем запись по ID с блокировкой строки (FOR UPDATE)


func (b *BookingRepo) GetBookingByID(ctx context.Context, tx pgx.Tx, bookingID int64) (models.Booking, error) {
	var booking models.Booking
	var masterID sql.NullInt64  // 🔥 Nullable тип
	
	query := `
		SELECT id, slot_id, service_id, master_id, client_telegram_id, client_name, status, created_at 
		FROM bookings 
		WHERE id = $1 
		FOR UPDATE
	`
	
	err := tx.QueryRow(ctx, query, bookingID).Scan(
		&booking.ID,
		&booking.SlotID,
		&booking.ServiceID,
		&masterID,  // 🔥 Сканируем в nullable тип
		&booking.ClientTelegramID,
		&booking.ClientName,
		&booking.Status,
		&booking.CreatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.Booking{}, fmt.Errorf("запись не найдена")
		}
		return models.Booking{}, err
	}
	
	// Преобразуем в int64 (0 если NULL)
	if masterID.Valid {
		booking.MasterID = masterID.Int64
	}
	
	return booking, nil
}

// 🔥 НОВЫЙ: Мягкое удаление — меняем статус вместо DELETE
func (b *BookingRepo) CancelBooking(ctx context.Context, tx pgx.Tx, bookingID int64) error {
	query := `
		UPDATE bookings 
		SET status = $1, updated_at = NOW() 
		WHERE id = $2
	`
	result, err := tx.Exec(ctx, query, models.BookingStatusCancelledByClient, bookingID)
	if err != nil {
		return fmt.Errorf("failed to cancel booking: %w", err)
	}

	// Проверяем, что реально обновили строку
	if result.RowsAffected() == 0 {
		return fmt.Errorf("запись не найдена или уже отменена")
	}

	return nil
}