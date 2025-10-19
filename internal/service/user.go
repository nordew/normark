package service

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/user/normark/internal/dto"
	"github.com/user/normark/internal/entity"
	"github.com/user/normark/internal/storage"
	"github.com/user/normark/pkg/auth"
	"go.uber.org/zap"
)

type UserService struct {
	storage    *storage.UserStorage
	jwtManager *auth.JWTManager
	logger     *zap.Logger
}

func NewUserService(
	storage *storage.UserStorage,
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
		return nil, errors.New("user with this email or username already exists")
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
		return nil, errors.Wrap(err, "invalid email or password")
	}

	if err := user.ComparePassword(req.Password); err != nil {
		s.logger.Error("invalid password attempt", zap.String("email", req.Email))
		return nil, errors.New("invalid email or password")
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
