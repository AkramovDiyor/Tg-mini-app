package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingService struct {
	pool        *pgxpool.Pool
	slotRepo    repositories.SlotRepository
	bookingRepo repositories.BookingRepository
	serviceRepo repositories.ServiceRepository
	clientRepo repositories.ClientRepository
}

func NewBookingService(pool *pgxpool.Pool, slotRepo repositories.SlotRepository, bookingRepo repositories.BookingRepository, serviceRepo repositories.ServiceRepository, clientRepo repositories.ClientRepository) *BookingService {
	return &BookingService{
		pool:        pool,
		slotRepo:    slotRepo,
		bookingRepo: bookingRepo,
		serviceRepo: serviceRepo,
		clientRepo: clientRepo,
	}
}

func (s *BookingService) BookSlot(ctx context.Context, clientTelegramID int64, serviceID int64, clientName string, priceLocked int, startTime time.Time) error {
    // 🔥 ШАГ 1: Создаем клиента ВНЕ транзакции бронирования
    // Это гарантирует, что клиент будет закоммичен ДО начала транзакции bookings
    _, err := s.clientRepo.GetOrCreateClient(ctx, clientTelegramID, "", clientName)
    if err != nil {
        log.Printf("❌ Ошибка создания клиента: %v", err)
        return fmt.Errorf("failed to create client: %w", err)
    }
    log.Printf("✅ Клиент %d создан/обновлен", clientTelegramID)

    // 🔥 ШАГ 2: Теперь начинаем транзакцию для бронирования
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

    // 3. Получаем услугу
    service, err := s.serviceRepo.GetServiceByID(ctx, serviceID)
    if err != nil {
        return err
    }

    // 4. Вычисляем конец слота
    endTime := startTime.Add(time.Duration(service.DurationMin) * time.Minute)

    // 5. Получаем или создаем слот
    slot, err := s.slotRepo.GetSlotByStartTimeAndMaster(ctx, tx, service.MasterID, startTime)
    if err != nil {
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
        if slot.Status != "free" && slot.Status != "booked" {
            return errors.New("слот уже занят")
        }
        if slot.Status == "free" {
            err = s.slotRepo.UpdateSlotStatus(ctx, tx, slot.ID, "booked")
            if err != nil {
                return err
            }
        }
    }

// 6. Создаем бронирование
newBooking := models.Booking{
    SlotID:           slot.ID,
    ServiceID:        serviceID,
    MasterID:         service.MasterID,  // 🔥 ДОБАВЛЕНО
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

func (s *BookingService) CancelBooking(ctx context.Context, bookingID int64, clientTelegramID int64) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			log.Println("Откат транзакции отмены:", err)
			tx.Rollback(ctx)
		}
	}()

	// 1. Получаем запись с блокировкой строки (FOR UPDATE)
	booking, err := s.bookingRepo.GetBookingByID(ctx, tx, bookingID)
	if err != nil {
		return err
	}

	// 2. ПРОВЕРКА БЕЗОПАСНОСТИ: клиент может отменить только СВОЮ запись
	if booking.ClientTelegramID != clientTelegramID {
		err = errors.New("вы не можете отменить чужую запись")
		return err
	}

	// 3. Проверяем статус записи
	if booking.Status == models.BookingStatusCancelledByClient {
		err = errors.New("запись уже отменена")
		return err
	}
	if booking.Status == models.BookingStatusCompleted {
		err = errors.New("нельзя отменить завершенную запись")
		return err
	}

	// 4. Проверяем, что запись еще не началась (получаем слот)
	slot, err := s.slotRepo.GetSlotByID(ctx, tx, booking.SlotID)
	if err != nil {
		return err
	}
	
	if slot.StartTime.Before(time.Now()) {
		err = errors.New("нельзя отменить запись, которая уже началась")
		return err
	}

	// 5. Меняем статус записи на "cancelled_by_client"
	err = s.bookingRepo.CancelBooking(ctx, tx, bookingID)
	if err != nil {
		return err
	}

	// 6. Освобождаем слот (теперь он снова доступен для записи)
	err = s.slotRepo.UpdateSlotStatus(ctx, tx, booking.SlotID, models.SlotStatusFree)
	if err != nil {
		return err
	}

	// 7. Коммитим транзакцию
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	log.Printf("✅ Запись #%d отменена клиентом %d", bookingID, clientTelegramID)
	return nil
}