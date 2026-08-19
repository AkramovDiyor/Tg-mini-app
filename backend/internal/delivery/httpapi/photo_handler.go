package httpapi

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PhotoHandler struct {
	photoRepo  repositories.PhotoRepository
	masterRepo repositories.MasterRepository
	uploadDir  string
}

func NewPhotoHandler(photoRepo repositories.PhotoRepository, masterRepo repositories.MasterRepository) *PhotoHandler {
    // 🔥 ИСПРАВЛЕНО: поднимаемся на уровень выше (в корень backend)
    workDir, _ := os.Getwd()
    uploadDir := filepath.Join(workDir, "..", "static", "uploads")
    
    // Нормализуем путь (убираем ..)
    uploadDir, _ = filepath.Abs(uploadDir)
    
    log.Printf("📁 Upload directory: %s", uploadDir)
    
    return &PhotoHandler{
        photoRepo:  photoRepo,
        masterRepo: masterRepo,
        uploadDir:  uploadDir,
    }
}

func (h *PhotoHandler) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	log.Println("📸 UploadPhoto: начало обработки")

	// 1. Авторизация
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

	// 2. Парсинг multipart формы (лимит 10 МБ)
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("❌ ParseMultipartForm error: %v", err)
		http.Error(w, "File too large (max 10MB)", http.StatusBadRequest)
		return
	}

	// 3. Получение файла из формы
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("❌ FormFile error: %v", err)
		http.Error(w, "Failed to get file from request. Make sure field name is 'file'", http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Printf("📎 Получен файл: %s, размер: %d байт", header.Filename, header.Size)

	// 4. Проверка расширения
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
	}
	if !allowedExts[ext] {
		http.Error(w, fmt.Sprintf("Invalid file type: %s. Allowed: jpg, jpeg, png, gif, webp", ext), http.StatusBadRequest)
		return
	}

	// 5. Генерация имени
	newFilename := uuid.New().String() + ext

	// 6. 🔥 Создание папки с подробным логированием
	log.Printf("📂 Создаем папку: %s", h.uploadDir)
	if err := os.MkdirAll(h.uploadDir, 0755); err != nil {
		log.Printf("❌ MkdirAll error: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create upload directory: %v", err), http.StatusInternalServerError)
		return
	}

	// 7. Проверка прав на запись (пробуем создать тестовый файл)
	testFile := filepath.Join(h.uploadDir, ".test")
	if f, err := os.Create(testFile); err != nil {
		log.Printf("❌ Нет прав на запись в %s: %v", h.uploadDir, err)
		http.Error(w, fmt.Sprintf("No write permission in upload directory: %v", err), http.StatusInternalServerError)
		return
	} else {
		f.Close()
		os.Remove(testFile)
	}

	// 8. Создание целевого файла
	dstPath := filepath.Join(h.uploadDir, newFilename)
	log.Printf("💾 Создаем файл: %s", dstPath)
	
	dst, err := os.Create(dstPath)
	if err != nil {
		log.Printf("❌ os.Create error: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create file: %v", err), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// 9. Копирование содержимого
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("❌ ReadAll error: %v", err)
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	if _, err := dst.Write(fileBytes); err != nil {
		log.Printf("❌ Write error: %v", err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Файл сохранен: %s (%d байт)", dstPath, len(fileBytes))

	// 10. Сохранение в БД
	url := "/static/uploads/" + newFilename
	photoID, err := h.photoRepo.AddPhoto(r.Context(), master.ID, url)
	if err != nil {
		os.Remove(dstPath)
		log.Printf("❌ AddPhoto error: %v", err)
		http.Error(w, fmt.Sprintf("Failed to save photo to database: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Фото сохранено в БД с ID: %d", photoID)

	// 11. Успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":  photoID,
		"url": url,
	})
}

// GetPhotos и DeletePhoto без изменений...
func (h *PhotoHandler) GetPhotos(w http.ResponseWriter, r *http.Request) {
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

	photos, err := h.photoRepo.GetPhotosByMaster(r.Context(), master.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get photos: %v", err), http.StatusInternalServerError)
		return
	}

	if photos == nil {
		photos = []models.MasterPhoto{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(photos)
}

func (h *PhotoHandler) DeletePhoto(w http.ResponseWriter, r *http.Request) {
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

	photoIDStr := chi.URLParam(r, "photo_id")
	photoID, err := strconv.ParseInt(photoIDStr, 10, 64)
	if err != nil || photoID <= 0 {
		http.Error(w, "bad request: invalid photo_id", http.StatusBadRequest)
		return
	}

	err = h.photoRepo.DeletePhoto(r.Context(), photoID, master.ID)
	if err != nil {
		if err.Error() == "photo not found or not owned by master" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to delete photo: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}