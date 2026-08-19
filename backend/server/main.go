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
    // 1. Устанавливаем часовой пояс Москвы глобально
    moscow, err := time.LoadLocation("Europe/Moscow")
    if err != nil {
        log.Fatalf("❌ Не удалось загрузить часовой пояс Москвы: %v", err)
    }
    time.Local = moscow
    log.Println("🕐 Часовой пояс установлен: Europe/Moscow")

    // 2. Загружаем конфигурацию
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("❌ Ошибка конфигурации: %v", err)
    }

    // 3. Подключаемся к базе данных
    pool, err := database.NewPostgres(cfg)
    if err != nil {
        log.Fatalf("❌ Ошибка подключения к БД: %v", err)
    }
    defer pool.Close()

    // 4. Инициализируем репозитории
    masterRepo := repositories.NewMasterRepo(pool)
    serviceRepo := repositories.NewServiceRepo(pool)
    slotRepo := repositories.NewSlotRepo(pool)
    bookingRepo := repositories.NewBookingRepo(pool)
    waitlistRepo := repositories.NewWaitlistRepo(pool)
    clientRepo := repositories.NewClientRepo(pool)
    photoRepo := repositories.NewPhotoRepo(pool)

    // 5. Инициализируем сервисы
    masterService := service.NewMasterService(masterRepo)
    bookingService := service.NewBookingService(pool, slotRepo, bookingRepo, serviceRepo, clientRepo)
    slotService := service.NewSlotService(masterRepo, serviceRepo, bookingRepo)

    // ============================================
    // 🌱 ПОСЕВ ТЕСТОВЫХ ДАННЫХ (SEED)
    // ============================================
        // ============================================
    // 🌱 ПОСЕВ ТЕСТОВЫХ ДАННЫХ (SEED)
    // ============================================
    ctxSeed := context.Background()

    // Создаем мастера (ID: 999999 для AuthMiddleware)
    master, err := masterService.RegisterMaster(ctxSeed, 999999, "Педро Барбер")
    if err != nil {
        log.Printf("⚠️ Ошибка при посеве мастера (возможно уже существует): %v", err)
        master, _ = masterRepo.GetMasterByTelegramID(ctxSeed, 999999) // Достаем, если уже есть
    }

    // Создаем клиента (ID: 777111222 для AuthMiddleware)
    _, err = clientRepo.GetOrCreateClient(ctxSeed, 777111222, "vasya_test", "Вася Тестовый")
    if err != nil {
        log.Printf("⚠️ Ошибка при посеве клиента: %v", err)
    }

    log.Println("----------------------------------------")
    log.Println("🌱 ТЕСТОВЫЕ ДАННЫЕ ГОТОВЫ!")
    log.Printf("🛠 Панель мастера: http://localhost:5173?startapp=master")
    log.Printf("✂️ Клиентская часть: http://localhost:5173?startapp=%s", master.InviteLink)
    log.Println("----------------------------------------")

    // 6. Инициализируем HTTP-хендлеры
    bookingHandler := httpapi.NewBookingHandler(
        masterRepo, serviceRepo, slotRepo, bookingRepo, photoRepo, bookingService, slotService,
    )

    masterHandler := httpapi.NewMasterHandler(masterRepo, serviceRepo, bookingRepo, waitlistRepo)
    photoHandler := httpapi.NewPhotoHandler(photoRepo, masterRepo)

    // 7. Создаем HTTP-роутер
    router := httpapi.NewRouter(bookingHandler, masterHandler, photoHandler, cfg.TgBotToken)

    // ============================================
    // 8. ЗАПУСК HTTP-СЕРВЕРА
    // ============================================
    httpServer := &http.Server{
        Addr:    ":8080",
        Handler: router,
    }

    go func() {
        log.Println("🚀 HTTP-сервер запущен на порту 8080")
        if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("❌ Ошибка HTTP-сервера: %v", err)
        }
    }()

    // ============================================
    // 9. GRACEFUL SHUTDOWN
    // ============================================
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