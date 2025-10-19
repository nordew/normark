package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/user/normark/internal/dto"
	"github.com/user/normark/internal/service"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService *service.UserService
	logger      *zap.Logger
	validate    *validator.Validate
}

func NewUserHandler(
	userService *service.UserService,
	logger *zap.Logger,
	validate *validator.Validate,
) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
		validate:    validate,
	}
}

func (h *UserHandler) InitRoutes(group *gin.RouterGroup) {
	group.POST("/sign-up", h.SignUp)
	group.POST("/sign-in", h.SignIn)
}

func (h *UserHandler) SignUp(c *gin.Context) {
	var req dto.SignUpRequest

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

	response, err := h.userService.SignUp(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("failed to sign up user", zap.Error(err))
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *UserHandler) SignIn(c *gin.Context) {
	var req dto.SignInRequest

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

	response, err := h.userService.SignIn(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("failed to sign in user", zap.Error(err))
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}
