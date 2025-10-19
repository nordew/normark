package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/user/normark/internal/service"
	"go.uber.org/zap"
)

type Handler struct {
	userService *service.UserService
	logger      *zap.Logger
	validate    *validator.Validate
	middleware  *Middleware
	rateLimiter *RateLimiter
}

func NewHandler(
	userService *service.UserService,
	logger *zap.Logger,
	middleware *Middleware,
	rateLimiter *RateLimiter,
) *Handler {
	return &Handler{
		userService: userService,
		logger:      logger,
		validate:    validator.New(),
		middleware:  middleware,
		rateLimiter: rateLimiter,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(h.rateLimiter.Limit())
	router.Use(h.middleware.CORS())
	router.Use(h.middleware.RequestLogger())

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			auth := v1.Group("/auth")
			{
				userHandler := NewUserHandler(h.userService, h.logger, h.validate)
				userHandler.InitRoutes(auth)
			}
		}
	}

	return router
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, ErrorResponse{Error: message})
}
