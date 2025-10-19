package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Handler struct {
	userService                UserService
	tradingJournalService      TradingJournalService
	tradingJournalEntryService TradingJournalEntryService
	logger                     *zap.Logger
	validate                   *validator.Validate
	middleware                 *Middleware
	rateLimiter                *RateLimiter
}

func NewHandler(
	userService UserService,
	tradingJournalService TradingJournalService,
	tradingJournalEntryService TradingJournalEntryService,
	logger *zap.Logger,
	middleware *Middleware,
	rateLimiter *RateLimiter,
) *Handler {
	return &Handler{
		userService:                userService,
		tradingJournalService:      tradingJournalService,
		tradingJournalEntryService: tradingJournalEntryService,
		logger:                     logger,
		validate:                   validator.New(),
		middleware:                 middleware,
		rateLimiter:                rateLimiter,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	h.setupMiddleware(router)

	api := router.Group("/api/v1")
	{
		h.initPublicRoutes(api)
		h.initAuthenticatedRoutes(api)
	}

	return router
}

func (h *Handler) setupMiddleware(router *gin.Engine) {
	router.Use(gin.Recovery())
	router.Use(h.rateLimiter.Limit())
	router.Use(h.middleware.CORS())
	router.Use(h.middleware.RequestLogger())
}

func (h *Handler) initPublicRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		userHandler := NewUserHandler(h.userService, h.logger, h.validate)
		userHandler.InitRoutes(auth)
	}
}

func (h *Handler) initAuthenticatedRoutes(api *gin.RouterGroup) {
	authenticated := api.Group("")
	authenticated.Use(h.middleware.Auth())
	{
		h.initJournalRoutes(authenticated)
	}
}

func (h *Handler) initJournalRoutes(group *gin.RouterGroup) {
	journals := group.Group("/journals")
	{
		journalHandler := NewTradingJournalHandler(h.tradingJournalService, h.logger, h.validate)
		journalHandler.InitRoutes(journals)

		h.initJournalEntryRoutes(journals)
	}
}

func (h *Handler) initJournalEntryRoutes(journals *gin.RouterGroup) {
	entries := journals.Group("/:id/entries")
	{
		entryHandler := NewTradingJournalEntryHandler(
			h.tradingJournalEntryService,
			h.tradingJournalService,
			h.logger,
			h.validate,
		)
		entryHandler.InitRoutes(entries)
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, ErrorResponse{Error: message})
}
