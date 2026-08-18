package telegram

import (
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// 🔥 ИСПРАВЛЕНО: StartBot теперь вызывает SetBot
func StartBot(token string, handler *Handler) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	bot.Debug = false

	log.Printf("✅ Бот авторизован как @%s", bot.Self.UserName)

	// 🔥 Устанавливаем бота в handler
	handler.SetBot(bot)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	log.Println("🚀 Telegram-бот запущен и слушает обновления...")

	for update := range updates {
		go func(upd tgbotapi.Update) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("💥 Паника в обработчике бота: %v", r)
				}
			}()

			if upd.Message != nil {
				handler.HandleMessage(&upd)
			}
		}(update)
	}

	return nil
}

func StartBotWithRetry(token string, handler *Handler) {
	for {
		err := StartBot(token, handler)
		if err != nil {
			log.Printf("❌ Ошибка запуска бота: %v. Повторная попытка через 5 секунд...", err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
}