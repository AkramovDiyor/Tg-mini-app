package database

import (
    "context"
    "fmt"
    "time"

    "backend/internal/config"

    "github.com/jackc/pgx/v5" 
    "github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres(cfg config.Config) (*pgxpool.Pool, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // 1. Парсим конфиг из строки подключения
    pgxConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
    if err != nil {
        return nil, fmt.Errorf("не удалось распарсить DATABASE_URL: %w", err)
    }

    // 2. ОТКЛЮЧАЕМ КЭШ ПОДГОТОВЛЕННЫХ ЗАПРОСОВ (ФИКС ОШИБКИ SUPABASE)
    pgxConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeExec

    // 3. Создаем пул с обновленным конфигом
    pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
    if err != nil {
        return nil, fmt.Errorf("не удалось создать конфигурацию пула Postgres: %w", err)
    }

    // Делаем обязательный Ping
    if err := pool.Ping(ctx); err != nil {
        pool.Close()
        return nil, fmt.Errorf("ошибка проверки связи (Ping) с PostgreSQL: %w", err)
    }

    return pool, nil
}