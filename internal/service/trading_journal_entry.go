package service

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/user/normark/internal/dto"
	"github.com/user/normark/internal/entity"
	"github.com/user/normark/internal/storage"
	"github.com/user/normark/internal/types"
	"go.uber.org/zap"
)

type TradingJournalEntryStorage interface {
	Create(ctx context.Context, entry *entity.TradingJournalEntry) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.TradingJournalEntry, error)
	GetByIDWithJournal(ctx context.Context, id uuid.UUID) (*entity.TradingJournalEntry, error)
	GetByJournalID(ctx context.Context, params storage.GetByJournalIDParams) ([]*entity.TradingJournalEntry, error)
	GetByDateRange(ctx context.Context, params storage.GetByDateRangeParams) ([]*entity.TradingJournalEntry, error)
	GetByAsset(ctx context.Context, params storage.GetByAssetParams) ([]*entity.TradingJournalEntry, error)
	GetBySession(ctx context.Context, params storage.GetBySessionParams) ([]*entity.TradingJournalEntry, error)
	GetByResult(ctx context.Context, params storage.GetByResultParams) ([]*entity.TradingJournalEntry, error)
	Update(ctx context.Context, entry *entity.TradingJournalEntry) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entity.TradingJournalEntry, error)
	Count(ctx context.Context) (int, error)
	CountByJournalID(ctx context.Context, journalID uuid.UUID) (int, error)
	Exists(ctx context.Context, id uuid.UUID, journalID uuid.UUID) (bool, error)
	GetStatistics(ctx context.Context, journalID uuid.UUID) (map[string]any, error)
}

type TradingJournalEntryService struct {
	storage        TradingJournalEntryStorage
	journalStorage TradingJournalStorage
	logger         *zap.Logger
}

func NewTradingJournalEntryService(
	storage TradingJournalEntryStorage,
	journalStorage TradingJournalStorage,
	logger *zap.Logger,
) *TradingJournalEntryService {
	return &TradingJournalEntryService{
		storage:        storage,
		journalStorage: journalStorage,
		logger:         logger,
	}
}

func (s *TradingJournalEntryService) Create(ctx context.Context, journalID uuid.UUID, req *dto.CreateTradingJournalEntryRequest) (*entity.TradingJournalEntry, error) {
	_, err := s.journalStorage.GetByID(ctx, journalID)
	if err != nil {
		s.logger.Error("failed to verify journal existence", zap.Error(err), zap.String("journal_id", journalID.String()))
		return nil, errors.Wrap(err, "journal not found")
	}

	entry := entity.NewTradingJournalEntry(
		journalID,
		req.Day,
		req.Asset,
		req.LTF,
		req.HTF,
		req.EntryCharts,
		req.Session,
		req.TradeType,
		req.Setup,
		req.Direction,
		req.EntryType,
		req.Realized,
		req.MaxRR,
		req.Result,
		req.Notes,
	)

	if err := entry.Validate(); err != nil {
		s.logger.Error("invalid trading journal entry data", zap.Error(err))
		return nil, errors.Wrap(err, "invalid trading journal entry data")
	}

	if err := s.storage.Create(ctx, entry); err != nil {
		s.logger.Error("failed to create trading journal entry", zap.Error(err))
		return nil, errors.Wrap(err, "failed to create trading journal entry")
	}

	return entry, nil
}

func (s *TradingJournalEntryService) GetByID(ctx context.Context, id uuid.UUID) (*entity.TradingJournalEntry, error) {
	entry, err := s.storage.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get trading journal entry by id", zap.Error(err), zap.String("id", id.String()))
		return nil, errors.Wrap(err, "failed to get trading journal entry")
	}

	return entry, nil
}

func (s *TradingJournalEntryService) GetByIDWithJournal(ctx context.Context, id uuid.UUID) (*entity.TradingJournalEntry, error) {
	entry, err := s.storage.GetByIDWithJournal(ctx, id)
	if err != nil {
		s.logger.Error("failed to get trading journal entry by id with journal", zap.Error(err), zap.String("id", id.String()))
		return nil, errors.Wrap(err, "failed to get trading journal entry with journal")
	}

	return entry, nil
}

