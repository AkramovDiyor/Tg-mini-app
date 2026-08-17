package repositories

import (
	"backend/internal/models"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceRepository interface {
	CreateService(ctx context.Context, service models.Service) error
	GetServicesByMasterID(ctx context.Context, masterID int64) ([]models.Service, error)
	GetServiceByID(ctx context.Context, serviceID int64) (models.Service, error) // НОВОЕ
}

type ServiceRepo struct {
	db *pgxpool.Pool
}

func NewServiceRepo(db *pgxpool.Pool) *ServiceRepo {
	return &ServiceRepo{db: db}
}

func (s *ServiceRepo) CreateService(ctx context.Context, service models.Service) error {
	_, err := s.db.Exec(ctx, "INSERT INTO services (master_id, name, duration_min, price) VALUES ($1, $2, $3, $4)", service.MasterID, service.Name, service.DurationMin, service.Price)
	return err
}

func (s *ServiceRepo) GetServicesByMasterID(ctx context.Context, masterID int64) ([]models.Service, error) {
	rows, err := s.db.Query(ctx, "SELECT id, master_id, name, duration_min, price FROM services WHERE master_id = $1", masterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []models.Service
	for rows.Next() {
		var srv models.Service
		err := rows.Scan(&srv.ID, &srv.MasterID, &srv.Name, &srv.DurationMin, &srv.Price)
		if err != nil {
			log.Println("Ошибка:", err)
			continue
		}
		services = append(services, srv)
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