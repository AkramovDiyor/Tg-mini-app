package httpapi

import (
	"backend/pkg/telegram"
	"context"
	"fmt"
	"log"
	"net/http"
)

type contextKey string

const (
	TgIDKey         contextKey = "tg_id"
	StartParamKey   contextKey = "start_param"
	TelegramUserKey contextKey = "telegram_user"
)

func GetIDFromContext(ctx context.Context) (int64, error) {
	tgID, ok := ctx.Value(TgIDKey).(int64)
	if !ok {
		return 0, fmt.Errorf("telegram ID not found in context")
	}
	return tgID, nil
}

func GetStartParamFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(StartParamKey).(string); ok {
		return v
	}
	return ""
}

func GetTelegramUserFromContext(ctx context.Context) telegram.WebAppUser {
	if user, ok := ctx.Value(TelegramUserKey).(telegram.WebAppUser); ok {
		return user
	}
	return telegram.WebAppUser{}
}

func AuthMiddleware(tgBotToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			initData := r.Header.Get("X-Telegram-Init-Data")
			if initData == "" {
				http.Error(w, "Missing Init Data", http.StatusUnauthorized)
				return
			}

			var tgID int64
			var startParam string
			var tgUser telegram.WebAppUser

			if initData == "test-vasya" {
				tgID = 777111222
				startParam = r.URL.Query().Get("startapp")
				tgUser = telegram.WebAppUser{ID: 777111222, FirstName: "Тестовый", LastName: "Клиент"}
			} else if initData == "test-master" {
				tgID = 999999
				startParam = "master"
				tgUser = telegram.WebAppUser{ID: 999999, FirstName: "Тестовый", LastName: "Мастер"}
			} else {
				var err error
				tgUser, startParam, err = telegram.ValidateInitData(initData, tgBotToken)
				if err != nil {
					log.Printf("❌ Ошибка валидации: %v", err)
					http.Error(w, "Invalid Signature", http.StatusUnauthorized)
					return
				}
				tgID = tgUser.ID
				log.Printf("✅ Авторизация: tgID=%d, name=%s %s", tgID, tgUser.FirstName, tgUser.LastName)
			}

			// 🔥 ИСПРАВЛЕНО: используем := для первого присвоения ctx
			ctx := context.WithValue(r.Context(), TgIDKey, tgID)
			ctx = context.WithValue(ctx, StartParamKey, startParam)
			ctx = context.WithValue(ctx, TelegramUserKey, tgUser)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}