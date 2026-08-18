package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	TgBotToken string
	WebAppURL  string // 🔥 НОВОЕ
}

func Load() (Config, error) {
	if err := godotenv.Load("../../.env"); err != nil {
		// .env может отсутствовать в контейнере
	}

	cfg := Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		TgBotToken: os.Getenv("TG_BOT_TOKEN"),
		WebAppURL:  os.Getenv("WEB_APP_URL"), // 🔥 НОВОЕ
	}

	if cfg.DBHost == "" || cfg.DBUser == "" || cfg.DBName == "" {
		return Config{}, fmt.Errorf("пропущены обязательные настройки БД")
	}

	if cfg.TgBotToken == "" {
		return Config{}, fmt.Errorf("пропущен TG_BOT_TOKEN")
	}

	if cfg.WebAppURL == "" {
		log.Println("⚠️  WEB_APP_URL не задан, Web App кнопка не будет работать")
	}

	return cfg, nil
}