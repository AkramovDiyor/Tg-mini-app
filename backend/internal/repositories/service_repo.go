package repositories

import (
	"backend/internal/models"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceRepository interface {
	CreateService(ctx context.Context, service models.Service) error
	GetServicesByMasterID(ctx context.Context, masterID int64) ([]models.Service, error)
}

type ServiceRepo struct {
	db *pgxpool.Pool
}

func NewServiceRepo(db *pgxpool.Pool) *ServiceRepo {
	return &ServiceRepo{
		db: db,
	}
}

func (s *ServiceRepo) CreateService(ctx context.Context, service models.Service) error {
	_, err := s.db.Exec(ctx, "INSERT INTO services (master_id, name, duration_min, price) VALUES ($1, $2, $3, $4)", service.MasterID, service.Name, service.DurationMin, service.Price)
	return err
}

func (s *ServiceRepo) GetServicesByMasterID(ctx context.Context, masterID int64) ([]models.Service, error){
	rows, err := s.db.Query(ctx, "SELECT id, master_id, name, duration_min, price FROM services WHERE master_id = $1", masterID)

	var services []models.Service

	for rows.Next() {
		var s models.Service
		err := rows.Scan(&s.ID, &s.MasterID, &s.Name, &s.DurationMin, &s.Price)
		if err != nil {
			log.Println("Ошибка:", err)
		}

		services = append(services, s)
	}
	defer rows.Close()
	return services, err
}
