package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config содержит все переменные окружения для работы приложения
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	TgBotToken string
}

// Load загружает файл .env и собирает структуру Config
func Load() (Config, error) {
	// Пытаемся загрузить .env файл.
	// Если его нет (например, в Docker-контейнере), godotenv вернет ошибку,
	// но переменные могут быть в самой системе, поэтому ошибку логируем, но не падаем.
	if err := godotenv.Load("../../.env"); err != nil {
		// Если нужно жестко требовать .env, можно раскомментировать строку ниже:
		// return Config{}, fmt.Errorf("ошибка загрузки .env файла: %w", err)
	}

	cfg := Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		TgBotToken: os.Getenv("TG_BOT_TOKEN"),
	}

	// Минимальная валидация критических полей
	if cfg.DBHost == "" || cfg.DBUser == "" || cfg.DBName == "" {
		return Config{}, fmt.Errorf("пропущены обязательные настройки базы данных в окружении")
	}
	// if cfg.TgBotToken == "" {
	// 	return Config{}, fmt.Errorf("пропущен токен Telegram-бота (TG_BOT_TOKEN)")
	// }

	return cfg, nil
}
