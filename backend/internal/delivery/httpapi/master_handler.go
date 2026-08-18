package httpapi

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

// NewMasterHandler — конструктор для создания обработчика
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

    // Получаем записи на сегодня
    today := time.Now()
    bookings, err := h.bookingRepo.GetBookingsByMasterAndDate(r.Context(), master.ID, today)
    if err != nil {
        http.Error(w, fmt.Sprintf("failed to fetch bookings: %v", err), http.StatusInternalServerError)
        return
    }

    // 🔥 Считаем статистику
    totalRevenue := 0
    for _, booking := range bookings {
        totalRevenue += booking.PriceLocked
    }

    response := map[string]interface{}{
        "stats": map[string]int{
            "count": len(bookings),
            "total": totalRevenue,  // 🔥 Теперь будет правильная сумма
        },
        "schedule": bookings,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (h *MasterHandler) GetWaitlist(w http.ResponseWriter, r *http.Request) {
    // 1. Получаем tg_id из контекста
    tgID, err := GetIDFromContext(r.Context())
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    // 2. Находим мастера по tg_id
    master, err := h.masterRepo.GetMasterByTelegramID(r.Context(), tgID)
    if err != nil {
        http.Error(w, "Master not found", http.StatusNotFound)
        return
    }

    // 3. Получаем лист ожидания для мастера
    waitlist, err := h.waitlistRepo.GetWaitlistByMaster(r.Context(), master.ID)
    if err != nil {
        http.Error(w, "Failed to get waitlist", http.StatusInternalServerError)
        return
    }

    // 4. Отправляем JSON ответ (даже если список пустой)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    
    // Если список пустой - вернется [] (пустой массив)
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
    // 1. Получаем tg_id и мастера (как у тебя уже есть)
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

    // 2. Читаем и декодируем JSON из тела запроса
    var req UpdateProfileRequest
    
    // Важно: закрываем тело после чтения
    defer r.Body.Close()
    
    // Декодируем JSON
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON format", http.StatusBadRequest)
        return
    }

    // 3. Теперь у нас есть данные в req
    // req.Name, req.Bio, req.Address

    // 4. Обновляем профиль
    err = h.masterRepo.UpdateMasterProfile(r.Context(), master.ID, req.Name, req.Bio, req.Address)
    if err != nil {
        http.Error(w, "Failed to update profile", http.StatusInternalServerError)
        return
    }

    // 5. Возвращаем успешный ответ
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

    // Возвращаем успешный ответ с созданной услугой
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated) // 201 Created
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status":  "success",
        "message": "Service created successfully",
        "service": service,
    })

}

func (h *MasterHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	// TODO: Реализовать обновление настроек
}
