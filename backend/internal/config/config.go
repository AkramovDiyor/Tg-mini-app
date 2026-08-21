package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// type Config struct {
// 	DBHost     string
// 	DBPort     string
// 	DBUser     string
// 	DBPassword string
// 	DBName     string
// 	TgBotToken string
// 	WebAppURL  string // 🔥 НОВОЕ
// }

type Config struct {
	DatabaseURL string
	TgBotToken  string
}

func Load() (Config, error) {
	if err := godotenv.Load("../../.env"); err != nil {
		// .env может отсутствовать в контейнере
	}

	cfg := Config{
		// DBHost:     os.Getenv("DB_HOST"),
		// DBPort:     os.Getenv("DB_PORT"),
		// DBUser:     os.Getenv("DB_USER"),
		// DBPassword: os.Getenv("DB_PASSWORD"),
		// DBName:     os.Getenv("DB_NAME"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		TgBotToken:  os.Getenv("TG_BOT_TOKEN"),
		// WebAppURL:  os.Getenv("WEB_APP_URL"), // 🔥 НОВОЕ
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("пропущены обязательные настройки БД (DATABASE_URL)")
	}

	// if cfg.TgBotToken == "" {
	// 	return Config{}, fmt.Errorf("пропущен TG_BOT_TOKEN")
	// }

	return cfg, nil
}