func (s *TradingJournalEntryService) GetJournalEntries(ctx context.Context, journalID uuid.UUID, limit, offset int) ([]*entity.TradingJournalEntry, error) {
	entries, err := s.storage.GetByJournalID(ctx, storage.GetByJournalIDParams{
		JournalID: journalID,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		s.logger.Error("failed to get journal entries", zap.Error(err), zap.String("journal_id", journalID.String()))
		return nil, errors.Wrap(err, "failed to get journal entries")
	}

	return entries, nil
}

func (s *TradingJournalEntryService) GetByDateRange(ctx context.Context, journalID uuid.UUID, startDate, endDate time.Time) ([]*entity.TradingJournalEntry, error) {
	entries, err := s.storage.GetByDateRange(ctx, storage.GetByDateRangeParams{
		JournalID: journalID,
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		s.logger.Error("failed to get entries by date range", zap.Error(err), zap.String("journal_id", journalID.String()))
		return nil, errors.Wrap(err, "failed to get entries by date range")
	}

	return entries, nil
}

func (s *TradingJournalEntryService) GetByAsset(ctx context.Context, journalID uuid.UUID, asset types.CurrencyPair, limit, offset int) ([]*entity.TradingJournalEntry, error) {
	entries, err := s.storage.GetByAsset(ctx, storage.GetByAssetParams{
		JournalID: journalID,
		Asset:     asset,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		s.logger.Error("failed to get entries by asset", zap.Error(err), zap.String("journal_id", journalID.String()), zap.String("asset", string(asset)))
		return nil, errors.Wrap(err, "failed to get entries by asset")
	}

	return entries, nil
}

func (s *TradingJournalEntryService) GetBySession(ctx context.Context, journalID uuid.UUID, session types.TradingSession, limit, offset int) ([]*entity.TradingJournalEntry, error) {
	entries, err := s.storage.GetBySession(ctx, storage.GetBySessionParams{
		JournalID: journalID,
		Session:   session,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		s.logger.Error("failed to get entries by session", zap.Error(err), zap.String("journal_id", journalID.String()), zap.String("session", string(session)))
		return nil, errors.Wrap(err, "failed to get entries by session")
	}

	return entries, nil
}

func (s *TradingJournalEntryService) GetByResult(ctx context.Context, journalID uuid.UUID, result types.TradeResult, limit, offset int) ([]*entity.TradingJournalEntry, error) {
	entries, err := s.storage.GetByResult(ctx, storage.GetByResultParams{
		JournalID: journalID,
		Result:    result,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		s.logger.Error("failed to get entries by result", zap.Error(err), zap.String("journal_id", journalID.String()), zap.String("result", string(result)))
		return nil, errors.Wrap(err, "failed to get entries by result")
	}

	return entries, nil
}

func (s *TradingJournalEntryService) Update(ctx context.Context, entry *entity.TradingJournalEntry) error {
	if err := entry.Validate(); err != nil {
		s.logger.Error("invalid trading journal entry data", zap.Error(err))
		return errors.Wrap(err, "invalid trading journal entry data")
	}

	if err := s.storage.Update(ctx, entry); err != nil {
		s.logger.Error("failed to update trading journal entry", zap.Error(err), zap.String("id", entry.ID.String()))
		return errors.Wrap(err, "failed to update trading journal entry")
	}

	return nil
}

func (s *TradingJournalEntryService) Delete(ctx context.Context, id uuid.UUID, journalID uuid.UUID) error {
	exists, err := s.storage.Exists(ctx, id, journalID)
	if err != nil {
		s.logger.Error("failed to check entry ownership", zap.Error(err))
		return errors.Wrap(err, "failed to verify entry ownership")
	}

	if !exists {
		return errors.New("trading journal entry not found or access denied")
	}

	if err := s.storage.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete trading journal entry", zap.Error(err), zap.String("id", id.String()))
		return errors.Wrap(err, "failed to delete trading journal entry")
	}

	return nil
}

func (s *TradingJournalEntryService) CountJournalEntries(ctx context.Context, journalID uuid.UUID) (int, error) {
	count, err := s.storage.CountByJournalID(ctx, journalID)
	if err != nil {
		s.logger.Error("failed to count journal entries", zap.Error(err), zap.String("journal_id", journalID.String()))
		return 0, errors.Wrap(err, "failed to count journal entries")
	}

	return count, nil
}

func (s *TradingJournalEntryService) GetStatistics(ctx context.Context, journalID uuid.UUID) (map[string]any, error) {
	stats, err := s.storage.GetStatistics(ctx, journalID)
	if err != nil {
		s.logger.Error("failed to get journal statistics", zap.Error(err), zap.String("journal_id", journalID.String()))
		return nil, errors.Wrap(err, "failed to get journal statistics")
	}

	if totalTrades, ok := stats["total_trades"].(int); ok && totalTrades > 0 {
		wins := 0
		if w, ok := stats["wins"].(int); ok {
			wins = w
		}
		winRate := float64(wins) / float64(totalTrades) * 100
		stats["win_rate"] = winRate
	} else {
		stats["win_rate"] = 0.0
	}

	return stats, nil
}

func (s *TradingJournalEntryService) VerifyAccess(ctx context.Context, entryID uuid.UUID, journalID uuid.UUID) (bool, error) {
	exists, err := s.storage.Exists(ctx, entryID, journalID)
	if err != nil {
		s.logger.Error("failed to verify entry access", zap.Error(err))
		return false, errors.Wrap(err, "failed to verify entry access")
	}

	return exists, nil
}
