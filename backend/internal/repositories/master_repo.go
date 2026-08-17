package repositories

import (
	"backend/internal/models"
	"context"
	"encoding/json"

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
	return &MasterRepo{db: db}
}

func (m *MasterRepo) CreateMaster(ctx context.Context, master models.Master) error {
	// Сериализуем WorkHours и Settings в JSON для сохранения в базу
	workHoursJSON, _ := json.Marshal(master.WorkHours)
	settingsJSON, _ := json.Marshal(master.Settings)
	
	_, err := m.db.Exec(ctx, 
		"INSERT INTO masters (telegram_id, name, bio, address, invite_link, work_hours, settings) VALUES ($1, $2, $3, $4, $5, $6, $7)", 
		master.TelegramID, master.Name, master.Bio, master.Address, master.InviteLink, workHoursJSON, settingsJSON,
	)
	return err
}

// ИСПРАВЛЕНО: добавляем work_hours и settings
func (m *MasterRepo) GetMasterByInviteLink(ctx context.Context, link string) (models.Master, error) {
	var master models.Master
	var workHoursJSON, settingsJSON []byte
	
	err := m.db.QueryRow(ctx, 
		"SELECT id, telegram_id, name, bio, address, invite_link, work_hours, settings FROM masters WHERE invite_link = $1", 
		link,
	).Scan(&master.ID, &master.TelegramID, &master.Name, &master.Bio, &master.Address, &master.InviteLink, &workHoursJSON, &settingsJSON)
	
	if err != nil {
		return master, err
	}
	
	// Десериализуем JSON обратно в структуры
	if len(workHoursJSON) > 0 {
		json.Unmarshal(workHoursJSON, &master.WorkHours)
	}
	if len(settingsJSON) > 0 {
		json.Unmarshal(settingsJSON, &master.Settings)
	}
	
	return master, nil
}