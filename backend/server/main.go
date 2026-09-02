package main

import (
	"backend/internal/config"
	"backend/internal/delivery/httpapi"
	tgbot "backend/internal/delivery/telegram"
	"backend/internal/repositories"
	service "backend/internal/services"
	"backend/internal/worker"
	"backend/pkg/database"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// Устанавливаем часовой пояс Москвы
	moscow, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Fatalf("Failed to load Moscow timezone: %v", err)
	}
	time.Local = moscow

	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Подключаемся к базе данных
	pool, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()
	log.Println("Database connected")

	// Инициализируем Telegram бота
	bot, err := tgbotapi.NewBotAPI(cfg.TgBotToken)
	if err != nil {
		log.Fatalf("Failed to initialize Telegram bot: %v", err)
	}
	log.Printf("Bot authorized as @%s", bot.Self.UserName)

	// Инициализируем репозитории
	masterRepo := repositories.NewMasterRepo(pool)
	serviceRepo := repositories.NewServiceRepo(pool)
	slotRepo := repositories.NewSlotRepo(pool)
	bookingRepo := repositories.NewBookingRepo(pool)
	waitlistRepo := repositories.NewWaitlistRepo(pool)
	clientRepo := repositories.NewClientRepo(pool)
	photoRepo := repositories.NewPhotoRepo(pool)

	// Инициализируем сервисы
	masterService := service.NewMasterService(masterRepo)
	bookingService := service.NewBookingService(pool, slotRepo, bookingRepo, serviceRepo, clientRepo, masterRepo, bot)
	slotService := service.NewSlotService(masterRepo, serviceRepo, bookingRepo)

	// Инициализируем HTTP хендлеры
	bookingHandler := httpapi.NewBookingHandler(
		masterRepo, serviceRepo, slotRepo, bookingRepo, photoRepo, bookingService, slotService,
	)
	masterHandler := httpapi.NewMasterHandler(masterRepo, serviceRepo, bookingRepo, waitlistRepo)
	photoHandler := httpapi.NewPhotoHandler(photoRepo, masterRepo)
	meHandler := httpapi.NewMeHandler(masterRepo)

	// Создаем роутер
	router := httpapi.NewRouter(bookingHandler, masterHandler, photoHandler, meHandler, cfg.TgBotToken)

	// Запускаем Telegram бота
	telegramHandler := tgbot.NewHandler(masterService, cfg.WebAppURL, bot)

	go func() {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		updates := bot.GetUpdatesChan(u)

		log.Println("Telegram bot started")

		for update := range updates {
			go func(upd tgbotapi.Update) {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Panic in bot handler: %v", r)
					}
				}()
				if upd.Message != nil {
					telegramHandler.HandleMessage(&upd)
				}
			}(update)
		}
	}()

	// Запускаем воркер напоминаний
	workerCtx, workerCancel := context.WithCancel(context.Background())
	scheduler := worker.NewScheduler(bot, bookingRepo, serviceRepo, masterRepo)
	go scheduler.Start(workerCtx)
	log.Println("Reminder worker started")

	// Запускаем HTTP сервер
	httpServer := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: router,
	}

	go func() {
		log.Println("HTTP server started on :8080")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Ждем сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Останавливаем воркер
	workerCancel()

	// Graceful shutdown HTTP сервера
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}