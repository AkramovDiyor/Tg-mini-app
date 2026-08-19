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

func GetIDFromContext(ctx context.Context) (int64, error) {
    tgID, ok := ctx.Value(TgIDKey).(int64)
    if !ok {
        return 0, fmt.Errorf("telegram ID not found in context")
    }
    return tgID, nil
}


// MasterInfoResponse — агрегированный ответ с инфо мастера и его фото
type MasterInfoResponse struct {
    ID         int64                `json:"id"`
    Name       string               `json:"name"`
    Bio        string               `json:"bio"`
    Address    string               `json:"address"`
    InviteLink string               `json:"invite_link"`
    Photos     []models.MasterPhoto `json:"photos"`
}

type BookingHandler struct {
    masterRepo  repositories.MasterRepository
    serviceRepo repositories.ServiceRepository
    slotRepo    repositories.SlotRepository
    bookingRepo repositories.BookingRepository
    photoRepo   repositories.PhotoRepository
    bookingServ *services.BookingService
    slotServ    *services.SlotService
}

func NewBookingHandler(
    masterRepo repositories.MasterRepository,
    serviceRepo repositories.ServiceRepository,
    slotRepo repositories.SlotRepository,
    bookingRepo repositories.BookingRepository,
    photoRepo   repositories.PhotoRepository,
    bookingServ *services.BookingService,
    slotServ *services.SlotService,
) *BookingHandler {
    return &BookingHandler{
        masterRepo:  masterRepo,
        serviceRepo: serviceRepo,
        slotRepo:    slotRepo,
        bookingRepo: bookingRepo,
        photoRepo:   photoRepo,
        bookingServ: bookingServ,
        slotServ:    slotServ,
    }
}

func (h *BookingHandler) GetServices(w http.ResponseWriter, r *http.Request) {
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

    // 🔥 Парсим дату с учетом московского времени
    moscow, err := time.LoadLocation("Europe/Moscow")
    if err != nil {
        http.Error(w, "failed to load timezone", http.StatusInternalServerError)
        return
    }

    targetDate, err := time.ParseInLocation("2006-01-02", dateStr, moscow)
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

type BookRequest struct {
    StartTime string `json:"start_time"`
    ServiceID int64  `json:"service_id"`
    Name      string `json:"name"`
    Price     int    `json:"price"`
}

func (h *BookingHandler) BookSlot(w http.ResponseWriter, r *http.Request) {
    tgID, err := GetIDFromContext(r.Context())
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

    // Парсим время начала (оно приходит в UTC формате от фронтенда)
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

func (h *BookingHandler) GetInfoMaster(w http.ResponseWriter, r *http.Request) {
    // 1. Достаем invite_link из URL
    inviteLink := chi.URLParam(r, "invite_link")
    if inviteLink == "" {
        http.Error(w, "bad request: missing invite link", http.StatusBadRequest)
        return
    }

    // 2. Ищем мастера в базе
    master, err := h.masterRepo.GetMasterByInviteLink(r.Context(), inviteLink)
    if err != nil {
        http.Error(w, "master not found", http.StatusNotFound)
        return
    }

    // 3. 🔥 АГРЕГАЦИЯ: достаем фото мастера
    photos, err := h.photoRepo.GetPhotosByMaster(r.Context(), master.ID)
    if err != nil {
        // Не фатальная ошибка — просто логируем и отдаем мастера без фото
        photos = []models.MasterPhoto{}
    }
    
    // Защита от null в JSON
    if photos == nil {
        photos = []models.MasterPhoto{}
    }

    // 4. Формируем агрегированный ответ
    response := MasterInfoResponse{
        ID:         master.ID,
        Name:       master.Name,
        Bio:        master.Bio,
        Address:    master.Address,
        InviteLink: master.InviteLink,
        Photos:     photos,
    }

    // 5. Отдаем JSON
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

func (h *BookingHandler) GetClientBookings(w http.ResponseWriter, r *http.Request) {
    tgID, err := GetIDFromContext(r.Context())
    if err != nil {
        http.Error(w, "Unauthorized: missing client id", http.StatusUnauthorized)
        return
    }

    bookings, err := h.bookingRepo.GetBookingsByClientTgID(r.Context(), tgID)
    if err != nil {
        http.Error(w, fmt.Sprintf("failed to fetch client bookings: %v", err), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(bookings)
}

func (h *BookingHandler) CancelBooking(w http.ResponseWriter, r *http.Request) {
    tgID, err := GetIDFromContext(r.Context())
    if err != nil {
        http.Error(w, "Unauthorized: missing client id", http.StatusUnauthorized)
        return
    }

    bookingIDStr := chi.URLParam(r, "booking_id")
    bookingID, err := strconv.ParseInt(bookingIDStr, 10, 64)
    if err != nil || bookingID <= 0 {
        http.Error(w, "bad request: invalid booking_id", http.StatusBadRequest)
        return
    }

    err = h.bookingServ.CancelBooking(r.Context(), bookingID, tgID)
    if err != nil {
        switch err.Error() {
        case "запись не найдена":
            http.Error(w, err.Error(), http.StatusNotFound)
        case "вы не можете отменить чужую запись":
            http.Error(w, err.Error(), http.StatusForbidden)
        case "запись уже отменена", "нельзя отменить завершенную запись", "нельзя отменить запись, которая уже началась":
            http.Error(w, err.Error(), http.StatusUnprocessableEntity)
        default:
            http.Error(w, fmt.Sprintf("failed to cancel booking: %v", err), http.StatusInternalServerError)
        }
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status":  "success",
        "message": "запись успешно отменена",
    })
}