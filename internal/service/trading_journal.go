package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/user/normark/internal/dto"
	"github.com/user/normark/internal/entity"
	"go.uber.org/zap"
)

type TradingJournalStorage interface {
	Create(ctx context.Context, journal *entity.TradingJournal) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.TradingJournal, error)
	GetByIDWithEntries(ctx context.Context, id uuid.UUID) (*entity.TradingJournal, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.TradingJournal, error)
	Update(ctx context.Context, journal *entity.TradingJournal) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entity.TradingJournal, error)
	Count(ctx context.Context) (int, error)
	CountByUserID(ctx context.Context, userID uuid.UUID) (int, error)
	Exists(ctx context.Context, id uuid.UUID, userID uuid.UUID) (bool, error)
}

type TradingJournalService struct {
	storage TradingJournalStorage
	cache   Cache
	logger  *zap.Logger
}

func NewTradingJournalService(
	storage TradingJournalStorage,
	logger *zap.Logger,
) *TradingJournalService {
	return &TradingJournalService{
		storage: storage,
		logger:  logger,
	}
}

func (s *TradingJournalService) WithCache(cache Cache) *TradingJournalService {
	s.cache = cache
	return s
}

func (s *TradingJournalService) Create(ctx context.Context, userID uuid.UUID, req *dto.CreateTradingJournalRequest) (*entity.TradingJournal, error) {
	journal := entity.NewTradingJournal(userID, req.Name, req.Description)

	if err := journal.Validate(); err != nil {
		s.logger.Error("invalid trading journal data", zap.Error(err))
		return nil, errors.Wrap(err, "invalid trading journal data")
	}

	if err := s.storage.Create(ctx, journal); err != nil {
		s.logger.Error("failed to create trading journal", zap.Error(err))
		return nil, errors.Wrap(err, "failed to create trading journal")
	}

	return journal, nil
}

func (s *TradingJournalService) GetByID(ctx context.Context, id uuid.UUID) (*entity.TradingJournal, error) {
	cacheKey := fmt.Sprintf("journal:%s", id.String())

	// Try to get from cache first
	if s.cache != nil {
		cached, err := s.cache.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			var journal entity.TradingJournal
			if err := json.Unmarshal([]byte(cached), &journal); err == nil {
				return &journal, nil
			}
		}
	}

	// Cache miss or error, fetch from database
	journal, err := s.storage.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get trading journal by id", zap.Error(err), zap.String("id", id.String()))
		return nil, errors.Wrap(err, "failed to get trading journal")
	}

	// Cache the result
	if s.cache != nil {
		if data, err := json.Marshal(journal); err == nil {
			if err := s.cache.Set(ctx, cacheKey, string(data), 15*time.Minute); err != nil {
				s.logger.Warn("failed to cache trading journal", zap.Error(err))
			}
		}
	}

	return journal, nil
}

func (s *TradingJournalService) GetByIDWithEntries(ctx context.Context, id uuid.UUID) (*entity.TradingJournal, error) {
	journal, err := s.storage.GetByIDWithEntries(ctx, id)
	if err != nil {
		s.logger.Error("failed to get trading journal by id with entries", zap.Error(err), zap.String("id", id.String()))
		return nil, errors.Wrap(err, "failed to get trading journal with entries")
	}

	return journal, nil
}

func (s *TradingJournalService) GetUserJournals(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.TradingJournal, error) {
	journals, err := s.storage.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		s.logger.Error("failed to get user journals", zap.Error(err), zap.String("user_id", userID.String()))
		return nil, errors.Wrap(err, "failed to get user journals")
	}

	return journals, nil
}

func (s *TradingJournalService) Update(ctx context.Context, journal *entity.TradingJournal) error {
	if err := journal.Validate(); err != nil {
		s.logger.Error("invalid trading journal data", zap.Error(err))
		return errors.Wrap(err, "invalid trading journal data")
	}

	if err := s.storage.Update(ctx, journal); err != nil {
		s.logger.Error("failed to update trading journal", zap.Error(err), zap.String("id", journal.ID.String()))
		return errors.Wrap(err, "failed to update trading journal")
	}

	if s.cache != nil {
		cacheKey := fmt.Sprintf("journal:%s", journal.ID.String())
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			s.logger.Warn("failed to invalidate cache after update", zap.Error(err))
		}
	}

	return nil
}

func (s *TradingJournalService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	exists, err := s.storage.Exists(ctx, id, userID)
	if err != nil {
		s.logger.Error("failed to check journal ownership", zap.Error(err))
		return errors.Wrap(err, "failed to verify journal ownership")
	}

	if !exists {
		return errors.New("trading journal not found or access denied")
	}

	if err := s.storage.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete trading journal", zap.Error(err), zap.String("id", id.String()))
		return errors.Wrap(err, "failed to delete trading journal")
	}

	if s.cache != nil {
		cacheKey := fmt.Sprintf("journal:%s", id.String())
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			s.logger.Warn("failed to invalidate cache after delete", zap.Error(err))
		}
	}

	return nil
}

func (s *TradingJournalService) CountUserJournals(ctx context.Context, userID uuid.UUID) (int, error) {
	count, err := s.storage.CountByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("failed to count user journals", zap.Error(err), zap.String("user_id", userID.String()))
		return 0, errors.Wrap(err, "failed to count user journals")
	}

	return count, nil
}

func (s *TradingJournalService) VerifyAccess(ctx context.Context, journalID uuid.UUID, userID uuid.UUID) (bool, error) {
	exists, err := s.storage.Exists(ctx, journalID, userID)
	if err != nil {
		s.logger.Error("failed to verify journal access", zap.Error(err))
		return false, errors.Wrap(err, "failed to verify journal access")
	}

	return exists, nil
}
