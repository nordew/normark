package bun

import (
	"context"
	"database/sql"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/user/normark/internal/entity"
	"github.com/user/normark/internal/types"
)

type TradingJournalEntryStorage struct {
	db *bun.DB
}

func NewTradingJournalEntryStorage(db *bun.DB) *TradingJournalEntryStorage {
	return &TradingJournalEntryStorage{
		db: db,
	}
}

type GetByJournalIDParams struct {
	JournalID uuid.UUID
	Limit     int
	Offset    int
}

type GetByDateRangeParams struct {
	JournalID uuid.UUID
	StartDate time.Time
	EndDate   time.Time
}

type GetByAssetParams struct {
	JournalID uuid.UUID
	Asset     types.CurrencyPair
	Limit     int
	Offset    int
}

type GetBySessionParams struct {
	JournalID uuid.UUID
	Session   types.TradingSession
	Limit     int
	Offset    int
}

type GetByResultParams struct {
	JournalID uuid.UUID
	Result    types.TradeResult
	Limit     int
	Offset    int
}

func (s *TradingJournalEntryStorage) Create(ctx context.Context, entry *entity.TradingJournalEntry) error {
	_, err := s.db.NewInsert().
		Model(entry).
		Exec(ctx)

	if err != nil {
		return errors.Wrap(err, "failed to create trading journal entry")
	}

	return nil
}

func (s *TradingJournalEntryStorage) GetByID(ctx context.Context, id uuid.UUID) (*entity.TradingJournalEntry, error) {
	entry := new(entity.TradingJournalEntry)

	err := s.db.NewSelect().
		Model(entry).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(err, "trading journal entry not found")
		}
		return nil, errors.Wrap(err, "failed to get trading journal entry by id")
	}

	return entry, nil
}

func (s *TradingJournalEntryStorage) GetByIDWithJournal(ctx context.Context, id uuid.UUID) (*entity.TradingJournalEntry, error) {
	entry := new(entity.TradingJournalEntry)

	err := s.db.NewSelect().
		Model(entry).
		Relation("Journal").
		Where("tje.id = ?", id).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(err, "trading journal entry not found")
		}
		return nil, errors.Wrap(err, "failed to get trading journal entry by id with journal")
	}

	return entry, nil
}

func (s *TradingJournalEntryStorage) GetByJournalID(ctx context.Context, params GetByJournalIDParams) ([]*entity.TradingJournalEntry, error) {
	var entries []*entity.TradingJournalEntry

	err := s.db.NewSelect().
		Model(&entries).
		Where("journal_id = ?", params.JournalID).
		Limit(params.Limit).
		Offset(params.Offset).
		Order("day DESC").
		Scan(ctx)

	if err != nil {
		return nil, errors.Wrap(err, "failed to get trading journal entries by journal id")
	}

	return entries, nil
}

func (s *TradingJournalEntryStorage) GetByDateRange(ctx context.Context, params GetByDateRangeParams) ([]*entity.TradingJournalEntry, error) {
	var entries []*entity.TradingJournalEntry

	err := s.db.NewSelect().
		Model(&entries).
		Where("journal_id = ?", params.JournalID).
		Where("day >= ?", params.StartDate).
		Where("day <= ?", params.EndDate).
		Order("day DESC").
		Scan(ctx)

	if err != nil {
		return nil, errors.Wrap(err, "failed to get trading journal entries by date range")
	}

	return entries, nil
}

func (s *TradingJournalEntryStorage) GetByAsset(ctx context.Context, params GetByAssetParams) ([]*entity.TradingJournalEntry, error) {
	var entries []*entity.TradingJournalEntry

	err := s.db.NewSelect().
		Model(&entries).
		Where("journal_id = ?", params.JournalID).
		Where("asset = ?", params.Asset).
		Limit(params.Limit).
		Offset(params.Offset).
		Order("day DESC").
		Scan(ctx)

	if err != nil {
		return nil, errors.Wrap(err, "failed to get trading journal entries by asset")
	}

	return entries, nil
}

func (s *TradingJournalEntryStorage) GetBySession(ctx context.Context, params GetBySessionParams) ([]*entity.TradingJournalEntry, error) {
	var entries []*entity.TradingJournalEntry

	err := s.db.NewSelect().
		Model(&entries).
		Where("journal_id = ?", params.JournalID).
		Where("session = ?", params.Session).
		Limit(params.Limit).
		Offset(params.Offset).
		Order("day DESC").
		Scan(ctx)

	if err != nil {
		return nil, errors.Wrap(err, "failed to get trading journal entries by session")
	}

	return entries, nil
}

