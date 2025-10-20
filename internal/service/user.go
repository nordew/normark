package service

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/user/normark/internal/dto"
	"github.com/user/normark/internal/entity"
	"github.com/user/normark/pkg/auth"
	"go.uber.org/zap"
)

type UserStorage interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entity.User, error)
	Count(ctx context.Context) (int, error)
	Exists(ctx context.Context, email, username string) (bool, error)
}

type UserService struct {
	storage    UserStorage
	jwtManager *auth.JWTManager
	logger     *zap.Logger
}

func NewUserService(
	storage UserStorage,
	jwtManager *auth.JWTManager,
	logger *zap.Logger,
) *UserService {
	return &UserService{
		storage:    storage,
		jwtManager: jwtManager,
		logger:     logger,
	}
}

func (s *UserService) SignUp(ctx context.Context, req *dto.SignUpRequest) (*dto.AuthResponse, error) {
	exists, err := s.storage.Exists(ctx, req.Email, req.Username)
	if err != nil {
		s.logger.Error("failed to check user existence", zap.Error(err))
		return nil, errors.Wrap(err, "failed to check user existence")
	}

	if exists {
		return nil, entity.ErrUserAlreadyExists
	}

	user, err := entity.NewUserFromSignUp(req)
	if err != nil {
		s.logger.Error("failed to create user entity", zap.Error(err))
		return nil, errors.Wrap(err, "failed to create user entity")
	}

	if err := s.storage.Create(ctx, user); err != nil {
		s.logger.Error("failed to create user in database", zap.Error(err))
		return nil, errors.Wrap(err, "failed to create user")
	}

	tokens, err := s.jwtManager.GenerateTokenPair(
		user.ID,
		user.Email,
		user.Username,
	)
	if err != nil {
		s.logger.Error("failed to generate tokens", zap.Error(err))
		return nil, errors.Wrap(err, "failed to generate tokens")
	}

	return &dto.AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
	}, nil
}

func (s *UserService) SignIn(ctx context.Context, req *dto.SignInRequest) (*dto.AuthResponse, error) {
	user, err := s.storage.GetByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Error("failed to get user by email", zap.Error(err))
		return nil, entity.ErrInvalidCredentials
	}

	if err := user.ComparePassword(req.Password); err != nil {
		s.logger.Error("invalid password attempt", zap.String("email", req.Email))
		return nil, entity.ErrInvalidCredentials
	}

	tokens, err := s.jwtManager.GenerateTokenPair(user.ID, user.Email, user.Username)
	if err != nil {
		s.logger.Error("failed to generate tokens", zap.Error(err))
		return nil, errors.Wrap(err, "failed to generate tokens")
	}

	return &dto.AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
	}, nil
}
