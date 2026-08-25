package telegram

import (
	"backend/internal/services"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ============================================================
// Кастомные структуры для Mini App кнопок
// (т.к. старая версия библиотеки не поддерживает WebApp поле)
// ============================================================

type WebAppInfo struct {
	URL string `json:"url"`
}

type InlineKeyboardButton struct {
	Text   string        `json:"text"`
	WebApp *WebAppInfo   `json:"web_app,omitempty"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

// ============================================================
// Handler
// ============================================================

type Handler struct {
	bot           *tgbotapi.BotAPI
	masterService *services.MasterService
	webAppURL     string
	botUsername   string
}

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
		h.sendTextMessage(chatID, "🤖 Используй команду /start для начала работы.")
		return
	}

	// 🔥 КЛЮЧЕВОЕ: достаем параметр из /start <параметр>
	startParam := ""
	args := msg.CommandArguments()
	if args != "" {
		startParam = strings.TrimSpace(args)
		log.Printf("📩 /start с параметром: %s от пользователя %d (%s)", startParam, chatID, firstName)
	} else {
		log.Printf("📩 /start БЕЗ параметра от пользователя %d (%s)", chatID, firstName)
	}

	ctx := context.Background()

	// ============================================================
	// СЦЕНАРИЙ A: Клиент открыл ссылку мастера (есть параметр)
	// ============================================================
	if startParam != "" && startParam != "master" {
		h.handleClientFlow(ctx, chatID, firstName, startParam)
		return
	}

	// ============================================================
	// СЦЕНАРИЙ B: Мастер зашёл в бота (нет параметра или "master")
	// ============================================================
	h.handleMasterFlow(ctx, chatID, firstName)
}

// handleMasterFlow — обработка входа мастера
func (h *Handler) handleMasterFlow(ctx context.Context, chatID int64, firstName string) {
	master, err := h.masterService.RegisterMaster(ctx, chatID, firstName)
	if err != nil {
		log.Printf("❌ Ошибка регистрации мастера: %v", err)
		h.sendTextMessage(chatID, "⚠️ Ошибка регистрации. Попробуйте позже.")
		return
	}

	log.Printf("✅ Мастер: ID=%d, invite_link=%s", master.ID, master.InviteLink)

	clientLink := fmt.Sprintf("https://t.me/%s?startapp=%s", h.botUsername, master.InviteLink)
	masterWebAppURL := fmt.Sprintf("%s?mode=master", h.webAppURL)

	welcomeText := fmt.Sprintf(
		"👋 Привет, *%s*!\n\n"+
			"🔗 *Твоя ссылка для клиентов:*\n"+
			"`%s`\n\n"+
			"Отправляй её клиентам, чтобы они могли записаться к тебе!",
		firstName, clientLink,
	)

	// 🔥 Формируем клавиатуру через кастомные структуры
	keyboard := InlineKeyboardMarkup{
		InlineKeyboard: [][]InlineKeyboardButton{
			{
				{
					Text:   "🛠 Открыть панель управления",
					WebApp: &WebAppInfo{URL: masterWebAppURL},
				},
			},
		},
	}

	h.sendWithKeyboard(chatID, welcomeText, keyboard)
}

// handleClientFlow — обработка входа клиента по ссылке мастера
func (h *Handler) handleClientFlow(ctx context.Context, chatID int64, firstName string, inviteLink string) {
	log.Printf("👤 Клиент %d открывает ссылку мастера: %s", chatID, inviteLink)

	clientWebAppURL := fmt.Sprintf("%s?startapp=%s", h.webAppURL, inviteLink)

	welcomeText := fmt.Sprintf(
		"👋 Привет, *%s*!\n\n"+
			"Нажми кнопку ниже, чтобы записаться к мастеру.",
		firstName,
	)

	// 🔥 Формируем клавиатуру через кастомные структуры
	keyboard := InlineKeyboardMarkup{
		InlineKeyboard: [][]InlineKeyboardButton{
			{
				{
					Text:   "✂️ Записаться",
					WebApp: &WebAppInfo{URL: clientWebAppURL},
				},
			},
		},
	}

	h.sendWithKeyboard(chatID, welcomeText, keyboard)
}

// sendWithKeyboard — универсальный метод отправки сообщения с кастомной клавиатурой
func (h *Handler) sendWithKeyboard(chatID int64, text string, keyboard InlineKeyboardMarkup) {
	keyboardJSON, err := json.Marshal(keyboard)
	if err != nil {
		log.Printf("❌ Ошибка маршалинга клавиатуры: %v", err)
		h.sendTextMessage(chatID, text) // fallback без клавиатуры
		return
	}

	params := tgbotapi.Params{
		"chat_id":      strconv.FormatInt(chatID, 10),
		"text":         text,
		"parse_mode":   "Markdown",
		"reply_markup": string(keyboardJSON),
	}

	_, err = h.bot.MakeRequest("sendMessage", params)
	if err != nil {
		log.Printf("❌ Ошибка отправки сообщения с клавиатурой: %v", err)
	}
}

// sendTextMessage — простое текстовое сообщение
func (h *Handler) sendTextMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("❌ Ошибка отправки: %v", err)
	}
}