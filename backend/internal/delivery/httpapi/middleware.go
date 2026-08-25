package httpapi

import (
	"backend/pkg/telegram"
	"context"
	"log"
	"net/http"
)

type contextKey string

const (
	TgIDKey       contextKey = "tg_id"
	StartParamKey contextKey = "start_param" // 🔥 НОВОЕ
)


// 🔥 НОВАЯ ФУНКЦИЯ: достаем start_param из контекста
func GetStartParamFromContext(ctx context.Context) string {
	startParam, ok := ctx.Value(StartParamKey).(string)
	if !ok {
		return ""
	}
	return startParam
}


func AuthMiddleware(tgBotToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			initData := r.Header.Get("X-Telegram-Init-Data")
			if initData == "" {
				http.Error(w, "Unauthorized: Missing Telegram Init Data", http.StatusUnauthorized)
				return
			}

			var tgID int64
			var startParam string

			// 🔥 ДЛЯ РАЗРАБОТКИ: тестовые токены
			if initData == "test-vasya" {
				tgID = 777111222
				// 🔥 ИСПРАВЛЕНО: берем startapp прямо из URL запроса!
				startParam = r.URL.Query().Get("startapp")
				if startParam == "" {
					startParam = "test-default-link" 
				}
				log.Printf("🧪 Тестовая авторизация: vasya (tgID=%d, start_param=%s)", tgID, startParam)
				
			} else if initData == "test-master" {
				tgID = 999999
				startParam = "master"
				log.Printf("🧪 Тестовая авторизация: master (tgID=%d)", tgID)
				
			} else {
				// 🔥 НАСТОЯЩАЯ ВАЛИДАЦИЯ (в Telegram)
				var err error
				tgID, startParam, err = telegram.ValidateInitData(initData, tgBotToken)
				if err != nil {
					log.Printf("❌ Ошибка валидации initData: %v", err)
					http.Error(w, "Unauthorized: Invalid Signature", http.StatusUnauthorized)
					return
				}
				log.Printf("✅ Реальная авторизация: tgID=%d, start_param=%s", tgID, startParam)
			}

			ctx := context.WithValue(r.Context(), TgIDKey, tgID)
			ctx = context.WithValue(ctx, StartParamKey, startParam)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}