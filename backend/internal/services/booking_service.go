package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingService struct {
	pool        *pgxpool.Pool
	slotRepo    repositories.SlotRepository
	bookingRepo repositories.BookingRepository
	serviceRepo repositories.ServiceRepository
}

func NewBookingService(pool *pgxpool.Pool, slotRepo repositories.SlotRepository, bookingRepo repositories.BookingRepository, serviceRepo repositories.ServiceRepository) *BookingService {
	return &BookingService{
		pool:        pool,
		slotRepo:    slotRepo,
		bookingRepo: bookingRepo,
		serviceRepo: serviceRepo,
	}
}

func (s *BookingService) BookSlot(ctx context.Context, clientTelegramID int64, serviceID int64, clientName string, priceLocked int, startTime time.Time) error {
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

	// 1. Получаем услугу по ID
	service, err := s.serviceRepo.GetServiceByID(ctx, serviceID)
	if err != nil {
		return err
	}

	// 2. Вычисляем конец слота
	endTime := startTime.Add(time.Duration(service.DurationMin) * time.Minute)

	// 3. Проверяем, существует ли уже слот на это время
	slot, err := s.slotRepo.GetSlotByStartTimeAndMaster(ctx, tx, service.MasterID, startTime)
	if err != nil {
		// Слота не существует — создаем новый
		newSlot := models.Slot{
			MasterID:  service.MasterID,
			StartTime: startTime,
			EndTime:   endTime,
			Status:    "booked",
		}
		slotID, err := s.slotRepo.CreateSlotWithID(ctx, tx, newSlot)
		if err != nil {
			return err
		}
		slot.ID = slotID
	} else {
		// Слот существует — проверяем статус
		if slot.Status != "free" && slot.Status != "booked" {
			return errors.New("слот уже занят")
		}
		// Если слот уже booked (например, из seed), используем его
		if slot.Status == "free" {
			err = s.slotRepo.UpdateSlotStatus(ctx, tx, slot.ID, "booked")
			if err != nil {
				return err
			}
		}
	}

	// 4. Создаем запись
	newBooking := models.Booking{
		SlotID:           slot.ID,
		ServiceID:        serviceID,
		ClientTelegramID: clientTelegramID,
		ClientName:       clientName,
		PriceLocked:      priceLocked,
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

	log.Println("✅ Клиент успешно записан!")
	return nil
}