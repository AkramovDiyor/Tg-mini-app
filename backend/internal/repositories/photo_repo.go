package repositories

import (
	"backend/internal/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PhotoRepository interface {
	AddPhoto(ctx context.Context, masterID int64, url string) (int64, error)
	GetPhotosByMaster(ctx context.Context, masterID int64) ([]models.MasterPhoto, error)
	DeletePhoto(ctx context.Context, photoID int64, masterID int64) error
}

type PhotoRepo struct {
	db *pgxpool.Pool
}

func NewPhotoRepo(db *pgxpool.Pool) *PhotoRepo {
	return &PhotoRepo{db: db}
}

// AddPhoto добавляет фото и возвращает его ID
func (p *PhotoRepo) AddPhoto(ctx context.Context, masterID int64, url string) (int64, error) {
	var photoID int64
	query := `INSERT INTO master_photos (master_id, url) VALUES ($1, $2) RETURNING id`
	
	err := p.db.QueryRow(ctx, query, masterID, url).Scan(&photoID)
	if err != nil {
		return 0, fmt.Errorf("failed to add photo: %w", err)
	}
	
	return photoID, nil
}

// GetPhotosByMaster возвращает все фото мастера
func (p *PhotoRepo) GetPhotosByMaster(ctx context.Context, masterID int64) ([]models.MasterPhoto, error) {
	query := `
		SELECT id, master_id, url, created_at 
		FROM master_photos 
		WHERE master_id = $1 
		ORDER BY created_at DESC
	`
	
	rows, err := p.db.Query(ctx, query, masterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get photos: %w", err)
	}
	defer rows.Close()

	var photos []models.MasterPhoto
	for rows.Next() {
		var photo models.MasterPhoto
		err := rows.Scan(&photo.ID, &photo.MasterID, &photo.URL, &photo.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan photo: %w", err)
		}
		photos = append(photos, photo)
	}

	return photos, nil
}

// DeletePhoto удаляет фото с проверкой принадлежности мастеру
func (p *PhotoRepo) DeletePhoto(ctx context.Context, photoID int64, masterID int64) error {
	query := `DELETE FROM master_photos WHERE id = $1 AND master_id = $2`
	
	result, err := p.db.Exec(ctx, query, photoID, masterID)
	if err != nil {
		return fmt.Errorf("failed to delete photo: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("photo not found or not owned by master")
	}

	return nil
}