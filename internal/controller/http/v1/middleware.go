package v1

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/user/normark/internal/config"
	"github.com/user/normark/pkg/auth"
	"go.uber.org/zap"
)

type JWTValidator interface {
	ValidateToken(tokenString string) (*auth.Claims, error)
}

type Middleware struct {
	logger       *zap.Logger
	jwtValidator JWTValidator
	corsConfig   *config.CORS
}

func NewMiddleware(
	logger *zap.Logger,
	jwtValidator JWTValidator,
	corsConfig *config.CORS,
) *Middleware {
	return &Middleware{
		logger:       logger,
		jwtValidator: jwtValidator,
		corsConfig:   corsConfig,
	}
}

func (m *Middleware) CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     m.corsConfig.AllowOrigins,
		AllowMethods:     m.corsConfig.AllowMethods,
		AllowHeaders:     m.corsConfig.AllowHeaders,
		AllowCredentials: m.corsConfig.AllowCredentials,
		MaxAge:           time.Duration(m.corsConfig.MaxAge) * time.Second,
	})
}

func (m *Middleware) RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			m.logger.Error(
				"request failed",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Int("status", c.Writer.Status()),
			)
		}
	}
}

func (m *Middleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			m.logger.Error("missing authorization header")
			newErrorResponse(c, http.StatusUnauthorized, "missing authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			m.logger.Error("invalid authorization header format")
			newErrorResponse(c, http.StatusUnauthorized, "invalid authorization header format")
			return
		}

		tokenString := parts[1]

		claims, err := m.jwtValidator.ValidateToken(tokenString)
		if err != nil {
			m.logger.Error("invalid token", zap.Error(err))
			newErrorResponse(c, http.StatusUnauthorized, "invalid token")
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("username", claims.Username)

		c.Next()
	}
}
