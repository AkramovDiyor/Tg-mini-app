package main

import (
    "backend/internal/config"
    "backend/internal/delivery/httpapi"
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
    // 1. Загружаем конфигурацию
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("❌ Ошибка конфигурации: %v", err)
    }

    // 2. Подключаемся к базе данных
    pool, err := database.NewPostgres(cfg)
    if err != nil {
        log.Fatalf("❌ Ошибка подключения к БД: %v", err)
    }
    defer pool.Close()

    // 3. Инициализируем репозитории
    masterRepo := repositories.NewMasterRepo(pool)
    serviceRepo := repositories.NewServiceRepo(pool)
    slotRepo := repositories.NewSlotRepo(pool)
    bookingRepo := repositories.NewBookingRepo(pool)
	waitlistRepo := repositories.NewWaitlistRepo(pool)
    clientRepo := repositories.NewClientRepo(pool)

    // 4. Инициализируем сервисы
    bookingService := service.NewBookingService(pool, slotRepo, bookingRepo, serviceRepo, clientRepo)
    slotService := service.NewSlotService(masterRepo, serviceRepo, bookingRepo)

    // 5. Инициализируем HTTP-хендлеры
    bookingHandler := httpapi.NewBookingHandler(
        masterRepo, serviceRepo, slotRepo, bookingRepo, bookingService, slotService,
    )

	masterHandler := httpapi.NewMasterHandler(masterRepo, serviceRepo, bookingRepo, waitlistRepo)

    // 6. Создаем HTTP-роутер
    router := httpapi.NewRouter(bookingHandler, masterHandler, cfg.TgBotToken)

    // ============================================
    // 7. ЗАПУСК HTTP-СЕРВЕРА
    // ============================================
    httpServer := &http.Server{
        Addr:    ":8080",
        Handler: router,
    }

    // Запускаем HTTP-сервер в горутине
    go func() {
        log.Println("🚀 HTTP-сервер запущен на порту 8080")
        if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("❌ Ошибка HTTP-сервера: %v", err)
        }
	}()

    // ============================================
    // 8. GRACEFUL SHUTDOWN (Корректное завершение)
    // ============================================
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    
    // Блокируемся до получения сигнала завершения
    <-quit
    log.Println("🛑 Получен сигнал завершения. Останавливаем сервер...")

    // Даем серверу 10 секунд на завершение текущих запросов
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := httpServer.Shutdown(ctx); err != nil {
        log.Printf("❌ Ошибка graceful shutdown HTTP: %v", err)
    }

    log.Println("✅ Сервер остановлен. До встречи!")
}