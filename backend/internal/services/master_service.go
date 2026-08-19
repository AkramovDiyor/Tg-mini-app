package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

type MasterService struct {
	masterRepo repositories.MasterRepository
}

// NewMasterService — конструктор сервиса мастера
func NewMasterService(masterRepo repositories.MasterRepository) *MasterService {
	return &MasterService{
		masterRepo: masterRepo,
	}
}

// RegisterMaster — идемпотентная регистрация мастера
// Если мастер уже существует — возвращаем его данные
// Если не существует — создаем нового с уникальной ссылкой
func (s *MasterService) RegisterMaster(ctx context.Context, tgID int64, firstName string) (models.Master, error) {
	// 1. Пытаемся найти существующего мастера
	existing, err := s.masterRepo.GetMasterByTelegramID(ctx, tgID)
	if err == nil {
		// Мастер найден — возвращаем его (идемпотентность!)
		log.Printf("✅ Найден существующий мастер: ID=%d, invite_link=%s", existing.ID, existing.InviteLink)
		return existing, nil
	}

	// 2. Проверяем, что ошибка именно "not found"
	if !strings.Contains(err.Error(), "not found") {
		// Какая-то другая ошибка (например, отвал БД)
		return models.Master{}, fmt.Errorf("failed to check existing master: %w", err)
	}

	// 3. Мастер не найден — создаем нового
	log.Printf("🆕 Мастер %d не найден, создаем нового", tgID)

	// 4. Генерируем уникальную ссылку
	inviteLink, err := s.generateUniqueInviteLink(ctx)
	if err != nil {
		return models.Master{}, err
	}

	// 5. Создаем структуру мастера
	newMaster := models.Master{
		TelegramID: tgID,
		Name:       firstName,
		Bio:        "",
		Address:    "",
		InviteLink: inviteLink,
	}

	// 6. Сохраняем в БД
	err = s.masterRepo.CreateMaster(ctx, newMaster)
	if err != nil {
		return models.Master{}, fmt.Errorf("failed to create master: %w", err)
	}

	// 7. Достаем созданного мастера (с ID из БД)
	created, err := s.masterRepo.GetMasterByTelegramID(ctx, tgID)
	if err != nil {
		return models.Master{}, fmt.Errorf("master created but failed to retrieve: %w", err)
	}

	log.Printf("✅ Создан новый мастер: ID=%d, invite_link=%s", created.ID, created.InviteLink)
	return created, nil
}

// generateUniqueInviteLink — генерирует уникальную 8-символьную ссылку
// С проверкой на коллизии (если вдруг такая ссылка уже есть — генерируем новую)
func (s *MasterService) generateUniqueInviteLink(ctx context.Context) (string, error) {
	const maxAttempts = 10

	for attempt := 0; attempt < maxAttempts; attempt++ {
		link := generateRandomString(8)

		// Проверяем, что ссылка еще не занята
		_, err := s.masterRepo.GetMasterByInviteLink(ctx, link)
		if err != nil {
			// Если ошибка "not found" — ссылка свободна, возвращаем её
			if strings.Contains(err.Error(), "not found") {
				return link, nil
			}
			// Другая ошибка — возвращаем её
			return "", fmt.Errorf("failed to check link uniqueness: %w", err)
		}

		// Ссылка уже занята, пробуем снова
		log.Printf("⚠️ Ссылка %s уже существует, генерируем новую (попытка %d/%d)", link, attempt+1, maxAttempts)
	}

	return "", fmt.Errorf("failed to generate unique invite link after %d attempts", maxAttempts)
}

// generateRandomString — генерирует случайную строку заданной длины
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[r.Intn(len(charset))]
	}
	return string(result)
}