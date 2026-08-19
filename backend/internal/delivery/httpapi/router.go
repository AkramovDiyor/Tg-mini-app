package httpapi

import (
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

func NewRouter(
	bookingHandler *BookingHandler,
	masterHandler *MasterHandler,
	photoHandler *PhotoHandler,
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
		Debug:            false,
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
