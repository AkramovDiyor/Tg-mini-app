package repositories

import (
	"backend/internal/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceRepository interface {
	CreateService(ctx context.Context, service models.Service) error
	GetServicesByMasterID(ctx context.Context, masterID int64) ([]models.Service, error)
	GetServiceByID(ctx context.Context, serviceID int64) (models.Service, error) 
	UpdateService(ctx context.Context, serviceID int64, masterID int64, name string, duration int, price int) error
	DeleteService(ctx context.Context, serviceID int64, masterID int64) error
}

type ServiceRepo struct {
	db *pgxpool.Pool
}

func NewServiceRepo(db *pgxpool.Pool) *ServiceRepo {
	return &ServiceRepo{db: db}
}

func (s *ServiceRepo) CreateService(ctx context.Context, service models.Service) error {
    query := `
        INSERT INTO services (master_id, name, duration_min, price, is_deleted) 
        VALUES ($1, $2, $3, $4, FALSE)
    `
    _, err := s.db.Exec(ctx, query, service.MasterID, service.Name, service.DurationMin, service.Price)
    return err
}

func (s *ServiceRepo) GetServicesByMasterID(ctx context.Context, masterID int64) ([]models.Service, error) {
    query := `
        SELECT id, master_id, name, duration_min, price, is_deleted, created_at, updated_at 
        FROM services 
        WHERE master_id = $1 AND is_deleted = FALSE
        ORDER BY created_at ASC
    `
    rows, err := s.db.Query(ctx, query, masterID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var services []models.Service
    for rows.Next() {
        var svc models.Service
        err := rows.Scan(
            &svc.ID, &svc.MasterID, &svc.Name, 
            &svc.DurationMin, &svc.Price, &svc.IsDeleted,
            &svc.CreatedAt, &svc.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        services = append(services, svc)
    }
    return services, nil
}

// НОВЫЙ МЕТОД
func (s *ServiceRepo) GetServiceByID(ctx context.Context, serviceID int64) (models.Service, error) {
	var service models.Service
	query := `SELECT id, master_id, name, duration_min, price FROM services WHERE id = $1`
	err := s.db.QueryRow(ctx, query, serviceID).Scan(
		&service.ID, &service.MasterID, &service.Name, &service.DurationMin, &service.Price,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.Service{}, fmt.Errorf("service with ID %d not found", serviceID)
		}
		return models.Service{}, err
	}
	return service, nil
}


func (s *ServiceRepo) UpdateService(ctx context.Context, serviceID int64, masterID int64, name string, duration int, price int) error {
    query := `UPDATE services SET name = $1, duration_min = $2, price = $3, updated_at = NOW() WHERE id = $4 AND master_id = $5`
    
    result, err := s.db.Exec(ctx, query, name, duration, price, serviceID, masterID)
    if err != nil {
        return fmt.Errorf("failed to update service: %w", err)
    }
    
    // Проверяем, что запись была обновлена
    rowsAffected := result.RowsAffected()
    if rowsAffected == 0 {
        return fmt.Errorf("service not found or not owned by master")
    }
    
    return nil
}

// 🔥 МЯГКОЕ УДАЛЕНИЕ: помечаем как удаленную вместо физического DELETE
func (s *ServiceRepo) DeleteService(ctx context.Context, serviceID int64, masterID int64) error {
    query := `
        UPDATE services 
        SET is_deleted = TRUE, updated_at = NOW() 
        WHERE id = $1 AND master_id = $2 AND is_deleted = FALSE
    `
    
    result, err := s.db.Exec(ctx, query, serviceID, masterID)
    if err != nil {
        return fmt.Errorf("failed to archive service: %w", err)
    }
    
    rowsAffected := result.RowsAffected()
    if rowsAffected == 0 {
        return fmt.Errorf("service not found, not owned by master, or already deleted")
    }
    
    return nil
}