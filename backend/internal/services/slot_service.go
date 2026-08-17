package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"context"
	"fmt"
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
	master, err := s.MasterRepo.GetMasterByInviteLink(ctx, inviteLink)
	if err != nil {
		return nil, err
	}

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

	dayName := getRussianWeekday(date.Weekday())
	isWorkingDay := false
	for _, d := range master.WorkHours.WorkingDays {
		if d == dayName {
			isWorkingDay = true
			break
		}
	}
	if !isWorkingDay {
		return []Slot{}, nil
	}

	bookings, err := s.BookingRepo.GetBookingsByMasterAndDate(ctx, master.ID, date)
	if err != nil {
		return nil, err
	}

	return generateVirtualSlots(master, service, bookings, date), nil
}

func generateVirtualSlots(master models.Master, service models.Service, bookings []models.Booking, date time.Time) []Slot {
	// Используем map чтобы избежать дубликатов слотов
	slotMap := make(map[int64]string) // key = Unix timestamp, value = status

	workStart := parseTimeOnDate(date, master.WorkHours.StartDay)
	workEnd := parseTimeOnDate(date, master.WorkHours.EndDay)
	breakStart := parseTimeOnDate(date, master.WorkHours.BreakStart)
	breakEnd := parseTimeOnDate(date, master.WorkHours.BreakEnd)

	serviceDuration := time.Duration(service.DurationMin) * time.Minute
	
	// ==========================================
	// ШАГ 1: Основная сетка (шаг = длительность услуги)
	// ==========================================
	for current := workStart; current.Add(serviceDuration).Before(workEnd) || current.Add(serviceDuration).Equal(workEnd); current = current.Add(serviceDuration) {
		slotEnd := current.Add(serviceDuration)
		status := "free"

		// Проверка обеда
		if master.WorkHours.BreakStart != "" && master.WorkHours.BreakEnd != "" {
			if isIntersecting(current, slotEnd, breakStart, breakEnd) {
				status = "booked"
			}
		}

		// Проверка пересечений с записями
		if status == "free" {
			for _, booking := range bookings {
				if isIntersecting(current, slotEnd, booking.StartTime.UTC(), booking.EndTime.UTC()) {
					status = "booked"
					break
				}
			}
		}

		key := current.Unix()
		slotMap[key] = status
	}

	// ==========================================
	// ШАГ 2: "Умные слоты" после каждой записи
	// Заполняем зазоры между записями
	// ==========================================
	for _, booking := range bookings {
		// Пропускаем отмененные записи
		if booking.Status == models.BookingStatusCancelledByClient || booking.Status == models.BookingStatusCancelledNoShow {
			continue
		}

		// Время, когда мастер освободится после этой записи
		freeAfter := booking.EndTime.UTC()
		smartSlotEnd := freeAfter.Add(serviceDuration)

		// Проверяем, что умный слот вписывается в рабочий день
		if smartSlotEnd.After(workEnd) {
			continue
		}

		// Проверяем, что умный слот не попадает на обед
		if master.WorkHours.BreakStart != "" && master.WorkHours.BreakEnd != "" {
			if isIntersecting(freeAfter, smartSlotEnd, breakStart, breakEnd) {
				continue
			}
		}

		// Проверяем, что умный слот не пересекается с другими записями
		isFree := true
		for _, otherBooking := range bookings {
			if otherBooking.ID == booking.ID {
				continue // Пропускаем саму себя
			}
			if isIntersecting(freeAfter, smartSlotEnd, otherBooking.StartTime.UTC(), otherBooking.EndTime.UTC()) {
				isFree = false
				break
			}
		}

		if isFree {
			key := freeAfter.Unix()
			// Добавляем только если такого слота еще нет
			if _, exists := slotMap[key]; !exists {
				slotMap[key] = "free"
			}
		}
	}

	// ==========================================
	// ШАГ 3: Сортируем слоты по времени
	// ==========================================
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

	// ==========================================
	// ШАГ 4: Формируем финальный массив
	// ==========================================
	var slots []Slot
	for _, entry := range entries {
		slots = append(slots, Slot{
			StartTime: time.Unix(entry.timestamp, 0).UTC(),
			Status:    entry.status,
		})
	}

	return slots
}

func isIntersecting(start1, end1, start2, end2 time.Time) bool {
	return start1.Before(end2) && start2.Before(end1)
}

func parseTimeOnDate(date time.Time, timeStr string) time.Time {
	if timeStr == "" {
		return time.Time{}
	}
	t, _ := time.Parse("15:04", timeStr)
	return time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC)
}

func getRussianWeekday(w time.Weekday) string {
	days := map[time.Weekday]string{
		time.Monday: "Пн", time.Tuesday: "Вт", time.Wednesday: "Ср",
		time.Thursday: "Чт", time.Friday: "Пт", time.Saturday: "Сб", time.Sunday: "Вс",
	}
	return days[w]
}