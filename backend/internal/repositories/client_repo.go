package repositories

import (
	"backend/internal/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClientRepository interface {
	GetOrCreateClient(ctx context.Context, telegramID int64, username string, firstName string) (models.Client, error)
	GetClientByTelegramID(ctx context.Context, telegramID int64) (models.Client, error)
}

type ClientRepo struct {
	db *pgxpool.Pool
}

func NewClientRepo(db *pgxpool.Pool) *ClientRepo {
	return &ClientRepo{
		db: db,
	}
}

func (c *ClientRepo) GetOrCreateClient(ctx context.Context, telegramID int64, username string, firstName string) (models.Client, error) {
	var client models.Client

	query := `
	INSERT INTO clients (telegram_id, username, first_name) 
	VALUES ($1, $2, $3)
	ON CONFLICT (telegram_id) 
	DO UPDATE SET username = EXCLUDED.username, first_name = EXCLUDED.first_name
	RETURNING telegram_id, username, first_name, last_name, strikes_count, is_banned, created_at;
`
	err := c.db.QueryRow(ctx, query, telegramID, username, firstName).Scan(
		&client.TelegramID,
		&client.Username,
		&client.FirstName,
		&client.LastName,
		&client.StrikesCount,
		&client.IsBanned,
		&client.CreatedAt,
	)

	return client, err
}

func (c *ClientRepo) GetClientByTelegramID(ctx context.Context, telegramID int64) (models.Client, error) {
	var client models.Client

	query := `SELECT telegram_id, username, first_name, last_name, strikes_count, is_banned, created_at 
FROM clients 
WHERE telegram_id = $1`

	err := c.db.QueryRow(ctx, query, telegramID).Scan(
		&client.TelegramID,
		&client.Username,
		&client.FirstName,
		&client.LastName,
		&client.StrikesCount,
		&client.IsBanned,
		&client.CreatedAt,
	)
	if err != nil {
		// Если строк не найдено, возвращаем пустую структуру и понятную ошибку
		if err == pgx.ErrNoRows {
			return models.Client{}, fmt.Errorf("client with telegram ID %d not found", telegramID)
		}
		// Если упало по другой причине (например, отвалилась база), возвращаем сырую ошибку
		return models.Client{}, err
	}

	return client, nil
}
