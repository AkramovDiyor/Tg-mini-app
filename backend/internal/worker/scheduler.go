package worker

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"context"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Scheduler — воркер, который раз в минуту проверяет записи
// и отправляет напоминания за 2 часа до начала
type Scheduler struct {
	bot         *tgbotapi.BotAPI
	bookingRepo repositories.BookingRepository
	serviceRepo repositories.ServiceRepository
	masterRepo  repositories.MasterRepository
}

// NewScheduler — конструктор воркера
func NewScheduler(
	bot *tgbotapi.BotAPI,
	bookingRepo repositories.BookingRepository,
	serviceRepo repositories.ServiceRepository,
	masterRepo repositories.MasterRepository,
) *Scheduler {
	return &Scheduler{
		bot:         bot,
		bookingRepo: bookingRepo,
		serviceRepo: serviceRepo,
		masterRepo:  masterRepo,
	}
}

// Start — запускает бесконечный цикл проверки
func (s *Scheduler) Start(ctx context.Context) {
	// Тикер срабатывает каждую 1 минуту
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	log.Println("⏰ [SCHEDULER] Воркер напоминаний запущен (проверка каждую минуту)")

	// Сразу делаем первую проверку при старте
	s.processReminders(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Println("⏰ [SCHEDULER] Остановка воркера...")
			return
		case <-ticker.C:
			s.processReminders(ctx)
		}
	}
}

// processReminders — основная логика: ищем записи и шлем пуши
func (s *Scheduler) processReminders(ctx context.Context) {
	// 1. Получаем записи, которым пора отправлять напоминание
	bookings, err := s.bookingRepo.GetBookingsForReminder(ctx)
	if err != nil {
		log.Printf("❌ [SCHEDULER] Ошибка получения записей: %v", err)
		return
	}

	if len(bookings) == 0 {
		// log.Printf("💤 [SCHEDULER] Нет записей для напоминания") // Закомментируй, если не хочешь спам в логах
		return
	}

	log.Printf("📬 [SCHEDULER] Найдено записей для напоминания: %d", len(bookings))

	// 2. Обрабатываем каждую запись
	for _, booking := range bookings {
		s.sendReminder(ctx, booking)
	}
}

// sendReminder — отправляет одно конкретное напоминание
func (s *Scheduler) sendReminder(ctx context.Context, booking models.Booking) {
	// Получаем услугу
	service, err := s.serviceRepo.GetServiceByID(ctx, booking.ServiceID)
	if err != nil {
		log.Printf("❌ [SCHEDULER] Не удалось получить услугу %d: %v", booking.ServiceID, err)
		return
	}

	// Получаем мастера (для имени)
	master, err := s.masterRepo.GetMasterByID(ctx, booking.MasterID)
	if err != nil {
		log.Printf("❌ [SCHEDULER] Не удалось получить мастера %d: %v", booking.MasterID, err)
		return
	}

	// Форматируем время в часовом поясе Москвы
	moscow, _ := time.LoadLocation("Europe/Moscow")
	localTime := booking.StartTime.In(moscow).Format("15:04")
	localDate := booking.StartTime.In(moscow).Format("2 января (пн)")

	// Определяем сколько осталось времени
	timeLeft := time.Until(booking.StartTime)
	hoursLeft := int(timeLeft.Hours())
	minutesLeft := int(timeLeft.Minutes()) % 60

	var timeLeftStr string
	if hoursLeft > 0 {
		timeLeftStr = fmt.Sprintf("через %d ч. %d мин.", hoursLeft, minutesLeft)
	} else {
		timeLeftStr = fmt.Sprintf("через %d мин.", minutesLeft)
	}

	// Формируем красивое сообщение
	message := fmt.Sprintf(
		"⏰ *Напоминание о записи!*\n\n"+
			"👤 *Мастер:* %s\n"+
			"✂️ *Услуга:* %s\n"+
			"📅 *Дата:* %s\n"+
			"🕒 *Время:* %s\n"+
			"⏳ *Осталось:* %s\n\n"+
			"💰 *Сумма:* %d ₽\n\n"+
			"Ждём вас! 🤝",
		master.Name,
		service.Name,
		localDate,
		localTime,
		timeLeftStr,
		booking.PriceLocked,
	)

	// Отправляем клиенту
	msg := tgbotapi.NewMessage(booking.ClientTelegramID, message)
	msg.ParseMode = tgbotapi.ModeMarkdown

	_, err = s.bot.Send(msg)
	if err != nil {
		log.Printf("❌ [SCHEDULER] Ошибка отправки клиенту %d: %v", booking.ClientTelegramID, err)
		// Не помечаем как отправленное, попробуем в следующий раз
		return
	}

	// Помечаем как отправленное в БД
	err = s.bookingRepo.MarkReminderSent(ctx, booking.ID)
	if err != nil {
		log.Printf("❌ [SCHEDULER] Не удалось пометить как отправленное (ID: %d): %v", booking.ID, err)
		return
	}

	log.Printf("✅ [SCHEDULER] Напоминание отправлено: booking=%d, client=%d, время=%s", 
		booking.ID, booking.ClientTelegramID, localTime)
}