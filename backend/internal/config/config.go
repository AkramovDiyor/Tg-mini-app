package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	TgBotToken  string
	WebAppURL   string // 🔥 ВОЗВРАЩАЕМ
}

func Load() (Config, error) {
	if err := godotenv.Load("../../.env"); err != nil {
		// .env может отсутствовать в контейнере
	}

	cfg := Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		TgBotToken:  os.Getenv("TG_BOT_TOKEN"),
		WebAppURL:   os.Getenv("WEB_APP_URL"), // 🔥 ВОЗВРАЩАЕМ
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("пропущены обязательные настройки БД (DATABASE_URL)")
	}

	if cfg.TgBotToken == "" {
		return Config{}, fmt.Errorf("пропущен TG_BOT_TOKEN")
	}

	return cfg, nil
}