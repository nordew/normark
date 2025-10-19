package auth

import (
	"time"

	"github.com/cockroachdb/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type JWTManager struct {
	secretKey          string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

func NewJWTManager(
	secretKey string,
	accessTokenExpiry, refreshTokenExpiry int,
) (*JWTManager, error) {
	if secretKey == "" {
		return nil, errors.New("secret key cannot be empty")
	}

	if len(secretKey) < 32 {
		return nil, errors.New("secret key must be at least 32 characters")
	}

	return &JWTManager{
		secretKey:          secretKey,
		accessTokenExpiry:  time.Duration(accessTokenExpiry) * time.Minute,
		refreshTokenExpiry: time.Duration(refreshTokenExpiry) * time.Minute,
	}, nil
}

func (m *JWTManager) GenerateTokenPair(
	userID uuid.UUID,
	email, username string,
) (*TokenPair, error) {
	accessToken, expiresAt, err := m.generateToken(
		userID,
		email,
		username,
		m.accessTokenExpiry,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate access token")
	}

	refreshToken, _, err := m.generateToken(
		userID,
		email,
		username,
		m.refreshTokenExpiry,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate refresh token")
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

func (m *JWTManager) generateToken(
	userID uuid.UUID,
	email, username string,
	expiry time.Duration,
) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(expiry)

	claims := &Claims{
		UserID:   userID,
		Email:    email,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "normark",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		return "", time.Time{}, errors.Wrap(err, "failed to sign token")
	}

	return tokenString, expiresAt, nil
}

func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.Newf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.secretKey), nil
		},
	)

	if err != nil {
		return nil, errors.Wrap(err, "failed to parse token")
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("failed to extract claims from token")
	}

	return claims, nil
}

func (m *JWTManager) RefreshAccessToken(refreshToken string) (string, time.Time, error) {
	claims, err := m.ValidateToken(refreshToken)
	if err != nil {
		return "", time.Time{}, errors.Wrap(err, "invalid refresh token")
	}

	accessToken, expiresAt, err := m.generateToken(
		claims.UserID,
		claims.Email,
		claims.Username,
		m.accessTokenExpiry,
	)
	if err != nil {
		return "", time.Time{}, errors.Wrap(err, "failed to generate new access token")
	}

	return accessToken, expiresAt, nil
}
