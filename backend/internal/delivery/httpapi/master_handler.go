package httpapi

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type UpdateProfileRequest struct {
    Name    string `json:"name"`
    Bio     string `json:"bio"`
    Address string `json:"address"`
}

type CreateServiceRequest struct {
    Name        string `json:"name"`
    DurationMin int    `json:"duration_min"`
    Price       int    `json:"price"`
}

type UpdateSettingsRequest struct {
    WorkHours map[string]interface{} `json:"work_hours"`
    Settings  map[string]interface{} `json:"settings"`
}

type MasterHandler struct {
	masterRepo   repositories.MasterRepository
	serviceRepo  repositories.ServiceRepository
	bookingRepo  repositories.BookingRepository
	waitlistRepo repositories.WaitlistRepository
}

func NewMasterHandler(
	masterRepo repositories.MasterRepository,
	serviceRepo repositories.ServiceRepository,
	bookingRepo repositories.BookingRepository,
	waitlistRepo repositories.WaitlistRepository,
) *MasterHandler {
	return &MasterHandler{
		masterRepo:   masterRepo,
		serviceRepo:  serviceRepo,
		bookingRepo:  bookingRepo,
		waitlistRepo: waitlistRepo,
	}
}

func (h *MasterHandler) GetToday(w http.ResponseWriter, r *http.Request) {
    tgID, err := GetIDFromContext(r.Context())
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    master, err := h.masterRepo.GetMasterByTelegramID(r.Context(), tgID)
    if err != nil {
        http.Error(w, "Master not found", http.StatusNotFound)
        return
    }

    // 🔥 Получаем текущую дату в Москве
    moscow, err := time.LoadLocation("Europe/Moscow")
    if err != nil {
        http.Error(w, "failed to load timezone", http.StatusInternalServerError)
        return
    }

    now := time.Now().In(moscow)
    
    // Нормализуем к началу дня в Москве
    todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, moscow)

    // Получаем записи на сегодня
    bookings, err := h.bookingRepo.GetBookingsByMasterAndDate(r.Context(), master.ID, todayStart)
    if err != nil {
        http.Error(w, fmt.Sprintf("failed to fetch bookings: %v", err), http.StatusInternalServerError)
        return
    }

    // Считаем статистику
    totalRevenue := 0
    for _, booking := range bookings {
        totalRevenue += booking.PriceLocked
    }

    response := map[string]interface{}{
        "stats": map[string]int{
            "count": len(bookings),
            "total": totalRevenue,
        },
        "schedule": bookings,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (h *MasterHandler) GetWaitlist(w http.ResponseWriter, r *http.Request) {
    tgID, err := GetIDFromContext(r.Context())
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    master, err := h.masterRepo.GetMasterByTelegramID(r.Context(), tgID)
    if err != nil {
        http.Error(w, "Master not found", http.StatusNotFound)
        return
    }

    waitlist, err := h.waitlistRepo.GetWaitlistByMaster(r.Context(), master.ID)
    if err != nil {
        http.Error(w, "Failed to get waitlist", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    
    if err := json.NewEncoder(w).Encode(waitlist); err != nil {
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
        return
    }
}

func (h *MasterHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	tgID, err := GetIDFromContext(r.Context())
	if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

	master, err := h.masterRepo.GetMasterByTelegramID(r.Context(), tgID)
	if err != nil {
        http.Error(w, "Master not found", http.StatusNotFound)
        return
    }

	json.NewEncoder(w).Encode(master)
}

func (h *MasterHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
    tgID, err := GetIDFromContext(r.Context())
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    master, err := h.masterRepo.GetMasterByTelegramID(r.Context(), tgID)
    if err != nil {
        http.Error(w, "Master not found", http.StatusNotFound)
        return
    }

    var req UpdateProfileRequest
    
    defer r.Body.Close()
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON format", http.StatusBadRequest)
        return
    }

    err = h.masterRepo.UpdateMasterProfile(r.Context(), master.ID, req.Name, req.Bio, req.Address)
    if err != nil {
        http.Error(w, "Failed to update profile", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *MasterHandler) GetServices(w http.ResponseWriter, r *http.Request) {
	tgID, err := GetIDFromContext(r.Context())
	if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

	master, err := h.masterRepo.GetMasterByTelegramID(r.Context(), tgID)
	if err != nil {
        http.Error(w, "Master not found", http.StatusNotFound)
        return
    }

	service, err := h.serviceRepo.GetServicesByMasterID(r.Context(), master.ID)
	if err != nil {
        http.Error(w, "service not found", http.StatusNotFound)
        return
    }

	json.NewEncoder(w).Encode(service)
}

func (h *MasterHandler) CreateService(w http.ResponseWriter, r *http.Request) {
	tgID, err := GetIDFromContext(r.Context())
	if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

	master, err := h.masterRepo.GetMasterByTelegramID(r.Context(), tgID)
	if err != nil {
        http.Error(w, "Master not found", http.StatusNotFound)
        return
    }

	var req CreateServiceRequest

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
        return
	}

	service := models.Service{
		MasterID: master.ID,
		Name: req.Name,
		DurationMin: req.DurationMin,
		Price: req.Price,
	}
	err = h.serviceRepo.CreateService(r.Context(), service)
    if err != nil {
        http.Error(w, "Failed to create service", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status":  "success",
        "message": "Service created successfully",
        "service": service,
    })
}

func (h *MasterHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
    tgID, err := GetIDFromContext(r.Context())
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    master, err := h.masterRepo.GetMasterByTelegramID(r.Context(), tgID)
    if err != nil {
        http.Error(w, "Master not found", http.StatusNotFound)
        return
    }

    var req UpdateSettingsRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON format", http.StatusBadRequest)
        return
    }

    // 🔥 Преобразуем map в структуры
    workHours, err := mapToWorkHours(req.WorkHours)
    if err != nil {
        http.Error(w, fmt.Sprintf("Invalid work_hours format: %v", err), http.StatusBadRequest)
        return
    }

    settings, err := mapToSettings(req.Settings)
    if err != nil {
        http.Error(w, fmt.Sprintf("Invalid settings format: %v", err), http.StatusBadRequest)
        return
    }

    // Обновляем в базе
    err = h.masterRepo.UpdateMasterSettings(r.Context(), master.ID, workHours, settings)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to update settings: %v", err), http.StatusInternalServerError)
        return
    }

    log.Printf("✅ Настройки обновлены для мастера %d: WorkDays=%v, StartTime=%s", 
        master.ID, workHours.WorkDays, workHours.StartTime)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// 🔥 Вспомогательные функции для преобразования map в структуры
func mapToWorkHours(data map[string]interface{}) (models.WorkHours, error) {
    var wh models.WorkHours

    if workDays, ok := data["work_days"].([]interface{}); ok {
        wh.WorkDays = make([]bool, len(workDays))
        for i, v := range workDays {
            if b, ok := v.(bool); ok {
                wh.WorkDays[i] = b
            }
        }
    }

    if startTime, ok := data["start_time"].(string); ok {
        wh.StartTime = startTime
    }
    if endTime, ok := data["end_time"].(string); ok {
        wh.EndTime = endTime
    }
    if lunchStart, ok := data["lunch_start"].(string); ok {
        wh.LunchStart = lunchStart
    }
    if lunchEnd, ok := data["lunch_end"].(string); ok {
        wh.LunchEnd = lunchEnd
    }

    return wh, nil
}

func mapToSettings(data map[string]interface{}) (models.MasterSettings, error) {
    var s models.MasterSettings

    if autoCancel, ok := data["auto_cancel"].(bool); ok {
        s.AutoCancel = autoCancel
    }
    if cancelHours, ok := data["cancel_hours"].(string); ok {
        s.CancelHours = cancelHours
    }
    if offerWaitlist, ok := data["offer_waitlist"].(bool); ok {
        s.OfferWaitlist = offerWaitlist
    }

    return s, nil
}



// UpdateServiceRequest описывает тело запроса для обновления услуги
type UpdateServiceRequest struct {
    Name        string `json:"name"`
    DurationMin int    `json:"duration_min"`
    Price       int    `json:"price"`
}

func (h *MasterHandler) UpdateService(w http.ResponseWriter, r *http.Request) {
    // 1. Получаем Telegram ID мастера из контекста
    tgID, err := GetIDFromContext(r.Context())
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    // 2. Находим мастера в базе
    master, err := h.masterRepo.GetMasterByTelegramID(r.Context(), tgID)
    if err != nil {
        http.Error(w, "Master not found", http.StatusNotFound)
        return
    }

    // 3. 🔥 ИСПРАВЛЕНО: правильно достаем service_id из URL
    serviceIDStr := chi.URLParam(r, "service_id")  // ← ВАЖНО: передать r первым аргументом
    serviceID, err := strconv.ParseInt(serviceIDStr, 10, 64)
    if err != nil || serviceID <= 0 {
        http.Error(w, "bad request: invalid service_id", http.StatusBadRequest)
        return
    }

    // 4. Парсим JSON из тела запроса
    var req UpdateServiceRequest
    defer r.Body.Close()
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON format", http.StatusBadRequest)
        return
    }

    // 5. Валидация полей
    if req.Name == "" || req.DurationMin <= 0 || req.Price <= 0 {
        http.Error(w, "bad request: name, duration_min and price are required and must be positive", http.StatusBadRequest)
        return
    }

    // 6. Вызываем метод репозитория (с проверкой master_id для безопасности)
    err = h.serviceRepo.UpdateService(r.Context(), serviceID, master.ID, req.Name, req.DurationMin, req.Price)
    if err != nil {
        // Разные типы ошибок — разные HTTP-коды
        if err.Error() == "service not found or not owned by master" {
            http.Error(w, err.Error(), http.StatusForbidden) // 403 — либо чужая, либо не существует
            return
        }
        http.Error(w, fmt.Sprintf("Failed to update service: %v", err), http.StatusInternalServerError)
        return
    }

    // 7. Возвращаем успешный ответ
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status":  "success",
        "message": "Service updated successfully",
    })
}

func (h *MasterHandler) DeleteService(w http.ResponseWriter, r *http.Request) {
    tgID, err := GetIDFromContext(r.Context())
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    master, err := h.masterRepo.GetMasterByTelegramID(r.Context(), tgID)
    if err != nil {
        http.Error(w, "Master not found", http.StatusNotFound)
        return
    }

    serviceIDStr := chi.URLParam(r, "service_id")
    serviceID, err := strconv.ParseInt(serviceIDStr, 10, 64)
    if err != nil || serviceID <= 0 {
        http.Error(w, "bad request: invalid service_id", http.StatusBadRequest)
        return
    }

    // Мягкое удаление (архивация)
    err = h.serviceRepo.DeleteService(r.Context(), serviceID, master.ID)
    if err != nil {
        if err.Error() == "service not found, not owned by master, or already deleted" {
            http.Error(w, err.Error(), http.StatusForbidden)
            return
        }
        http.Error(w, fmt.Sprintf("Failed to delete service: %v", err), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status":  "deleted",
        "message": "Service archived successfully",
    })
}