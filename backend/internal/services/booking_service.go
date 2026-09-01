package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingService struct {
	pool        *pgxpool.Pool
	slotRepo    repositories.SlotRepository
	bookingRepo repositories.BookingRepository
	serviceRepo repositories.ServiceRepository
	clientRepo  repositories.ClientRepository
	masterRepo  repositories.MasterRepository // 🔥 ДОБАВЛЕНО: для получения Telegram ID мастера
	bot         *tgbotapi.BotAPI
}

func NewBookingService(
	pool *pgxpool.Pool,
	slotRepo repositories.SlotRepository,
	bookingRepo repositories.BookingRepository,
	serviceRepo repositories.ServiceRepository,
	clientRepo repositories.ClientRepository,
	masterRepo repositories.MasterRepository, // 🔥 ДОБАВЛЕНО
	bot *tgbotapi.BotAPI,
) *BookingService {
	return &BookingService{
		pool:        pool,
		slotRepo:    slotRepo,
		bookingRepo: bookingRepo,
		serviceRepo: serviceRepo,
		clientRepo:  clientRepo,
		masterRepo:  masterRepo, // 🔥 ДОБАВЛЕНО
		bot:         bot,
	}
}

func (s *BookingService) BookSlot(ctx context.Context, clientTelegramID int64, serviceID int64, clientName string, priceLocked int, startTime time.Time) error {
	// ШАГ 1: Создаем или обновляем клиента
	// GetOrCreateClient должен вернуть нам актуальные данные клиента
	client, err := s.clientRepo.GetOrCreateClient(ctx, clientTelegramID, "", clientName)
	if err != nil {
		log.Printf("❌ Ошибка создания клиента: %v", err)
		return fmt.Errorf("failed to create client: %w", err)
	}

	// 🔥 ШАГ 1.5: Определяем реальное имя для уведомления
	// Если в БД есть имя, используем его. Если нет (или оно "Клиент"), используем то, что пришло.
	realClientName := clientName
	if client.FirstName != "" {
		realClientName = client.FirstName
		if client.LastName != "" {
			realClientName += " " + client.LastName
		}
	}
	
	// Если имя всё еще "Клиент", попробуем взять из Telegram ID (если у вас есть такой метод)
	// Но обычно GetOrCreateClient уже сохраняет правильное имя, если фронтенд его передал.

	// ШАГ 2: Начинаем транзакцию для бронирования
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

	// ... (твой существующий код получения услуги, слота и создания записи без изменений) ...
	service, err := s.serviceRepo.GetServiceByID(ctx, serviceID)
	if err != nil {
		return err
	}

	endTime := startTime.Add(time.Duration(service.DurationMin) * time.Minute)

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

	newBooking := models.Booking{
		SlotID:           slot.ID,
		ServiceID:        serviceID,
		MasterID:         service.MasterID,
		ClientTelegramID: clientTelegramID,
		ClientName:       realClientName, // 🔥 Используем реальное имя
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

	// 🔥 ШАГ 7: Отправляем уведомление с ПРАВИЛЬНЫМ именем
	go s.notifyMaster(service.MasterID, realClientName, service.Name, startTime, priceLocked, "Новая запись")

	log.Printf("✅ Клиент %s (ID: %d) успешно записан!", realClientName, clientTelegramID)
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

	// 1. Получаем запись
	booking, err := s.bookingRepo.GetBookingByID(ctx, tx, bookingID)
	if err != nil {
		return err
	}

	// 2. Проверка безопасности
	if booking.ClientTelegramID != clientTelegramID {
		return errors.New("вы не можете отменить чужую запись")
	}

	// 3. Проверяем статус
	if booking.Status == models.BookingStatusCancelledByClient {
		return errors.New("запись уже отменена")
	}
	if booking.Status == models.BookingStatusCompleted {
		return errors.New("нельзя отменить завершенную запись")
	}

	// 4. Проверяем, что запись еще не началась
	slot, err := s.slotRepo.GetSlotByID(ctx, tx, booking.SlotID)
	if err != nil {
		return err
	}
	if slot.StartTime.Before(time.Now()) {
		return errors.New("нельзя отменить запись, которая уже началась")
	}

	// 5. Меняем статус записи
	err = s.bookingRepo.CancelBooking(ctx, tx, bookingID)
	if err != nil {
		return err
	}

	// 6. Освобождаем слот
	err = s.slotRepo.UpdateSlotStatus(ctx, tx, booking.SlotID, models.SlotStatusFree)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	// 🔥 ШАГ 7: Отправляем уведомление об отмене в горутине
	// Нам нужно название услуги, поэтому делаем легкий запрос (или можно добавить его в GetBookingByID через JOIN)
	service, _ := s.serviceRepo.GetServiceByID(context.Background(), booking.ServiceID)
	serviceName := "Неизвестная услуга"
	if service.Name != "" {
		serviceName = service.Name
	}

	go s.notifyMaster(booking.MasterID, booking.ClientName, serviceName, slot.StartTime, booking.PriceLocked, "Запись отменена")

	log.Printf("✅ Запись #%d отменена клиентом %d", bookingID, clientTelegramID)
	return nil
}

// 🔥 УНИВЕРСАЛЬНЫЙ МЕТОД УВЕДОМЛЕНИЯ
func (s *BookingService) notifyMaster(masterID int64, clientName, serviceName string, startTime time.Time, price int, action string) {
	// 1. Получаем telegram_id мастера
	master, err := s.masterRepo.GetMasterByID(context.Background(), masterID)
	if err != nil {
		log.Printf("❌ Не удалось получить мастера для уведомления (ID: %d): %v", masterID, err)
		return
	}

	// 2. Форматируем время по Москве
	moscow, _ := time.LoadLocation("Europe/Moscow")
	localTime := startTime.In(moscow).Format("15:04, 2 января")

	// 3. Выбираем эмодзи в зависимости от действия
	emoji := "🔔"
	if action == "Запись отменена" {
		emoji = "❌"
	}

	// 4. Формируем сообщение
	message := fmt.Sprintf(
		"%s *%s!*\n\n"+
			"👤 *Клиент:* %s\n"+
			"✂️ *Услуга:* %s\n"+
			"🕒 *Время:* %s\n"+
			"💰 *Сумма:* %d ₽",
		emoji, action, clientName, serviceName, localTime, price,
	)

	// 5. Отправляем сообщение
	msg := tgbotapi.NewMessage(master.TelegramID, message)
	msg.ParseMode = tgbotapi.ModeMarkdown

	if _, err := s.bot.Send(msg); err != nil {
		log.Printf("❌ Ошибка отправки уведомления мастеру %d: %v", master.TelegramID, err)
	} else {
		log.Printf("✅ Уведомление '%s' отправлено мастеру %d", action, master.TelegramID)
	}
}