package httpapi

import (
	// "backend/pkg/telegram"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)


func NewRouter(
    bookingHandler *BookingHandler, 
    masterHandler *MasterHandler, 
    tgBotToken string,
) *chi.Mux {
	r := chi.NewRouter()



    corsHandler := cors.New(cors.Options{
        AllowedOrigins: []string{
            "http://localhost:5173", 
			"http://127.0.0.1:5173",
            "http://10.21.33.135",
            "https://tg-mini-app-pink.vercel.app",
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
		v1.Get("/master/{invite_link}/services", bookingHandler.GetServices)
		v1.Get("/master/{invite_link}/slots", bookingHandler.GetSlots)

		// Защищенные маршруты (требуют заголовок X-Telegram-Init-Data)
		v1.Group(func(protected chi.Router) {
			protected.Use(AuthMiddleware(tgBotToken))

			protected.Post("/book", bookingHandler.BookSlot)
            protected.Get("/client/bookings", bookingHandler.GetClientBookings)
            protected.Post("/client/bookings/{booking_id}/cancel", bookingHandler.CancelBooking)
		})

        v1.Route("/master", func(m chi.Router) {
            m.Use(AuthMiddleware(tgBotToken)) // Тоже защищена!
            
            m.Get("/today", masterHandler.GetToday)
            m.Get("/waitlist", masterHandler.GetWaitlist)
            m.Get("/profile", masterHandler.GetProfile)
            m.Put("/profile", masterHandler.UpdateProfile)
            m.Get("/services", masterHandler.GetServices)
            m.Post("/services", masterHandler.CreateService)
            m.Put("/settings", masterHandler.UpdateSettings)
        })
	})

	return r
}
