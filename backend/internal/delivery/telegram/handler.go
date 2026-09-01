package telegram

import (
	"backend/internal/services"
	"context"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	bot           *tgbotapi.BotAPI
	masterService *services.MasterService
	webAppURL     string // 🔥 Теперь это https://t.me/book1ngMin1App_bot/slotify
	botUsername   string
}

func NewHandler(masterService *services.MasterService, webAppURL string, bot *tgbotapi.BotAPI) *Handler {
	return &Handler{
		bot:           bot,
		masterService: masterService,
		webAppURL:     webAppURL,
		botUsername:   bot.Self.UserName,
	}
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

	startParam := ""
	args := msg.CommandArguments()
	if args != "" {
		startParam = strings.TrimSpace(args)
	}

	ctx := context.Background()

	if startParam != "" && startParam != "master" {
		h.handleClientFlow(ctx, chatID, firstName, startParam)
		return
	}

	h.handleMasterFlow(ctx, chatID, firstName)
}

func (h *Handler) handleMasterFlow(ctx context.Context, chatID int64, firstName string) {
	master, err := h.masterService.RegisterMaster(ctx, chatID, firstName)
	if err != nil {
		log.Printf("❌ Ошибка регистрации мастера: %v", err)
		h.sendTextMessage(chatID, "⚠️ Ошибка регистрации. Попробуйте позже.")
		return
	}

	log.Printf("✅ Мастер: ID=%d, invite_link=%s", master.ID, master.InviteLink)

	// 🔥 Ссылка для клиентов (Named Web App с параметром)
	clientLink := fmt.Sprintf("%s?startapp=%s", h.webAppURL, master.InviteLink)
	
	// 🔥 Ссылка для панели мастера (без параметра)
	masterLink := h.webAppURL

	welcomeText := fmt.Sprintf(
		"👋 Привет, *%s*!\n\n"+
			"🔗 *Твоя ссылка для клиентов:*\n"+
			"`%s`\n\n"+
			"Отправляй её клиентам. При переходе они сразу увидят форму записи!\n\n"+
			"👇 Нажми кнопку ниже, чтобы открыть свою панель управления:",
		firstName, clientLink,
	)

	// 🔥 Кнопка для мастера
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("🛠 Открыть панель управления", masterLink),
		),
	)

	reply := tgbotapi.NewMessage(chatID, welcomeText)
	reply.ParseMode = tgbotapi.ModeMarkdown
	reply.ReplyMarkup = keyboard

	if _, err := h.bot.Send(reply); err != nil {
		log.Printf("❌ Ошибка отправки: %v", err)
	}
}

func (h *Handler) handleClientFlow(ctx context.Context, chatID int64, firstName string, inviteLink string) {
	log.Printf("👤 Клиент %d открывает ссылку мастера: %s", chatID, inviteLink)

	// 🔥 Ссылка для клиента (Named Web App с параметром)
	clientLink := fmt.Sprintf("%s?startapp=%s", h.webAppURL, inviteLink)

	welcomeText := fmt.Sprintf(
		"👋 Привет, *%s*!\n\n"+
			"Нажми кнопку ниже, чтобы записаться к мастеру.",
		firstName,
	)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("✂️ Записаться", clientLink),
		),
	)

	reply := tgbotapi.NewMessage(chatID, welcomeText)
	reply.ParseMode = tgbotapi.ModeMarkdown
	reply.ReplyMarkup = keyboard

	if _, err := h.bot.Send(reply); err != nil {
		log.Printf("❌ Ошибка отправки: %v", err)
	}
}

func (h *Handler) sendTextMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("❌ Ошибка отправки: %v", err)
	}
}