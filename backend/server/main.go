package main

import (
	"backend/internal/config"
	"backend/internal/delivery/httpapi"
	tgbot "backend/internal/delivery/telegram"
	"backend/internal/repositories"
	service "backend/internal/services"
	"backend/pkg/database"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	moscow, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Fatalf("❌ Не удалось загрузить часовой пояс Москвы: %v", err)
	}
	time.Local = moscow
	log.Println("🕐 Часовой пояс установлен: Europe/Moscow")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Ошибка конфигурации: %v", err)
	}

	pool, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к БД: %v", err)
	}
	defer pool.Close()

	// Репозитории
	masterRepo := repositories.NewMasterRepo(pool)
	serviceRepo := repositories.NewServiceRepo(pool)
	slotRepo := repositories.NewSlotRepo(pool)
	bookingRepo := repositories.NewBookingRepo(pool)
	waitlistRepo := repositories.NewWaitlistRepo(pool)
	clientRepo := repositories.NewClientRepo(pool)
	photoRepo := repositories.NewPhotoRepo(pool)

	// Сервисы
	masterService := service.NewMasterService(masterRepo)
	bookingService := service.NewBookingService(pool, slotRepo, bookingRepo, serviceRepo, clientRepo)
	slotService := service.NewSlotService(masterRepo, serviceRepo, bookingRepo)

	// Seed тестовых данных
	ctxSeed := context.Background()
	master, err := masterService.RegisterMaster(ctxSeed, 999999, "Педро Барбер")
	if err != nil {
		log.Printf("⚠️ Ошибка при посеве мастера: %v", err)
		master, _ = masterRepo.GetMasterByTelegramID(ctxSeed, 999999)
	}
	_, err = clientRepo.GetOrCreateClient(ctxSeed, 777111222, "vasya_test", "Вася Тестовый")
	if err != nil {
		log.Printf("⚠️ Ошибка при посеве клиента: %v", err)
	}

	log.Println("----------------------------------------")
	log.Println("🌱 ТЕСТОВЫЕ ДАННЫЕ ГОТОВЫ!")
	log.Printf("🛠 Панель мастера: ?mode=master")
	log.Printf("✂️ Клиентская часть: ?startapp=%s", master.InviteLink)
	log.Println("----------------------------------------")

	// Хендлеры
	bookingHandler := httpapi.NewBookingHandler(
		masterRepo, serviceRepo, slotRepo, bookingRepo, photoRepo, bookingService, slotService,
	)
	masterHandler := httpapi.NewMasterHandler(masterRepo, serviceRepo, bookingRepo, waitlistRepo)
	photoHandler := httpapi.NewPhotoHandler(photoRepo, masterRepo)
	meHandler := httpapi.NewMeHandler(masterRepo) // 🔥 ДОБАВЛЕНО

	// Роутер
	router := httpapi.NewRouter(bookingHandler, masterHandler, photoHandler, meHandler, cfg.TgBotToken)

	// Telegram-бот
	telegramHandler := tgbot.NewHandler(masterService, cfg.WebAppURL)
	go tgbot.StartBotWithRetry(cfg.TgBotToken, telegramHandler)

	// HTTP-сервер
	httpServer := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: router,
	}

	go func() {
		log.Println("🚀 HTTP-сервер запущен на порту 8080")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Ошибка HTTP-сервера: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Получен сигнал завершения. Останавливаем сервер...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("❌ Ошибка graceful shutdown HTTP: %v", err)
	}

	log.Println("✅ Сервер остановлен. До встречи!")
}