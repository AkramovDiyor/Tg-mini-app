package repositories

import (
	"backend/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ========== INTERFACE ==========
type MasterRepository interface {
	CreateMaster(ctx context.Context, master models.Master) error
	GetMasterByInviteLink(ctx context.Context, link string) (models.Master, error)
	GetMasterByTelegramID(ctx context.Context, telegramID int64) (models.Master, error)
	UpdateMasterProfile(ctx context.Context, masterID int64, name, bio, address string) error
	UpdateMasterSettings(ctx context.Context, masterID int64, workHours models.WorkHours, settings models.MasterSettings) error
}

// ========== REPOSITORY ==========
type MasterRepo struct {
	db *pgxpool.Pool
}

func NewMasterRepo(db *pgxpool.Pool) *MasterRepo {
	return &MasterRepo{db: db}
}

// ========== CREATE ==========
func (m *MasterRepo) CreateMaster(ctx context.Context, master models.Master) error {
	_, err := m.db.Exec(ctx,
		"INSERT INTO masters (telegram_id, name, bio, address, invite_link) VALUES ($1, $2, $3, $4, $5)",
		master.TelegramID, master.Name, master.Bio, master.Address, master.InviteLink,
	)
	return err
}

// 🔥 ВСПОМОГАТЕЛЬНАЯ ФУНКЦИЯ: парсит JSONB поля мастера
func (m *MasterRepo) parseMasterJSON(master *models.Master, workHoursJSON, settingsJSON []byte) {
	if len(workHoursJSON) > 0 {
		err := json.Unmarshal(workHoursJSON, &master.WorkHours)
		if err != nil {
			log.Printf("❌ Ошибка парсинга work_hours: %v (raw: %s)", err, string(workHoursJSON))
		} else {
			log.Printf("✅ Распарсен WorkHours: WorkDays=%v, StartTime=%s, EndTime=%s",
				master.WorkHours.WorkDays, master.WorkHours.StartTime, master.WorkHours.EndTime)
		}
	}
	if len(settingsJSON) > 0 {
		err := json.Unmarshal(settingsJSON, &master.Settings)
		if err != nil {
			log.Printf("❌ Ошибка парсинга settings: %v", err)
		}
	}
}

// ========== GET BY INVITE LINK ==========
func (m *MasterRepo) GetMasterByInviteLink(ctx context.Context, link string) (models.Master, error) {
	var master models.Master
	var workHoursJSON, settingsJSON []byte

	err := m.db.QueryRow(ctx,
		"SELECT id, telegram_id, name, bio, address, invite_link, work_hours, settings FROM masters WHERE invite_link = $1",
		link,
	).Scan(&master.ID, &master.TelegramID, &master.Name, &master.Bio, &master.Address, &master.InviteLink, &workHoursJSON, &settingsJSON)

	if err != nil {
		if err == pgx.ErrNoRows {
			return models.Master{}, fmt.Errorf("master not found")
		}
		return master, err
	}

	log.Printf("📊 Сырой JSON work_hours: %s", string(workHoursJSON))
	log.Printf("📊 Сырой JSON settings: %s", string(settingsJSON))

	// Парсим JSON в структуры
	m.parseMasterJSON(&master, workHoursJSON, settingsJSON)

	return master, nil
}

// ========== GET BY TELEGRAM ID ==========
func (m *MasterRepo) GetMasterByTelegramID(ctx context.Context, telegramID int64) (models.Master, error) {
	var master models.Master
	var workHoursJSON, settingsJSON []byte  // 🔥 ВАЖНО: сканируем в []byte

	err := m.db.QueryRow(ctx,
		"SELECT id, telegram_id, name, bio, address, invite_link, work_hours, settings FROM masters WHERE telegram_id = $1",
		telegramID,
	).Scan(
		&master.ID,
		&master.TelegramID,
		&master.Name,
		&master.Bio,
		&master.Address,
		&master.InviteLink,
		&workHoursJSON,   // 🔥 в []byte, а не в структуру!
		&settingsJSON,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return models.Master{}, fmt.Errorf("master not found")
		}
		return models.Master{}, err
	}

	// 🔥 Парсим JSON в структуры (используем ту же вспомогательную функцию)
	m.parseMasterJSON(&master, workHoursJSON, settingsJSON)

	return master, nil
}

// ========== UPDATE PROFILE ==========
func (m *MasterRepo) UpdateMasterProfile(ctx context.Context, masterID int64, name, bio, address string) error {
	query := `
		UPDATE masters 
		SET name = $1, bio = $2, address = $3, updated_at = NOW()
		WHERE id = $4
	`

	result, err := m.db.Exec(ctx, query, name, bio, address, masterID)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("master with id %d not found", masterID)
	}

	return nil
}

// ========== UPDATE SETTINGS ==========
func (m *MasterRepo) UpdateMasterSettings(ctx context.Context, masterID int64, workHours models.WorkHours, settings models.MasterSettings) error {
	workHoursJSON, err := json.Marshal(workHours)
	if err != nil {
		return fmt.Errorf("failed to marshal work_hours: %w", err)
	}

	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	query := `
		UPDATE masters 
		SET work_hours = $1::jsonb, settings = $2::jsonb, updated_at = NOW()
		WHERE id = $3
	`
	result, err := m.db.Exec(ctx, query, workHoursJSON, settingsJSON, masterID)
	if err != nil {
		return fmt.Errorf("failed to update settings: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("master with id %d not found", masterID)
	}

	log.Printf("✅ Настройки обновлены для мастера %d", masterID)
	return nil
}