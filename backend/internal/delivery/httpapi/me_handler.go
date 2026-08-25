package httpapi

import (
	"backend/internal/repositories"
	"encoding/json"
	"log"
	"net/http"
)

type MeHandler struct {
	masterRepo repositories.MasterRepository
}

func NewMeHandler(masterRepo repositories.MasterRepository) *MeHandler {
	return &MeHandler{masterRepo: masterRepo}
}

// GetMe определяет, кто текущий пользователь: мастер или клиент
func (h *MeHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	// 1. Достаем tgID из контекста (уже валидирован через AuthMiddleware)
	tgID, err := GetIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	log.Printf("🔍 Определяем роль для tgID=%d", tgID)

	// 2. Проверяем, есть ли такой мастер
	isMaster, err := h.masterRepo.ExistsByTelegramID(r.Context(), tgID)
	if err != nil {
		log.Printf("❌ Ошибка проверки мастера: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var response map[string]interface{}

	if isMaster {
		// 3a. Пользователь — мастер, достаем полные данные
		master, err := h.masterRepo.GetMasterByTelegramID(r.Context(), tgID)
		if err != nil {
			log.Printf("❌ Не удалось получить мастера: %v", err)
			http.Error(w, "Failed to fetch master data", http.StatusInternalServerError)
			return
		}

		response = map[string]interface{}{
			"role":   "master",
			"master": master,
		}
		log.Printf("✅ Определен МАСТЕР: %s (invite_link=%s)", master.Name, master.InviteLink)
	} else {
		// 3b. Пользователь — клиент
		response = map[string]interface{}{
			"role": "client",
			"user": map[string]interface{}{
				"telegram_id": tgID,
			},
		}
		log.Printf("✅ Определен КЛИЕНТ: tgID=%d", tgID)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}