package main

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/repositories"
	"backend/internal/service"
	"backend/pkg/database"
	"context"
	"time"

	// "internal/poll"
	"log"
)

func main() {
	// 1. Загружаем переменные из .env
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("%v", err)
	}

	// 2. Создаем тот самый pool (электростанцию)
	pool, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer pool.Close() // Закроет соединения, когда сервер выключится

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	// 4. Передаем pool в репозиторий!
	masterRepo := repositories.NewMasterRepo(pool)

	testMaster := models.Get()
	testMaster.InviteLink = "LINK////1232431212"
	testMaster.TelegramID = time.Now().UnixNano()

	err = masterRepo.CreateMaster(ctx, testMaster)
	if err != nil {
		log.Fatalf("Ошибка создания: %v", err)
	}
	log.Println("Мастер успешно создан!")

	savedMaster, err := masterRepo.GetMasterByInviteLink(ctx, "LINK////123243")
	if err != nil {
		log.Fatalf("Ошибка чтения по ссылке: %v", err)
	}

	// Выводим имя прочитанного мастера, чтобы убедиться, что данные совпали
	log.Printf("Успешно прочитали из базы! Имя мастера: %s, ID в базе: %d", savedMaster.Name, savedMaster.ID)

	// serviceRepo := repositories.NewServiceRepo(pool)

	// 2. Создаем тестовую услугу для нашего сохраненного мастера
	// Берем savedMaster.ID, который нам только что вернул Postgres!
	// testService := models.Service{
	// 	MasterID:    savedMaster.ID,
	// 	Name:        "Мужская стрижка",
	// 	DurationMin: 45,
	// 	Price:       1500,
	// }

	// 3. Тестируем создание услуги
	// err = serviceRepo.CreateService(ctx, testService)
	// if err != nil {
	// 	log.Fatalf("Ошибка создания услуги: %v", err)
	// }
	// log.Println("✂️ Услуга успешно создана!")

	// 4. Тестируем получение списка услуг этого мастера
	// services, err := serviceRepo.GetServicesByMasterID(ctx, savedMaster.ID)
	// if err != nil {
	// 	log.Printf("Ошибка получения услуг: %v", err)
	// }

	// log.Printf("📚 Успешно! Найдено услуг у мастера: %d", len(services))
	// for _, s := range services {
	// 	log.Printf(" - %s (%d мин) — %d ₽", s.Name, s.DurationMin, s.Price)
	// }

	// ==========================================
	// ТЕСТ РЕПОЗИТОРИЯ СЛОТОВ (SlotRepo)
	// ==========================================

	// 1. Инициализируем репозиторий слотов
	// slotRepo := repositories.NewSlotRepo(pool)

	// 2. Создаем временные точки на СЕГОДНЯ
	// Пусть первый слот будет сегодня в 15:00
	// now := time.Now()
	// slotStart := time.Date(now.Year(), now.Month(), now.Day(), 15, 0, 0, 0, time.Local)
	// slotEnd := slotStart.Add(45 * time.Minute) // Длительность 45 минут

	// testSlot := models.Slot{
	// 	MasterID:  savedMaster.ID, // Привязываем к нашему мастеру из базы
	// 	StartTime: slotStart,
	// 	EndTime:   slotEnd,
	// 	Status:    "free", // Слот свободен для записи
	// }

	// // 3. Тестируем создание слота
	// // Оборачиваем в проверку, так как при повторном запуске база может выдать ошибку
	// // из-за нашего триггера EXCLUDE (защиты от пересечений слотов), и это нормально!
	// err = slotRepo.CreateSlot(ctx, testSlot)
	// if err != nil {
	// 	log.Printf("Предупреждение при создании слота (возможно, уже существует или пересекается): %v", err)
	// } else {
	// 	log.Println("📅 Слот расписания успешно создан!")
	// }

	// // 4. Тестируем получение слотов на сегодняшний день
	// // Передаем текущее время (метод сам сбросит часы до 00:00:00)
	// slots, err := slotRepo.GetSlotsByMasterAndDate(ctx, savedMaster.ID, now)
	// if err != nil {
	// 	log.Fatalf("Ошибка получения слотов: %v", err)
	// }

	// log.Printf("📊 Успешно! Найдено слотов на сегодня у мастера: %d", len(slots))
	// for _, slot := range slots {
	// 	log.Printf(" - Слот №%d: %s -> %s [Статус: %s]",
	// 		slot.ID,
	// 		slot.StartTime.Format("15:04"),
	// 		slot.EndTime.Format("15:04"),
	// 		slot.Status,
	// 	)
	// }

	// ==========================================
	// ТЕСТ БИЗНЕС-ЛОГИКИ (BookingService)
	// ==========================================

	// 1. Инициализируем репозитории
	slotRepo := repositories.NewSlotRepo(pool)
	serviceRepo := repositories.NewServiceRepo(pool)
	bookingRepo := repositories.NewBookingRepo(pool)

	// 2. Инициализируем СЕРВИС (Мозг)
	bookingService := service.NewBookingService(pool, slotRepo, bookingRepo)

	// 3. Подготавливаем данные (Услуга и Слот)
	// Создаем услугу
	testService := models.Service{
		MasterID:    savedMaster.ID,
		Name:        "Мужская стрижка",
		DurationMin: 45,
		Price:       1500,
	}
	err = serviceRepo.CreateService(ctx, testService)
	if err != nil {
		log.Fatalf("Ошибка создания услуги: %v", err)
	}
	services, _ := serviceRepo.GetServicesByMasterID(ctx, savedMaster.ID)
	serviceID := services[0].ID // Берем ID первой услуги

	// Создаем слот на сегодня 15:00
	now := time.Now()
	slotStart := time.Date(now.Year(), now.Month(), now.Day(), 15, 0, 0, 0, time.Local)
	testSlot := models.Slot{
		MasterID:  savedMaster.ID,
		StartTime: slotStart,
		EndTime:   slotStart.Add(45 * time.Minute),
		Status:    "free",
	}

	// Если слот уже есть (при повторном запуске), просто получим его
	err = slotRepo.CreateSlot(ctx, testSlot)
	if err != nil {
		log.Printf("Слот уже существует (это ок): %v", err)
	}
	slots, _ := slotRepo.GetSlotsByMasterAndDate(ctx, savedMaster.ID, now)
	if len(slots) == 0 {
		log.Fatalf("Слоты не найдены, нечего бронировать!")
	}
	slotID := slots[0].ID // Берем ID первого слота на сегодня

	// 4. САМОЕ ГЛАВНОЕ: ПЫТАЕМСЯ ЗАПИСАТЬ КЛИЕНТА (Вася №1)
	log.Println("=> Попытка записи Васи №1...")
	err = bookingService.BookSlot(ctx, slotID, serviceID, 777111222) // 777111222 - Telegram ID Васи
	if err != nil {
		log.Fatalf("❌ Ошибка записи Васи №1: %v", err)
	}
	log.Println("✅ Вася №1 успешно записан! Слот заблокирован.")

	// 5. ПЫТАЕМСЯ ЗАПИСАТЬ ВТОРОГО КЛИЕНТА НА ТОТ ЖЕ СЛОТ (Вася №2)
	log.Println("=> Попытка записи Васи №2 на тот же слот...")
	err = bookingService.BookSlot(ctx, slotID, serviceID, 888999000) // 888999000 - Telegram ID второго Васи
	if err != nil {
		// Мы ЖДЕМ эту ошибку! Система должна сказать "слот уже занят"
		log.Printf("✅ Система отработала верно! Отказ Вася №2: %v", err)
	} else {
		log.Fatalf("❌ КРИТИЧЕСКАЯ ОШИБКА! Вася №2 записался на занятый слот!")
	}

	defer cancel()

}
