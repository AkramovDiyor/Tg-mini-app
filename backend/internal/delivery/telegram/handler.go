package telegram

import (
	"backend/internal/models"
	"backend/internal/services"
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	bot           *tgbotapi.BotAPI
	masterService *services.MasterService // 🔥 Используем сервис вместо репозитория
	webAppURL     string
	botUsername   string
}

// 🔥 ИСПРАВЛЕНО: принимаем MasterService вместо MasterRepository
func NewHandler(masterService *services.MasterService, webAppURL string) *Handler {
	return &Handler{
		bot:           nil,
		masterService: masterService,
		webAppURL:     webAppURL,
	}
}

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
	
	// 🔥 Используем сервис с идемпотентной регистрацией
	master, err := h.masterService.RegisterMaster(ctx, chatID, firstName)
	if err != nil {
		log.Printf("❌ Ошибка регистрации мастера: %v", err)
		h.sendTextMessage(chatID, "⚠️ Произошла ошибка при регистрации. Попробуйте позже.")
		return
	}

	log.Printf("✅ Мастер готов: ID=%d, invite_link=%s", master.ID, master.InviteLink)

	// Формируем приветственное сообщение
	welcomeText := fmt.Sprintf(
		"👋 Привет, *%s*!\n\n"+
			"Ты успешно зарегистрирован в системе.\n\n"+
			"🔗 *Твоя ссылка для клиентов:*\n"+
			"`https://t.me/%s?startapp=%s`\n\n"+
			"Отправь её клиентам, чтобы они могли записаться к тебе!",
		firstName, h.botUsername, master.InviteLink,
	)

	// 🔥 Создаем кнопку Web App
	webAppURLWithParam := fmt.Sprintf("%s?startapp=master", h.webAppURL)
	
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{
				Text: "🛠 Открыть панель управления",
				WebApp: &tgbotapi.WebAppInfo{
					URL: webAppURLWithParam,
				},
			},
		),
	)

	reply := tgbotapi.NewMessage(chatID, welcomeText)
	reply.ParseMode = tgbotapi.ModeMarkdown
	reply.ReplyMarkup = keyboard

	if _, err := h.bot.Send(reply); err != nil {
		log.Printf("❌ Ошибка отправки сообщения: %v", err)
	}
}

func (h *Handler) sendTextMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("❌ Ошибка отправки сообщения: %v", err)
	}
}