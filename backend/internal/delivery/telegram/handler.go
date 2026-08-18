package telegram

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	bot         *tgbotapi.BotAPI
	masterRepo  repositories.MasterRepository
	webAppURL   string
	botUsername string // Инициализируется позже, в StartBot
}

// 🔥 ИСПРАВЛЕНО: Убираем botUsername из конструктора
func NewHandler(masterRepo repositories.MasterRepository, webAppURL string) *Handler {
	return &Handler{
		bot:        nil, // Будет установлен в StartBot
		masterRepo: masterRepo,
		webAppURL:  webAppURL,
	}
}

// SetBot устанавливает бота после инициализации
func (h *Handler) SetBot(bot *tgbotapi.BotAPI) {
	h.bot = bot
	h.botUsername = bot.Self.UserName
}

func (h *Handler) HandleMessage(update *tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	msg := update.Message
	chatID := msg.Chat.ID
	firstName := msg.From.FirstName

	if !msg.IsCommand() || msg.Command() != "start" {
		h.sendTextMessage(chatID, "🤖 Я тебя не понимаю. Используй команду /start для начала работы.")
		return
	}

	log.Printf("📩 Получена команда /start от пользователя: %d (%s)", chatID, firstName)

	ctx := context.Background()
	master, err := h.masterRepo.GetMasterByTelegramID(ctx, chatID)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			master, err = h.createNewMaster(ctx, chatID, firstName)
			if err != nil {
				log.Printf("❌ Ошибка создания мастера: %v", err)
				h.sendTextMessage(chatID, "⚠️ Произошла ошибка при регистрации. Попробуйте позже.")
				return
			}
			log.Printf("✅ Создан новый мастер: ID=%d, invite_link=%s", master.ID, master.InviteLink)
		} else {
			log.Printf("❌ Ошибка БД: %v", err)
			h.sendTextMessage(chatID, "⚠️ Ошибка подключения к базе данных. Попробуйте позже.")
			return
		}
	} else {
		log.Printf("✅ Найден существующий мастер: ID=%d, invite_link=%s", master.ID, master.InviteLink)
	}

	welcomeText := fmt.Sprintf(
		"👋 Привет, *%s*!\n\n"+
			"Ты успешно зарегистрирован в системе.\n\n"+
			"🔗 *Твоя ссылка для клиентов:*\n"+
			"`https://t.me/%s?startapp=%s`\n\n"+
			"Отправь её клиентам, чтобы они могли записаться к тебе!",
		firstName, h.botUsername, master.InviteLink,
	)





	reply := tgbotapi.NewMessage(chatID, welcomeText)
	reply.ParseMode = tgbotapi.ModeMarkdown
	

	if _, err := h.bot.Send(reply); err != nil {
		log.Printf("❌ Ошибка отправки сообщения: %v", err)
	}
}

func (h *Handler) createNewMaster(ctx context.Context, chatID int64, firstName string) (models.Master, error) {
	inviteLink := generateInviteLink(8)

	newMaster := models.Master{
		TelegramID: chatID,
		Name:       firstName,
		Bio:        "",
		Address:    "",
		InviteLink: inviteLink,
	}

	err := h.masterRepo.CreateMaster(ctx, newMaster)
	if err != nil {
		return models.Master{}, fmt.Errorf("failed to create master: %w", err)
	}

	return h.masterRepo.GetMasterByTelegramID(ctx, chatID)
}

func (h *Handler) sendTextMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("❌ Ошибка отправки сообщения: %v", err)
	}
}

func generateInviteLink(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[r.Intn(len(charset))]
	}
	return string(result)
}