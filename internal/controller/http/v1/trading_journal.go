package v1

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/user/normark/internal/dto"
	"github.com/user/normark/internal/dto/mapper"
	"github.com/user/normark/internal/entity"
	"go.uber.org/zap"
)

type TradingJournalService interface {
	Create(ctx context.Context, userID uuid.UUID, req *dto.CreateTradingJournalRequest) (*entity.TradingJournal, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.TradingJournal, error)
	GetByIDWithEntries(ctx context.Context, id uuid.UUID) (*entity.TradingJournal, error)
	GetUserJournals(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.TradingJournal, error)
	Update(ctx context.Context, journal *entity.TradingJournal) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	CountUserJournals(ctx context.Context, userID uuid.UUID) (int, error)
	VerifyAccess(ctx context.Context, journalID uuid.UUID, userID uuid.UUID) (bool, error)
}

type TradingJournalHandler struct {
	journalService TradingJournalService
	logger         *zap.Logger
	validate       *validator.Validate
}

func NewTradingJournalHandler(
	journalService TradingJournalService,
	logger *zap.Logger,
	validate *validator.Validate,
) *TradingJournalHandler {
	return &TradingJournalHandler{
		journalService: journalService,
		logger:         logger,
		validate:       validate,
	}
}

func (h *TradingJournalHandler) InitRoutes(group *gin.RouterGroup) {
	group.POST("", h.Create)
	group.GET("", h.List)
	group.GET("/:id", h.GetByID)
	group.GET("/:id/with-entries", h.GetByIDWithEntries)
	group.PUT("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
}

// Create godoc
// @Summary      Create a new trading journal
// @Description  Create a new trading journal for the authenticated user
// @Tags         Trading Journals
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateTradingJournalRequest true "Trading journal details"
// @Success      201 {object} dto.TradingJournalResponse "Successfully created trading journal"
// @Failure      400 {object} ErrorResponse "Invalid request body or validation failed"
// @Failure      401 {object} ErrorResponse "Unauthorized - missing or invalid token"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /api/v1/journals [post]
func (h *TradingJournalHandler) Create(c *gin.Context) {
	var req dto.CreateTradingJournalRequest

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

	userID, exists := c.Get("userID")
	if !exists {
		h.logger.Error("user id not found in context")
		newErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		h.logger.Error("invalid user id type in context")
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	journal, err := h.journalService.Create(c.Request.Context(), uid, &req)
	if err != nil {
		h.logger.Error("failed to create trading journal", zap.Error(err))
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := mapper.ToTradingJournalResponse(journal)
	c.JSON(http.StatusCreated, response)
}

// List godoc
// @Summary      List user's trading journals
// @Description  Get a paginated list of all trading journals for the authenticated user
// @Tags         Trading Journals
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        limit query int false "Maximum number of journals to return (default: 20, max: 100)"
// @Param        offset query int false "Number of journals to skip (default: 0)"
// @Success      200 {object} dto.TradingJournalListResponse "Successfully retrieved journals list"
// @Failure      401 {object} ErrorResponse "Unauthorized - missing or invalid token"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /api/v1/journals [get]
func (h *TradingJournalHandler) List(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		h.logger.Error("user id not found in context")
		newErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		h.logger.Error("invalid user id type in context")
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
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

	journals, err := h.journalService.GetUserJournals(c.Request.Context(), uid, limit, offset)
	if err != nil {
		h.logger.Error("failed to get user journals", zap.Error(err))
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	total, err := h.journalService.CountUserJournals(c.Request.Context(), uid)
	if err != nil {
		h.logger.Error("failed to count user journals", zap.Error(err))
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := &dto.TradingJournalListResponse{
		Journals: mapper.ToTradingJournalResponses(journals),
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	}

	c.JSON(http.StatusOK, response)
}

// GetByID godoc
// @Summary      Get trading journal by ID
// @Description  Retrieve a specific trading journal by its ID
// @Tags         Trading Journals
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Trading Journal ID (UUID)"
// @Success      200 {object} dto.TradingJournalResponse "Successfully retrieved trading journal"
// @Failure      400 {object} ErrorResponse "Invalid journal ID"
// @Failure      401 {object} ErrorResponse "Unauthorized - missing or invalid token"
// @Failure      404 {object} ErrorResponse "Journal not found"
// @Router       /api/v1/journals/{id} [get]
func (h *TradingJournalHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Error("invalid journal id", zap.Error(err))
		newErrorResponse(c, http.StatusBadRequest, "invalid journal id")
		return
	}

	journal, err := h.journalService.GetByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("failed to get trading journal", zap.Error(err))
		newErrorResponse(c, http.StatusNotFound, "journal not found")
		return
	}

	response := mapper.ToTradingJournalResponse(journal)
	c.JSON(http.StatusOK, response)
}

// GetByIDWithEntries godoc
// @Summary      Get trading journal with entries
// @Description  Retrieve a specific trading journal by its ID including all associated entries
// @Tags         Trading Journals
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Trading Journal ID (UUID)"
// @Success      200 {object} dto.TradingJournalWithEntriesResponse "Successfully retrieved trading journal with entries"
// @Failure      400 {object} ErrorResponse "Invalid journal ID"
// @Failure      401 {object} ErrorResponse "Unauthorized - missing or invalid token"
// @Failure      404 {object} ErrorResponse "Journal not found"
// @Router       /api/v1/journals/{id}/with-entries [get]
func (h *TradingJournalHandler) GetByIDWithEntries(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Error("invalid journal id", zap.Error(err))
		newErrorResponse(c, http.StatusBadRequest, "invalid journal id")
		return
	}

	journal, err := h.journalService.GetByIDWithEntries(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("failed to get trading journal with entries", zap.Error(err))
		newErrorResponse(c, http.StatusNotFound, "journal not found")
		return
	}

	response := mapper.ToTradingJournalWithEntriesResponse(journal)
	c.JSON(http.StatusOK, response)
}

// Update godoc
// @Summary      Update trading journal
// @Description  Update an existing trading journal's name and description
// @Tags         Trading Journals
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Trading Journal ID (UUID)"
// @Param        request body dto.UpdateTradingJournalRequest true "Updated journal details"
// @Success      200 {object} dto.TradingJournalResponse "Successfully updated trading journal"
// @Failure      400 {object} ErrorResponse "Invalid request body, validation failed, or invalid journal ID"
// @Failure      401 {object} ErrorResponse "Unauthorized - missing or invalid token"
// @Failure      404 {object} ErrorResponse "Journal not found"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /api/v1/journals/{id} [put]
func (h *TradingJournalHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Error("invalid journal id", zap.Error(err))
		newErrorResponse(c, http.StatusBadRequest, "invalid journal id")
		return
	}

	var req dto.UpdateTradingJournalRequest

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

	journal, err := h.journalService.GetByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("failed to get trading journal", zap.Error(err))
		newErrorResponse(c, http.StatusNotFound, "journal not found")
		return
	}

	journal.Name = req.Name
	journal.Description = req.Description

	if err := h.journalService.Update(c.Request.Context(), journal); err != nil {
		h.logger.Error("failed to update trading journal", zap.Error(err))
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := mapper.ToTradingJournalResponse(journal)
	c.JSON(http.StatusOK, response)
}

// Delete godoc
// @Summary      Delete trading journal
// @Description  Delete a trading journal and all its associated entries
// @Tags         Trading Journals
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Trading Journal ID (UUID)"
// @Success      200 {object} map[string]string "Successfully deleted journal"
// @Failure      400 {object} ErrorResponse "Invalid journal ID"
// @Failure      401 {object} ErrorResponse "Unauthorized - missing or invalid token"
// @Failure      500 {object} ErrorResponse "Internal server error or access denied"
// @Router       /api/v1/journals/{id} [delete]
func (h *TradingJournalHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Error("invalid journal id", zap.Error(err))
		newErrorResponse(c, http.StatusBadRequest, "invalid journal id")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		h.logger.Error("user id not found in context")
		newErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		h.logger.Error("invalid user id type in context")
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := h.journalService.Delete(c.Request.Context(), id, uid); err != nil {
		h.logger.Error("failed to delete trading journal", zap.Error(err))
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "journal deleted successfully"})
}
