package repositories

import (
	"backend/internal/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MasterRepository interface {
	CreateMaster(ctx context.Context, master models.Master) error
	GetMasterByInviteLink(ctx context.Context, link string) (models.Master, error)
}

type MasterRepo struct {
	db *pgxpool.Pool
}

func NewMasterRepo(db *pgxpool.Pool) *MasterRepo {
	return &MasterRepo{
		db: db,
	}
}

func (m *MasterRepo) CreateMaster(ctx context.Context, master models.Master) error {
	_, err := m.db.Exec(ctx, "INSERT INTO masters (telegram_id, name, bio, address, invite_link) VALUES ($1, $2, $3, $4, $5)", master.TelegramID, master.Name, master.Bio, master.Address, master.InviteLink)
	return err
}

func (m *MasterRepo) GetMasterByInviteLink(ctx context.Context, link string) (models.Master, error) {
	var master models.Master
	err := m.db.QueryRow(ctx, "SELECT id, telegram_id, name, bio, address, invite_link FROM masters WHERE invite_link = $1", link).Scan(&master.ID, &master.TelegramID, &master.Name, &master.Bio, &master.Address, &master.InviteLink)
	return master, err
}
