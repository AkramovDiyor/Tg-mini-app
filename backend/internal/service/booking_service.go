package service

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingService struct {
	pool        *pgxpool.Pool
	slotRepo    repositories.SlotRepository
	bookingRepo repositories.BookingRepository
}

func NewBookingService(pool *pgxpool.Pool, slotRepo repositories.SlotRepository, bookingRepo repositories.BookingRepository) *BookingService {
	return &BookingService{
		pool:        pool,
		slotRepo:    slotRepo,
		bookingRepo: bookingRepo,
	}
}

func (s *BookingService) BookSlot(ctx context.Context, slotID int64, serviceID int64, clientTelegramID int64) error {
	tx, err := s.pool.Begin(ctx)

	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			log.Println("Откат транзакции из-за ошибки:", err)
			tx.Rollback(ctx)
		}
	}()

	slot, err := s.slotRepo.GetSlotByID(ctx, tx, slotID)
	if err != nil {
		return err
	}

	if slot.Status != "free" {
		err = errors.New("слот уже занят")
		return  err
	}

	err = s.slotRepo.UpdateSlotStatus(ctx, tx, slotID, "booked")
	if err != nil {
		return err
	}

	newBooking := models.Booking{
		SlotID:           slotID,
		ServiceID:        serviceID,
		ClientTelegramID: clientTelegramID,
		Status:           "active",
	}
	err = s.bookingRepo.CreateBooking(ctx, tx, newBooking)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	log.Println("Клиент успешно записан!")
    return nil
}
