package httpapi

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

func NewRouter(
	bookingHandler *BookingHandler,
	masterHandler *MasterHandler,
	photoHandler *PhotoHandler,
	meHandler *MeHandler,
	tgBotToken string,
) *chi.Mux {
	r := chi.NewRouter()

	corsHandler := cors.New(cors.Options{
		// 🔥 Разрешаем ЛЮБЫЕ домены через функцию
		AllowOriginRequestFunc: func(r *http.Request, origin string) bool {
			// Локальные домены
			if strings.HasPrefix(origin, "http://localhost:") ||
				strings.HasPrefix(origin, "https://localhost:") ||
				strings.HasPrefix(origin, "http://127.0.0.1:") ||
				strings.HasPrefix(origin, "http://192.168.") ||
				strings.HasPrefix(origin, "https://192.168.") {
				return true
			}
			// 🎯 КЛЮЧЕВОЕ: разрешаем ВСЕ ngrok, vercel, loca.lt
			if strings.HasSuffix(origin, ".ngrok-free.app") ||
				strings.HasSuffix(origin, ".ngrok.io") ||
				strings.HasSuffix(origin, ".vercel.app") ||
				strings.HasSuffix(origin, ".loca.lt") ||
				strings.HasSuffix(origin, ".trycloudflare.com") {
				return true
			}
			return false
		},
		AllowedMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH",
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
			"X-Requested-With",
			"X-Telegram-Init-Data",
			"ngrok-skip-browser-warning",
			"Bypass-Tunnel-Reminder",
			"Cache-Control",
		},
		ExposedHeaders: []string{
			"Content-Length",
			"Content-Type",
		},
		AllowCredentials: true,
		MaxAge:           300,
		Debug:            true, // 🔥 Включено для отладки
	})

	r.Use(corsHandler.Handler)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// 🔥 ИСПРАВЛЕНО: поднимаемся на уровень выше
	staticDir := filepath.Join(".", "..", "static", "uploads")
	staticDir, _ = filepath.Abs(staticDir)

	r.Handle("/static/uploads/*", http.StripPrefix("/static/uploads/", http.FileServer(http.Dir(staticDir))))
	// Группа эндпоинтов V1
	r.Route("/api/v1", func(v1 chi.Router) {

		// ============================================
		// 🔓 ПУБЛИЧНЫЕ МАРШРУТЫ (для клиентов)
		// Используют префикс /invite/ чтобы не конфликтовать
		// ============================================
		v1.Get("/invite/{invite_link}/services", bookingHandler.GetServices)
		v1.Get("/invite/{invite_link}/slots", bookingHandler.GetSlots)
		v1.Get("/invite/{invite_link}/info", bookingHandler.GetInfoMaster)

		    v1.Group(func(auth chi.Router) {
        auth.Use(AuthMiddleware(tgBotToken))
        auth.Get("/me", meHandler.GetMe)
    })

		// ============================================
		// 🔐 ЗАЩИЩЕННЫЕ МАРШРУТЫ КЛИЕНТА
		// ============================================
		v1.Group(func(client chi.Router) {
			client.Use(AuthMiddleware(tgBotToken))

			client.Post("/book", bookingHandler.BookSlot)
			client.Get("/client/bookings", bookingHandler.GetClientBookings)
			client.Post("/client/bookings/{booking_id}/cancel", bookingHandler.CancelBooking)
		})

		// ============================================
		// 🔐 ЗАЩИЩЕННЫЕ МАРШРУТЫ МАСТЕРА
		// ============================================
		v1.Route("/master", func(m chi.Router) {
			m.Use(AuthMiddleware(tgBotToken))

			m.Get("/today", masterHandler.GetToday)
			m.Get("/waitlist", masterHandler.GetWaitlist)
			m.Get("/profile", masterHandler.GetProfile)
			m.Put("/profile", masterHandler.UpdateProfile)
			m.Get("/services", masterHandler.GetServices)
			m.Post("/services", masterHandler.CreateService)
			m.Put("/services/{service_id}", masterHandler.UpdateService)
			m.Delete("/services/{service_id}", masterHandler.DeleteService)
			m.Put("/settings", masterHandler.UpdateSettings)

			m.Post("/photos", photoHandler.UploadPhoto)
			m.Get("/photos", photoHandler.GetPhotos)
			m.Delete("/photos/{photo_id}", photoHandler.DeletePhoto)
		})
	})

	return r
}
