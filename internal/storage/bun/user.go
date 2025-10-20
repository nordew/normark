package bun

import (
	"context"
	"database/sql"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/user/normark/internal/entity"
)

type UserStorage struct {
	db *bun.DB
}

func NewUserStorage(db *bun.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (s *UserStorage) Create(ctx context.Context, user *entity.User) error {
	_, err := s.db.NewInsert().
		Model(user).
		Exec(ctx)

	if err != nil {
		return errors.Wrap(err, "failed to create user")
	}

	return nil
}

func (s *UserStorage) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	user := new(entity.User)

	err := s.db.NewSelect().
		Model(user).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(err, "user not found")
		}
		return nil, errors.Wrap(err, "failed to get user by id")
	}

	return user, nil
}

func (s *UserStorage) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := new(entity.User)

	err := s.db.NewSelect().
		Model(user).
		Where("email = ?", email).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(err, "user not found")
		}
		return nil, errors.Wrap(err, "failed to get user by email")
	}

	return user, nil
}

func (s *UserStorage) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	user := new(entity.User)

	err := s.db.NewSelect().
		Model(user).
		Where("username = ?", username).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(err, "user not found")
		}
		return nil, errors.Wrap(err, "failed to get user by username")
	}

	return user, nil
}

func (s *UserStorage) Update(ctx context.Context, user *entity.User) error {
	result, err := s.db.NewUpdate().
		Model(user).
		WherePK().
		Exec(ctx)

	if err != nil {
		return errors.Wrap(err, "failed to update user")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (s *UserStorage) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := s.db.NewDelete().
		Model((*entity.User)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	if err != nil {
		return errors.Wrap(err, "failed to delete user")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (s *UserStorage) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	var users []*entity.User

	err := s.db.NewSelect().
		Model(&users).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Scan(ctx)

	if err != nil {
		return nil, errors.Wrap(err, "failed to list users")
	}

	return users, nil
}

func (s *UserStorage) Count(ctx context.Context) (int, error) {
	count, err := s.db.NewSelect().
		Model((*entity.User)(nil)).
		Count(ctx)

	if err != nil {
		return 0, errors.Wrap(err, "failed to count users")
	}

	return count, nil
}

func (s *UserStorage) Exists(ctx context.Context, email, username string) (bool, error) {
	count, err := s.db.NewSelect().
		Model((*entity.User)(nil)).
		Where("email = ? OR username = ?", email, username).
		Count(ctx)

	if err != nil {
		return false, errors.Wrap(err, "failed to check if user exists")
	}

	return count > 0, nil
}
