package main

import (
    "backend/internal/config"
    "backend/internal/delivery/httpapi"
    "backend/internal/models"
    "backend/internal/repositories"
    service "backend/internal/services"
    "backend/pkg/database"
    "context"
    "log"
    "net/http"
    "time"

    "github.com/jackc/pgx/v5/pgxpool" // ДОБАВЛЕНО
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Ошибка конфигурации: %v", err)
    }

    pool, err := database.NewPostgres(cfg)
    if err != nil {
        log.Fatalf("Ошибка подключения к БД: %v", err)
    }
    defer pool.Close()

    // 3. Репозитории
    masterRepo := repositories.NewMasterRepo(pool)
    serviceRepo := repositories.NewServiceRepo(pool)
    slotRepo := repositories.NewSlotRepo(pool)
    bookingRepo := repositories.NewBookingRepo(pool)

    // Seed тестовых данных
    seedTestData(context.Background(), pool, masterRepo, serviceRepo, slotRepo, bookingRepo)

    // 4. Сервисы
    bookingService := service.NewBookingService(pool, slotRepo, bookingRepo, serviceRepo)
    slotService := service.NewSlotService(masterRepo, serviceRepo, bookingRepo)

    // 5. Хендлеры
    bookingHandler := httpapi.NewBookingHandler(masterRepo, serviceRepo, slotRepo, bookingRepo, bookingService, slotService)

    // 6. Роутер
    router := httpapi.NewRouter(bookingHandler, cfg.TgBotToken)

    // 7. Запуск
    log.Println("🚀 Сервер запущен на порту 8080. Готов к запросам от фронтенда!")
    log.Fatal(http.ListenAndServe(":8080", router))
}

// ИСПРАВЛЕНО: *pgxpool.Pool вместо *database.PgxPool
func seedTestData(ctx context.Context, pool *pgxpool.Pool, masterRepo *repositories.MasterRepo, serviceRepo *repositories.ServiceRepo, slotRepo *repositories.SlotRepo, bookingRepo *repositories.BookingRepo) {
    log.Println("🌱 Создаем тестовые данные...")

    // 1. Создаем мастера С ГРАФИКОМ РАБОТЫ
    master := models.Master{
        TelegramID: 111222333,
        Name:       "Педро Барбер",
        Bio:        "Барбер · 6 лет опыта",
        Address:    "ул. Центральная, 1",
        InviteLink: "LINK123243",
        WorkHours: models.WorkHours{
            WorkingDays: []string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"}, // Работает каждый день
            StartDay:    "10:00",
            EndDay:      "20:00",
            BreakStart:  "13:00",
            BreakEnd:    "14:00",
        },
        Settings: models.MasterSettings{
            AutoCancel:        true,
            CancelBeforeHours: 2,
            UseWaitlist:       true,
        },
    }

    err := masterRepo.CreateMaster(ctx, master)
    if err != nil {
        log.Printf("⚠️  Мастер уже существует, получаем существующего...")
    } else {
        log.Println("✅ Мастер создан с графиком работы")
    }

    // Получаем мастера
    master, err = masterRepo.GetMasterByInviteLink(ctx, "LINK123243")
    if err != nil {
        log.Fatalf("❌ Не удалось получить мастера: %v", err)
    }
    log.Printf("✅ Мастер ID: %d, Invite: %s", master.ID, master.InviteLink)
    log.Printf("📅 График: %s-%s, Обед: %s-%s", master.WorkHours.StartDay, master.WorkHours.EndDay, master.WorkHours.BreakStart, master.WorkHours.BreakEnd)

    // 2. Проверяем услуги
    services, err := serviceRepo.GetServicesByMasterID(ctx, master.ID)
    var serviceID int64
    
    if len(services) == 0 {
        service := models.Service{
            MasterID:    master.ID,
            Name:        "Мужская стрижка",
            DurationMin: 45,
            Price:       1500,
        }
        err = serviceRepo.CreateService(ctx, service)
        if err != nil {
            log.Printf("⚠️  Не удалось создать услугу: %v", err)
        } else {
            log.Println("✅ Услуга 'Мужская стрижка' создана")
        }
        
        services, _ = serviceRepo.GetServicesByMasterID(ctx, master.ID)
        if len(services) > 0 {
            serviceID = services[0].ID
        }
    } else {
        serviceID = services[0].ID
        log.Printf("✅ Услуга уже существует (ID: %d)", serviceID)
    }

    if serviceID == 0 {
        log.Fatal("❌ Не удалось получить ID услуги!")
    }

    // 3. Проверяем, какой день недели 2026-08-20
    testDate := time.Date(2026, 8, 20, 10, 0, 0, 0, time.UTC)
    weekday := testDate.Weekday()
    log.Printf("📅 Дата 2026-08-20 это %s", weekday)

    // 4. Создаем тестовый слот на 10:00
    existingSlots, _ := slotRepo.GetSlotsByMasterAndDate(ctx, master.ID, testDate)
    
    var testSlotID int64
    slotExists := false
    
    for _, s := range existingSlots {
        if s.StartTime.Hour() == 10 && s.StartTime.Minute() == 0 {
            testSlotID = s.ID
            slotExists = true
            log.Printf("✅ Тестовый слот уже существует (ID: %d)", testSlotID)
            break
        }
    }
    
    if !slotExists {
        testSlot := models.Slot{
            MasterID:  master.ID,
            StartTime: testDate,
            EndTime:   testDate.Add(45 * time.Minute),
            Status:    "booked",
        }
        err = slotRepo.CreateSlot(ctx, testSlot)
        if err != nil {
            log.Printf("⚠️  Не удалось создать слот: %v", err)
        } else {
            log.Println("✅ Тестовый слот на 10:00 создан")
            
            slots, _ := slotRepo.GetSlotsByMasterAndDate(ctx, master.ID, testDate)
            for _, s := range slots {
                if s.StartTime.Hour() == 10 && s.StartTime.Minute() == 0 {
                    testSlotID = s.ID
                    break
                }
            }
        }
    }

    // 5. Создаем бронирование
    if testSlotID > 0 {
        tx, err := pool.Begin(ctx)
        if err != nil {
            log.Printf("⚠️  Не удалось начать транзакцию: %v", err)
        } else {
            testBooking := models.Booking{
                SlotID:           testSlotID,
                ServiceID:        serviceID,
                ClientTelegramID: 777111222,
                ClientName:       "Вася №1",
                PriceLocked:      1500,
                Status:           "active",
            }
            
            err = bookingRepo.CreateBooking(ctx, tx, testBooking)
            if err != nil {
                log.Printf("⚠️  Бронирование уже существует: %v", err)
                tx.Rollback(ctx)
            } else {
                tx.Commit(ctx)
                log.Println("✅ Тестовое бронирование 'Вася №1' создано")
            }
        }
    }

    log.Println("🎯 Тестовые данные готовы!")
    log.Printf("📅 Тестируй с датой: 2026-08-20")
    log.Printf("🔗 Invite link: LINK123243")
    log.Printf("🆔 Service ID: %d", serviceID)
}