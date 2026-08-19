package services

import (
    "backend/internal/models"
    "backend/internal/repositories"
    "context"
    "fmt"
    "log"
    "sort"
    "time"
)

type Slot struct {
    StartTime time.Time `json:"start_time"`
    Status    string    `json:"status"`
}

type SlotService struct {
    MasterRepo  repositories.MasterRepository
    ServiceRepo repositories.ServiceRepository
    BookingRepo repositories.BookingRepository
}

func NewSlotService(mRepo repositories.MasterRepository, sRepo repositories.ServiceRepository, bRepo repositories.BookingRepository) *SlotService {
    return &SlotService{MasterRepo: mRepo, ServiceRepo: sRepo, BookingRepo: bRepo}
}

func (s *SlotService) GetAvailableSlots(ctx context.Context, inviteLink string, date time.Time, serviceID int64) ([]Slot, error) {
    moscow, _ := time.LoadLocation("Europe/Moscow")

    // 1. Получаем мастера
    master, err := s.MasterRepo.GetMasterByInviteLink(ctx, inviteLink)
    if err != nil {
        return nil, err
    }

    log.Printf("📊 Мастер %s: WorkDays=%v, StartTime=%s, EndTime=%s",
        master.Name, master.WorkHours.WorkDays, master.WorkHours.StartTime, master.WorkHours.EndTime)

    // 2. Получаем услугу
    services, err := s.ServiceRepo.GetServicesByMasterID(ctx, master.ID)
    if err != nil {
        return nil, err
    }

    var service models.Service
    found := false
    for _, srv := range services {
        if srv.ID == serviceID {
            service = srv
            found = true
            break
        }
    }
    if !found {
        return nil, fmt.Errorf("service not found for this master")
    }

    // 3. 🔥 ПРОВЕРКА: работает ли мастер в этот день
    // date.Weekday() возвращает: 0=Sunday, 1=Monday, ..., 6=Saturday
    // Но массив work_days: [Пн(0), Вт(1), Ср(2), Чт(3), Пт(4), Сб(5), Вс(6)]
    weekday := date.Weekday()
    dayIndex := int(weekday) - 1 // Пн=0, Вт=1, ..., Сб=5
    if weekday == time.Sunday {
        dayIndex = 6 // Воскресенье = индекс 6
    }

    // Проверяем, что массив существует и индекс валидный
    isWorkingDay := false
    if len(master.WorkHours.WorkDays) > dayIndex {
        isWorkingDay = master.WorkHours.WorkDays[dayIndex]
    }

    if !isWorkingDay {
        log.Printf("❌ День %s (индекс %d) - выходной для мастера %s", 
            weekday, dayIndex, master.Name)
        return []Slot{}, nil
    }

    log.Printf("✅ День %s (индекс %d) - рабочий для мастера %s", 
        weekday, dayIndex, master.Name)

    // 4. Нормализуем дату к началу дня в Москве
    dateInMoscow := time.Date(
        date.Year(), date.Month(), date.Day(),
        0, 0, 0, 0,
        moscow,
    )

    // 5. Получаем активные записи мастера на эту дату
    bookings, err := s.BookingRepo.GetBookingsByMasterAndDate(ctx, master.ID, dateInMoscow)
    if err != nil {
        return nil, err
    }

    // 6. Генерируем виртуальные слоты
    return generateVirtualSlots(master, service, bookings, dateInMoscow), nil
}

// generateVirtualSlots - ядро алгоритма генерации слотов
func generateVirtualSlots(master models.Master, service models.Service, bookings []models.Booking, date time.Time) []Slot {
    moscow, _ := time.LoadLocation("Europe/Moscow")

    slotMap := make(map[int64]string)

    workStart := parseTimeOnDate(date, master.WorkHours.StartTime, moscow)
    workEnd := parseTimeOnDate(date, master.WorkHours.EndTime, moscow)
    lunchStart := parseTimeOnDate(date, master.WorkHours.LunchStart, moscow)
    lunchEnd := parseTimeOnDate(date, master.WorkHours.LunchEnd, moscow)

    serviceDuration := time.Duration(service.DurationMin) * time.Minute

    // ШАГ 1: Основная сетка
    for current := workStart; current.Add(serviceDuration).Before(workEnd) || current.Add(serviceDuration).Equal(workEnd); current = current.Add(serviceDuration) {
        slotEnd := current.Add(serviceDuration)
        status := "free"

        // Проверка обеда
        if master.WorkHours.LunchStart != "" && master.WorkHours.LunchEnd != "" {
            if isIntersecting(current, slotEnd, lunchStart, lunchEnd) {
                status = "booked"
            }
        }

        // Проверка пересечений с записями
        if status == "free" {
            for _, booking := range bookings {
                if isIntersecting(current, slotEnd, booking.StartTime, booking.EndTime) {
                    status = "booked"
                    break
                }
            }
        }

        key := current.Unix()
        slotMap[key] = status
    }

    // ШАГ 2: "Умные слоты" после каждой записи
    for _, booking := range bookings {
        if booking.Status == models.BookingStatusCancelledByClient || booking.Status == models.BookingStatusCancelledNoShow {
            continue
        }

        freeAfter := booking.EndTime
        smartSlotEnd := freeAfter.Add(serviceDuration)

        if smartSlotEnd.After(workEnd) {
            continue
        }

        if master.WorkHours.LunchStart != "" && master.WorkHours.LunchEnd != "" {
            if isIntersecting(freeAfter, smartSlotEnd, lunchStart, lunchEnd) {
                continue
            }
        }

        isFree := true
        for _, otherBooking := range bookings {
            if otherBooking.ID == booking.ID {
                continue
            }
            if isIntersecting(freeAfter, smartSlotEnd, otherBooking.StartTime, otherBooking.EndTime) {
                isFree = false
                break
            }
        }

        if isFree {
            key := freeAfter.Unix()
            if _, exists := slotMap[key]; !exists {
                slotMap[key] = "free"
            }
        }
    }

    // ШАГ 3: Сортируем слоты
    type slotEntry struct {
        timestamp int64
        status    string
    }
    var entries []slotEntry
    for ts, status := range slotMap {
        entries = append(entries, slotEntry{timestamp: ts, status: status})
    }
    sort.Slice(entries, func(i, j int) bool {
        return entries[i].timestamp < entries[j].timestamp
    })

    // ШАГ 4: Формируем финальный массив
    var slots []Slot
    for _, entry := range entries {
        slots = append(slots, Slot{
            StartTime: time.Unix(entry.timestamp, 0).In(moscow),
            Status:    entry.status,
        })
    }

    return slots
}

func isIntersecting(start1, end1, start2, end2 time.Time) bool {
    return start1.Before(end2) && start2.Before(end1)
}

func parseTimeOnDate(date time.Time, timeStr string, loc *time.Location) time.Time {
    if timeStr == "" {
        return time.Time{}
    }
    t, _ := time.Parse("15:04", timeStr)
    return time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), 0, 0, loc)
}