package v1

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/user/normark/internal/types"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/user/normark/internal/dto"
	"github.com/user/normark/internal/dto/mapper"
	"github.com/user/normark/internal/entity"
	"go.uber.org/zap"
)

type TradingJournalEntryService interface {
	Create(ctx context.Context, journalID uuid.UUID, req *dto.CreateTradingJournalEntryRequest) (*entity.TradingJournalEntry, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.TradingJournalEntry, error)
	GetByIDWithJournal(ctx context.Context, id uuid.UUID) (*entity.TradingJournalEntry, error)
	GetJournalEntries(ctx context.Context, journalID uuid.UUID, limit, offset int) ([]*entity.TradingJournalEntry, error)
	GetByDateRange(ctx context.Context, journalID uuid.UUID, startDate, endDate time.Time) ([]*entity.TradingJournalEntry, error)
	GetByAsset(ctx context.Context, journalID uuid.UUID, asset types.CurrencyPair, limit, offset int) ([]*entity.TradingJournalEntry, error)
	GetBySession(ctx context.Context, journalID uuid.UUID, session types.TradingSession, limit, offset int) ([]*entity.TradingJournalEntry, error)
	GetByResult(ctx context.Context, journalID uuid.UUID, result types.TradeResult, limit, offset int) ([]*entity.TradingJournalEntry, error)
	Update(ctx context.Context, entry *entity.TradingJournalEntry) error
	Delete(ctx context.Context, id uuid.UUID, journalID uuid.UUID) error
	CountJournalEntries(ctx context.Context, journalID uuid.UUID) (int, error)
	GetStatistics(ctx context.Context, journalID uuid.UUID) (map[string]any, error)
	VerifyAccess(ctx context.Context, entryID uuid.UUID, journalID uuid.UUID) (bool, error)
}

type TradingJournalEntryHandler struct {
	entryService   TradingJournalEntryService
	journalService TradingJournalService
	logger         *zap.Logger
	validate       *validator.Validate
}

func NewTradingJournalEntryHandler(
	entryService TradingJournalEntryService,
	journalService TradingJournalService,
	logger *zap.Logger,
	validate *validator.Validate,
) *TradingJournalEntryHandler {
	return &TradingJournalEntryHandler{
		entryService:   entryService,
		journalService: journalService,
		logger:         logger,
		validate:       validate,
	}
}

func (h *TradingJournalEntryHandler) InitRoutes(group *gin.RouterGroup) {
	group.POST("", h.Create)
	group.GET("", h.List)
	group.GET("/statistics", h.GetStatistics)
	group.GET("/:entryId", h.GetByID)
	group.PUT("/:entryId", h.Update)
	group.DELETE("/:entryId", h.Delete)
}

