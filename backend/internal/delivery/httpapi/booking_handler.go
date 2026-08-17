package httpapi

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"backend/internal/services"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

func GetClientIDFromContext(ctx context.Context) (int64, error) {
    tgID, ok := ctx.Value(TgIDKey).(int64)
    if !ok {
        return 0, fmt.Errorf("telegram ID not found in context")
    }
    return tgID, nil
}

type BookingHandler struct {
    masterRepo  repositories.MasterRepository
    serviceRepo repositories.ServiceRepository
    slotRepo    repositories.SlotRepository
    bookingRepo repositories.BookingRepository // ДОБАВЛЕНО
    bookingServ *services.BookingService
    slotServ    *services.SlotService
}

func NewBookingHandler(
    masterRepo repositories.MasterRepository,
    serviceRepo repositories.ServiceRepository,
    slotRepo repositories.SlotRepository,
    bookingRepo repositories.BookingRepository, // ДОБАВЛЕНО
    bookingServ *services.BookingService,
    slotServ *services.SlotService,
) *BookingHandler {
    return &BookingHandler{
        masterRepo:  masterRepo,
        serviceRepo: serviceRepo,
        slotRepo:    slotRepo,
        bookingRepo: bookingRepo, // ДОБАВЛЕНО
        bookingServ: bookingServ,
        slotServ:    slotServ,
    }
}

func (h *BookingHandler) GetServices(w http.ResponseWriter, r *http.Request) {
    // ... (код без изменений)
    inviteLink := chi.URLParam(r, "invite_link")
    if inviteLink == "" {
        http.Error(w, "bad request: missing invite link", http.StatusBadRequest)
        return
    }
    master, err := h.masterRepo.GetMasterByInviteLink(r.Context(), inviteLink)
    if err != nil {
        http.Error(w, "master not found", http.StatusNotFound)
        return
    }
    serviceList, err := h.serviceRepo.GetServicesByMasterID(r.Context(), master.ID)
    if err != nil {
        http.Error(w, "failed to fetch services", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(serviceList)
}

// ПЕРЕПИСАННЫЙ МЕТОД
func (h *BookingHandler) GetSlots(w http.ResponseWriter, r *http.Request) {
    inviteLink := chi.URLParam(r, "invite_link")
    dateStr := r.URL.Query().Get("date")
    serviceIDStr := r.URL.Query().Get("service_id")

    if inviteLink == "" || dateStr == "" || serviceIDStr == "" {
        http.Error(w, "bad request: missing parameters (invite_link, date, service_id)", http.StatusBadRequest)
        return
    }

    serviceID, err := strconv.ParseInt(serviceIDStr, 10, 64)
    if err != nil {
        http.Error(w, "bad request: invalid service_id", http.StatusBadRequest)
        return
    }

    targetDate, err := time.Parse("2006-01-02", dateStr)
    if err != nil {
        http.Error(w, "bad request: invalid date format, use YYYY-MM-DD", http.StatusBadRequest)
        return
    }

    // Вызываем сервис генерации виртуальных слотов
    slotsList, err := h.slotServ.GetAvailableSlots(r.Context(), inviteLink, targetDate, serviceID)
    if err != nil {
        http.Error(w, fmt.Sprintf("failed to fetch slots: %v", err), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(slotsList)
}

// ... (BookSlot и GetInfoMaster остаются без изменений)
type BookRequest struct {
    StartTime string `json:"start_time"` // ISO 8601: "2026-08-20T10:00:00Z"
    ServiceID int64  `json:"service_id"`
    Name      string `json:"name"`
    Price     int    `json:"price"`
}

func (h *BookingHandler) BookSlot(w http.ResponseWriter, r *http.Request) {
    tgID, err := GetClientIDFromContext(r.Context())
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    var req BookRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "bad request: invalid JSON body", http.StatusBadRequest)
        return
    }

    if req.StartTime == "" || req.ServiceID == 0 || req.Name == "" {
        http.Error(w, "bad request: missing required fields", http.StatusBadRequest)
        return
    }

    // Парсим время начала
    startTime, err := time.Parse(time.RFC3339, req.StartTime)
    if err != nil {
        http.Error(w, "bad request: invalid start_time format", http.StatusBadRequest)
        return
    }

    // Вызываем сервис бронирования
    err = h.bookingServ.BookSlot(r.Context(), tgID, req.ServiceID, req.Name, req.Price, startTime)
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnprocessableEntity)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    w.Write([]byte(`{"status":"success","message":"slot booked successfully"}`))
}

func (h *BookingHandler) GetInfoMaster(w http.ResponseWriter, r *http.Request) (models.Master, error) {
    // ... (код без изменений)
    inviteLink := chi.URLParam(r, "invite_link")
    if inviteLink == "" {
        http.Error(w, "bad request: missing invite link", http.StatusBadRequest)
    }
    master, err := h.masterRepo.GetMasterByInviteLink(r.Context(), inviteLink)
    if err != nil {
        http.Error(w, "master not found", http.StatusNotFound)
    }
    return master, err
}



// Добавь этот метод в конец файла booking_handler.go

func (h *BookingHandler) GetClientBookings(w http.ResponseWriter, r *http.Request) {
	// 1. Достаем Telegram ID из контекста
	tgID, err := GetClientIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized: missing client id", http.StatusUnauthorized)
		return
	}

	// 2. Вызываем ПРАВИЛЬНЫЙ метод
	bookings, err := h.bookingRepo.GetBookingsByClientTgID(r.Context(), tgID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to fetch client bookings: %v", err), http.StatusInternalServerError)
		return
	}

	// 3. Отдаем JSON фронтенду
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bookings)
}