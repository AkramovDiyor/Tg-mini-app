package services

import (
	"backend/internal/repositories"
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WaitlistService struct {
	pool        *pgxpool.Pool
	slotRepo    repositories.SlotRepository
	waitlistRepo repositories.WaitlistRepository
}

func NewWaitlistService(pool *pgxpool.Pool, slotRepo repositories.SlotRepository, waitlistRepo repositories.WaitlistRepository) *WaitlistService {
    return &WaitlistService{
        pool:         pool,
        slotRepo:     slotRepo,
        waitlistRepo: waitlistRepo,
    }
}



func (s *WaitlistService) JoinWaitlist(ctx context.Context, masterID, clientTgID int64, desiredDate time.Time) error{

	tx, err := s.pool.Begin(ctx)
	defer tx.Rollback(ctx)

	err = s.waitlistRepo.AddToWaitlist(ctx, tx, masterID, clientTgID, desiredDate)
	if err != nil{
		return err
	}
	err = tx.Commit(ctx)
	return err
}

func (s *WaitlistService) OfferSlotToFirstInLine(ctx context.Context, masterID, slotID int64, desiredDate time.Time) error{

	req, err := s.waitlistRepo.GetFirstInWaitlist(ctx, masterID, desiredDate)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Println("🤖 Авто-забиватор: Очередь пуста, совпадений по дате нет.")
			return nil 
		}
		// ЭТОТ ЛОГ ПОКАЖЕТ НАМ НАСТОЯЩУЮ ПРИЧИНУ!
		log.Println("🚨 Критическая ошибка внутри GetFirstInWaitlist:", err)
		return err 
	}


	tx, err := s.pool.Begin(ctx)
	defer tx.Rollback(ctx)

	err = s.waitlistRepo.UpdateWaitlistStatus(ctx, tx, req.ID, "offered", slotID)
	if err != nil {
		return nil
	}
	err = tx.Commit(ctx)
	// Здесь в будущем я допишем отправку сообщения ботом
	return err
}