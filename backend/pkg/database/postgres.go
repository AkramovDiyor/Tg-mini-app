package database

import (
	"context"
	"fmt"
	"time"

	"backend/internal/config" // Внимание: убедитесь, что имя модуля в go.mod совпадает с 'backend'

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgres инициализирует пул соединений с базой данных PostgreSQL
func NewPostgres(cfg config.Config) (*pgxpool.Pool, error) {
	// Сборка DSN строки подключения
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	// Создаем контекст с таймаутом на случай, если база "лежит" и долго не отвечает
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Инициализируем пул соединений (соединение ленивое, физически база еще не проверяется)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать конфигурацию пула Postgres: %w", err)
	}

	// Делаем обязательный Ping, чтобы физически проверить доступность СУБД
	if err := pool.Ping(ctx); err != nil {
		pool.Close() // Закрываем пул, если подключиться не удалось
		return nil, fmt.Errorf("ошибка проверки связи (Ping) с PostgreSQL: %w", err)
	}

	return pool, nil
}
