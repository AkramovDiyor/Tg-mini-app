package httpapi

import (
	// "backend/pkg/telegram"
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

// Кастомный тип для ключа, чтобы безопасно работать с контекстом запроса
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
                tgID = 777111222
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

func NewRouter(handler *BookingHandler, tgBotToken string) *chi.Mux {
	r := chi.NewRouter()



    corsHandler := cors.New(cors.Options{
        AllowedOrigins: []string{
            "http://localhost:5173", 
			"http://127.0.0.1:5173",
        },
        AllowedMethods: []string{
            http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
        },
        AllowedHeaders: []string{
            "Origin",
			"Accept",
			"Content-Type",
			"X-Requested-With",
			"X-Telegram-Init-Data",
        },
        AllowCredentials: true,
        Debug: false,
    })


    r.Use(corsHandler.Handler)

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Группа эндпоинтов V1 (для чистоты архитектуры)
	r.Route("/api/v1", func(v1 chi.Router) {

		// Публичные маршруты (доступны без авторизации через ТГ)
		v1.Get("/master/{invite_link}/services", handler.GetServices)
		v1.Get("/master/{invite_link}/slots", handler.GetSlots)

		// Защищенные маршруты (требуют заголовок X-Telegram-Init-Data)
		v1.Group(func(protected chi.Router) {
			protected.Use(AuthMiddleware(tgBotToken))

			protected.Post("/book", handler.BookSlot)
            protected.Get("/client/bookings", handler.GetClientBookings)
		})
	})

	return r
}