func (s *TradingJournalEntryStorage) GetByResult(ctx context.Context, params GetByResultParams) ([]*entity.TradingJournalEntry, error) {
	var entries []*entity.TradingJournalEntry

	err := s.db.NewSelect().
		Model(&entries).
		Where("journal_id = ?", params.JournalID).
		Where("result = ?", params.Result).
		Limit(params.Limit).
		Offset(params.Offset).
		Order("day DESC").
		Scan(ctx)

	if err != nil {
		return nil, errors.Wrap(err, "failed to get trading journal entries by result")
	}

	return entries, nil
}

func (s *TradingJournalEntryStorage) Update(ctx context.Context, entry *entity.TradingJournalEntry) error {
	result, err := s.db.NewUpdate().
		Model(entry).
		WherePK().
		Exec(ctx)

	if err != nil {
		return errors.Wrap(err, "failed to update trading journal entry")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.New("trading journal entry not found")
	}

	return nil
}

func (s *TradingJournalEntryStorage) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := s.db.NewDelete().
		Model((*entity.TradingJournalEntry)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	if err != nil {
		return errors.Wrap(err, "failed to delete trading journal entry")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.New("trading journal entry not found")
	}

	return nil
}

func (s *TradingJournalEntryStorage) List(ctx context.Context, limit, offset int) ([]*entity.TradingJournalEntry, error) {
	var entries []*entity.TradingJournalEntry

	err := s.db.NewSelect().
		Model(&entries).
		Limit(limit).
		Offset(offset).
		Order("day DESC").
		Scan(ctx)

	if err != nil {
		return nil, errors.Wrap(err, "failed to list trading journal entries")
	}

	return entries, nil
}

func (s *TradingJournalEntryStorage) Count(ctx context.Context) (int, error) {
	count, err := s.db.NewSelect().
		Model((*entity.TradingJournalEntry)(nil)).
		Count(ctx)

	if err != nil {
		return 0, errors.Wrap(err, "failed to count trading journal entries")
	}

	return count, nil
}

func (s *TradingJournalEntryStorage) CountByJournalID(ctx context.Context, journalID uuid.UUID) (int, error) {
	count, err := s.db.NewSelect().
		Model((*entity.TradingJournalEntry)(nil)).
		Where("journal_id = ?", journalID).
		Count(ctx)

	if err != nil {
		return 0, errors.Wrap(err, "failed to count trading journal entries by journal id")
	}

	return count, nil
}

func (s *TradingJournalEntryStorage) Exists(ctx context.Context, id uuid.UUID, journalID uuid.UUID) (bool, error) {
	count, err := s.db.NewSelect().
		Model((*entity.TradingJournalEntry)(nil)).
		Where("id = ? AND journal_id = ?", id, journalID).
		Count(ctx)

	if err != nil {
		return false, errors.Wrap(err, "failed to check if trading journal entry exists")
	}

	return count > 0, nil
}

func (s *TradingJournalEntryStorage) GetStatistics(ctx context.Context, journalID uuid.UUID) (map[string]any, error) {
	stats := make(map[string]any)

	totalTrades, err := s.CountByJournalID(ctx, journalID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to count total trades")
	}
	stats["total_trades"] = totalTrades

	var resultStats []struct {
		Result types.TradeResult
		Count  int
	}
	err = s.db.NewSelect().
		Model((*entity.TradingJournalEntry)(nil)).
		Column("result").
		ColumnExpr("COUNT(*) as count").
		Where("journal_id = ?", journalID).
		Group("result").
		Scan(ctx, &resultStats)

	if err != nil {
		return nil, errors.Wrap(err, "failed to get result statistics")
	}

	for _, stat := range resultStats {
		switch stat.Result {
		case types.TradeResultTakeProfit:
			stats["wins"] = stat.Count
		case types.TradeResultStopLoss:
			stats["losses"] = stat.Count
		case types.TradeResultBreakEven:
			stats["break_even"] = stat.Count
		}
	}

	var totalRealized float64
	err = s.db.NewSelect().
		Model((*entity.TradingJournalEntry)(nil)).
		ColumnExpr("COALESCE(SUM(realized), 0) as total").
		Where("journal_id = ?", journalID).
		Scan(ctx, &totalRealized)

	if err != nil {
		return nil, errors.Wrap(err, "failed to calculate total realized")
	}
	stats["total_realized"] = totalRealized

	var avgRR float64
	err = s.db.NewSelect().
		Model((*entity.TradingJournalEntry)(nil)).
		ColumnExpr("COALESCE(AVG(max_rr), 0) as avg").
		Where("journal_id = ?", journalID).
		Scan(ctx, &avgRR)

	if err != nil {
		return nil, errors.Wrap(err, "failed to calculate average RR")
	}
	stats["avg_risk_reward"] = avgRR

	return stats, nil
}