// Create godoc
// @Summary      Create a new trading journal entry
// @Description  Create a new trade entry in a specific trading journal
// @Tags         Trading Journal Entries
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Trading Journal ID (UUID)"
// @Param        request body dto.CreateTradingJournalEntryRequest true "Trading entry details"
// @Success      201 {object} dto.TradingJournalEntryResponse "Successfully created trading entry"
// @Failure      400 {object} ErrorResponse "Invalid request body, validation failed, or invalid journal ID"
// @Failure      401 {object} ErrorResponse "Unauthorized - missing or invalid token"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /api/v1/journals/{id}/entries [post]
func (h *TradingJournalEntryHandler) Create(c *gin.Context) {
	journalIDStr := c.Param("id")
	journalID, err := uuid.Parse(journalIDStr)
	if err != nil {
		h.logger.Error("invalid journal id", zap.Error(err))
		newErrorResponse(c, http.StatusBadRequest, "invalid journal id")
		return
	}

	var req dto.CreateTradingJournalEntryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("failed to bind request", zap.Error(err))
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(&req); err != nil {
		h.logger.Error("validation failed", zap.Error(err))
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	entry, err := h.entryService.Create(c.Request.Context(), journalID, &req)
	if err != nil {
		h.logger.Error("failed to create trading journal entry", zap.Error(err))
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := mapper.ToTradingJournalEntryResponse(entry)
	c.JSON(http.StatusCreated, response)
}

// List godoc
// @Summary      List trading journal entries
// @Description  Get a paginated list of all entries for a specific trading journal
// @Tags         Trading Journal Entries
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Trading Journal ID (UUID)"
// @Param        limit query int false "Maximum number of entries to return (default: 20, max: 100)"
// @Param        offset query int false "Number of entries to skip (default: 0)"
// @Success      200 {object} dto.TradingJournalEntryListResponse "Successfully retrieved entries list"
// @Failure      400 {object} ErrorResponse "Invalid journal ID"
// @Failure      401 {object} ErrorResponse "Unauthorized - missing or invalid token"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /api/v1/journals/{id}/entries [get]
func (h *TradingJournalEntryHandler) List(c *gin.Context) {
	journalIDStr := c.Param("id")
	journalID, err := uuid.Parse(journalIDStr)
	if err != nil {
		h.logger.Error("invalid journal id", zap.Error(err))
		newErrorResponse(c, http.StatusBadRequest, "invalid journal id")
		return
	}

	limit := 20
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	entries, err := h.entryService.GetJournalEntries(c.Request.Context(), journalID, limit, offset)
	if err != nil {
		h.logger.Error("failed to get journal entries", zap.Error(err))
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	total, err := h.entryService.CountJournalEntries(c.Request.Context(), journalID)
	if err != nil {
		h.logger.Error("failed to count journal entries", zap.Error(err))
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := &dto.TradingJournalEntryListResponse{
		Entries: mapper.ToTradingJournalEntryResponses(entries),
		Total:   total,
		Limit:   limit,
		Offset:  offset,
	}

	c.JSON(http.StatusOK, response)
}

// GetByID godoc
// @Summary      Get trading journal entry by ID
// @Description  Retrieve a specific trading journal entry by its ID
// @Tags         Trading Journal Entries
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Trading Journal ID (UUID)"
// @Param        entryId path string true "Trading Entry ID (UUID)"
// @Success      200 {object} dto.TradingJournalEntryResponse "Successfully retrieved trading entry"
// @Failure      400 {object} ErrorResponse "Invalid journal ID or entry ID"
// @Failure      401 {object} ErrorResponse "Unauthorized - missing or invalid token"
// @Failure      403 {object} ErrorResponse "Access denied - entry does not belong to journal"
// @Failure      404 {object} ErrorResponse "Entry not found"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /api/v1/journals/{id}/entries/{entryId} [get]
func (h *TradingJournalEntryHandler) GetByID(c *gin.Context) {
	journalIDStr := c.Param("id")
	journalID, err := uuid.Parse(journalIDStr)
	if err != nil {
		h.logger.Error("invalid journal id", zap.Error(err))
		newErrorResponse(c, http.StatusBadRequest, "invalid journal id")
		return
	}

	entryIDStr := c.Param("entryId")
	entryID, err := uuid.Parse(entryIDStr)
	if err != nil {
		h.logger.Error("invalid entry id", zap.Error(err))
		newErrorResponse(c, http.StatusBadRequest, "invalid entry id")
		return
	}

	entryAccess, err := h.entryService.VerifyAccess(c.Request.Context(), entryID, journalID)
	if err != nil {
		h.logger.Error("failed to verify entry access", zap.Error(err))
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if !entryAccess {
		h.logger.Error("entry does not belong to journal")
		newErrorResponse(c, http.StatusForbidden, "access denied")
		return
	}

	entry, err := h.entryService.GetByID(c.Request.Context(), entryID)
	if err != nil {
		h.logger.Error("failed to get trading journal entry", zap.Error(err))
		newErrorResponse(c, http.StatusNotFound, "entry not found")
		return
	}

	response := mapper.ToTradingJournalEntryResponse(entry)
	c.JSON(http.StatusOK, response)
}

// Update godoc
// @Summary      Update trading journal entry
// @Description  Update an existing trading journal entry
// @Tags         Trading Journal Entries
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Trading Journal ID (UUID)"
// @Param        entryId path string true "Trading Entry ID (UUID)"
// @Param        request body dto.UpdateTradingJournalEntryRequest true "Updated entry details"
// @Success      200 {object} dto.TradingJournalEntryResponse "Successfully updated trading entry"
// @Failure      400 {object} ErrorResponse "Invalid request body, validation failed, invalid journal ID, or invalid entry ID"
// @Failure      401 {object} ErrorResponse "Unauthorized - missing or invalid token"
// @Failure      403 {object} ErrorResponse "Access denied - entry does not belong to journal"
// @Failure      404 {object} ErrorResponse "Entry not found"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /api/v1/journals/{id}/entries/{entryId} [put]
func (h *TradingJournalEntryHandler) Update(c *gin.Context) {
	journalIDStr := c.Param("id")
	journalID, err := uuid.Parse(journalIDStr)
	if err != nil {
		h.logger.Error("invalid journal id", zap.Error(err))
		newErrorResponse(c, http.StatusBadRequest, "invalid journal id")
		return
	}

	entryIDStr := c.Param("entryId")
	entryID, err := uuid.Parse(entryIDStr)
	if err != nil {
		h.logger.Error("invalid entry id", zap.Error(err))
		newErrorResponse(c, http.StatusBadRequest, "invalid entry id")
		return
	}

	var req dto.UpdateTradingJournalEntryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("failed to bind request", zap.Error(err))
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(&req); err != nil {
		h.logger.Error("validation failed", zap.Error(err))
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	entryAccess, err := h.entryService.VerifyAccess(c.Request.Context(), entryID, journalID)
	if err != nil {
		h.logger.Error("failed to verify entry access", zap.Error(err))
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if !entryAccess {
		h.logger.Error("entry does not belong to journal")
		newErrorResponse(c, http.StatusForbidden, "access denied")
		return
	}

	entry, err := h.entryService.GetByID(c.Request.Context(), entryID)
	if err != nil {
		h.logger.Error("failed to get trading journal entry", zap.Error(err))
		newErrorResponse(c, http.StatusNotFound, "entry not found")
		return
	}

	entry.Day = req.Day
	entry.Asset = req.Asset
	entry.LTF = req.LTF
	entry.HTF = req.HTF
	entry.EntryCharts = req.EntryCharts
	entry.Session = req.Session
	entry.TradeType = req.TradeType
	entry.Setup = req.Setup
	entry.Direction = req.Direction
	entry.EntryType = req.EntryType
	entry.Realized = req.Realized
	entry.MaxRR = req.MaxRR
	entry.Result = req.Result
	entry.Notes = req.Notes

	if err := h.entryService.Update(c.Request.Context(), entry); err != nil {
		h.logger.Error("failed to update trading journal entry", zap.Error(err))
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := mapper.ToTradingJournalEntryResponse(entry)
	c.JSON(http.StatusOK, response)
}

// Delete godoc
// @Summary      Delete trading journal entry
// @Description  Delete a specific trading journal entry
// @Tags         Trading Journal Entries
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Trading Journal ID (UUID)"
// @Param        entryId path string true "Trading Entry ID (UUID)"
// @Success      200 {object} map[string]string "Successfully deleted entry"
// @Failure      400 {object} ErrorResponse "Invalid journal ID or entry ID"
// @Failure      401 {object} ErrorResponse "Unauthorized - missing or invalid token"
// @Failure      500 {object} ErrorResponse "Internal server error or access denied"
// @Router       /api/v1/journals/{id}/entries/{entryId} [delete]
func (h *TradingJournalEntryHandler) Delete(c *gin.Context) {
	journalIDStr := c.Param("id")
	journalID, err := uuid.Parse(journalIDStr)
	if err != nil {
		h.logger.Error("invalid journal id", zap.Error(err))
		newErrorResponse(c, http.StatusBadRequest, "invalid journal id")
		return
	}

	entryIDStr := c.Param("entryId")
	entryID, err := uuid.Parse(entryIDStr)
	if err != nil {
		h.logger.Error("invalid entry id", zap.Error(err))
		newErrorResponse(c, http.StatusBadRequest, "invalid entry id")
		return
	}

	if err := h.entryService.Delete(c.Request.Context(), entryID, journalID); err != nil {
		h.logger.Error("failed to delete trading journal entry", zap.Error(err))
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "entry deleted successfully"})
}

// GetStatistics godoc
// @Summary      Get trading journal statistics
// @Description  Retrieve statistical data for a specific trading journal including win rate, total trades, and performance metrics
// @Tags         Trading Journal Entries
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Trading Journal ID (UUID)"
// @Success      200 {object} dto.TradingJournalStatisticsResponse "Successfully retrieved journal statistics"
// @Failure      400 {object} ErrorResponse "Invalid journal ID"
// @Failure      401 {object} ErrorResponse "Unauthorized - missing or invalid token"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /api/v1/journals/{id}/entries/statistics [get]
func (h *TradingJournalEntryHandler) GetStatistics(c *gin.Context) {
	journalIDStr := c.Param("id")
	journalID, err := uuid.Parse(journalIDStr)
	if err != nil {
		h.logger.Error("invalid journal id", zap.Error(err))
		newErrorResponse(c, http.StatusBadRequest, "invalid journal id")
		return
	}

	stats, err := h.entryService.GetStatistics(c.Request.Context(), journalID)
	if err != nil {
		h.logger.Error("failed to get journal statistics", zap.Error(err))
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := mapper.ToStatisticsResponse(stats)
	c.JSON(http.StatusOK, response)
}
