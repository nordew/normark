package bun

import (
	"context"
	"database/sql"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/user/normark/internal/entity"
)

type TradingJournalStorage struct {
	db *bun.DB
}

func NewTradingJournalStorage(db *bun.DB) *TradingJournalStorage {
	return &TradingJournalStorage{
		db: db,
	}
}

func (s *TradingJournalStorage) Create(ctx context.Context, journal *entity.TradingJournal) error {
	_, err := s.db.NewInsert().
		Model(journal).
		Exec(ctx)

	if err != nil {
		return errors.Wrap(err, "failed to create trading journal")
	}

	return nil
}

func (s *TradingJournalStorage) GetByID(ctx context.Context, id uuid.UUID) (*entity.TradingJournal, error) {
	journal := new(entity.TradingJournal)

	err := s.db.NewSelect().
		Model(journal).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(err, "trading journal not found")
		}
		return nil, errors.Wrap(err, "failed to get trading journal by id")
	}

	return journal, nil
}

func (s *TradingJournalStorage) GetByIDWithEntries(ctx context.Context, id uuid.UUID) (*entity.TradingJournal, error) {
	journal := new(entity.TradingJournal)

	err := s.db.NewSelect().
		Model(journal).
		Relation("Entries").
		Where("tj.id = ?", id).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(err, "trading journal not found")
		}
		return nil, errors.Wrap(err, "failed to get trading journal by id with entries")
	}

	return journal, nil
}

func (s *TradingJournalStorage) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.TradingJournal, error) {
	var journals []*entity.TradingJournal

	err := s.db.NewSelect().
		Model(&journals).
		Where("user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Scan(ctx)

	if err != nil {
		return nil, errors.Wrap(err, "failed to get trading journals by user id")
	}

	return journals, nil
}

func (s *TradingJournalStorage) Update(ctx context.Context, journal *entity.TradingJournal) error {
	result, err := s.db.NewUpdate().
		Model(journal).
		WherePK().
		Exec(ctx)

	if err != nil {
		return errors.Wrap(err, "failed to update trading journal")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.New("trading journal not found")
	}

	return nil
}

func (s *TradingJournalStorage) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := s.db.NewDelete().
		Model((*entity.TradingJournal)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	if err != nil {
		return errors.Wrap(err, "failed to delete trading journal")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.New("trading journal not found")
	}

	return nil
}

func (s *TradingJournalStorage) List(ctx context.Context, limit, offset int) ([]*entity.TradingJournal, error) {
	var journals []*entity.TradingJournal

	err := s.db.NewSelect().
		Model(&journals).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Scan(ctx)

	if err != nil {
		return nil, errors.Wrap(err, "failed to list trading journals")
	}

	return journals, nil
}

func (s *TradingJournalStorage) Count(ctx context.Context) (int, error) {
	count, err := s.db.NewSelect().
		Model((*entity.TradingJournal)(nil)).
		Count(ctx)

	if err != nil {
		return 0, errors.Wrap(err, "failed to count trading journals")
	}

	return count, nil
}

func (s *TradingJournalStorage) CountByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	count, err := s.db.NewSelect().
		Model((*entity.TradingJournal)(nil)).
		Where("user_id = ?", userID).
		Count(ctx)

	if err != nil {
		return 0, errors.Wrap(err, "failed to count trading journals by user id")
	}

	return count, nil
}

func (s *TradingJournalStorage) Exists(ctx context.Context, id uuid.UUID, userID uuid.UUID) (bool, error) {
	count, err := s.db.NewSelect().
		Model((*entity.TradingJournal)(nil)).
		Where("id = ? AND user_id = ?", id, userID).
		Count(ctx)

	if err != nil {
		return false, errors.Wrap(err, "failed to check if trading journal exists")
	}

	return count > 0, nil
}
