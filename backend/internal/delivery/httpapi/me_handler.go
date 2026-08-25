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
	// 1. Достаем tgID и start_param из контекста
	tgID, err := GetIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	startParam := GetStartParamFromContext(r.Context())
	log.Printf("🔍 Определяем роль для tgID=%d, start_param=%s", tgID, startParam)

	// 2. Проверяем, есть ли такой мастер
	isMaster, err := h.masterRepo.ExistsByTelegramID(r.Context(), tgID)
	if err != nil {
		log.Printf("❌ Ошибка проверки мастера: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var response map[string]interface{}

	if isMaster {
		// 3a. Пользователь — мастер
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
		// 🔥 Возвращаем invite_link из start_param
		response = map[string]interface{}{
			"role": "client",
			"user": map[string]interface{}{
				"telegram_id": tgID,
			},
			"invite_link": startParam, // 🔥 КЛЮЧЕВОЕ: клиент узнает, к какому мастеру пришел
		}
		log.Printf("✅ Определен КЛИЕНТ: tgID=%d, invite_link=%s", tgID, startParam)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}