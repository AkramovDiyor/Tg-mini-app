package httpapi

import (
	"context"
	"net/http"
)

type contextKey string

const TgIDKey contextKey = "tg_id"

// AuthMiddleware — перехватывает заголовок, проверяет его и кладёт Telegram ID в контекст
// func AuthMiddleware(tgBotToken string) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			// 1. Достаем заголовок X-Telegram-Init-Data из HTTP-запроса
// 			initData := r.Header.Get("X-Telegram-Init-Data")
// 			if initData == "" {
// 				http.Error(w, "Unauthorized: Missing Telegram Init Data", http.StatusUnauthorized)
// 				return
// 			}

//             tgID, err := telegram.ValidateInitData(initData, tgBotToken)
//             if err != nil {
//                 http.Error(w, "Unauthorized: Invalid Signature", http.StatusUnauthorized)
//                 return
//             }

// 			// 3. Создаем новый контекст на основе старого и кладем туда наш Telegram ID
// 			ctx := context.WithValue(r.Context(), TgIDKey, tgID)

// 			// 4. Пропускаем запрос дальше по цепочке, передавая обновленный контекст
// 			next.ServeHTTP(w, r.WithContext(ctx))
// 		})
// 	}
// }

func AuthMiddleware(tgBotToken string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            initData := r.Header.Get("X-Telegram-Init-Data")
            if initData == "" {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }

            // ВРЕМЕННО ДЛЯ ПОСТМАНА:
            var tgID int64
            if initData == "test-vasya" {
                tgID = 111222333
            } else {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }

            // РЕАЛЬНУЮ КРИПТУ ВОТ ТУТ ВРЕМЕННО ОТКЛЮЧИ:
            /*
            tgID, err := telegram.ValidateInitData(initData, tgBotToken)
            if err != nil {
                http.Error(w, "Unauthorized: Invalid Signature", http.StatusUnauthorized)
                return
            }
            */
            // И исправь имя переменной ниже с err на обычное присвоение, если компилятор ругается.
            ctx := context.WithValue(r.Context(), TgIDKey, tgID)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}