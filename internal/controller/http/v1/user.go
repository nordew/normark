package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/user/normark/internal/dto"
	"go.uber.org/zap"
)

type UserService interface {
	SignUp(ctx context.Context, req *dto.SignUpRequest) (*dto.AuthResponse, error)
	SignIn(ctx context.Context, req *dto.SignInRequest) (*dto.AuthResponse, error)
}

type UserHandler struct {
	userService UserService
	logger      *zap.Logger
	validate    *validator.Validate
}

func NewUserHandler(
	userService UserService,
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

// SignUp godoc
// @Summary      Register a new user
// @Description  Create a new user account with email, username and password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body dto.SignUpRequest true "User registration details"
// @Success      201 {object} dto.AuthResponse "Successfully registered user with access and refresh tokens"
// @Failure      400 {object} ErrorResponse "Invalid request body or validation failed"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /api/v1/auth/sign-up [post]
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

// SignIn godoc
// @Summary      User login
// @Description  Authenticate user with email and password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body dto.SignInRequest true "User login credentials"
// @Success      200 {object} dto.AuthResponse "Successfully authenticated with access and refresh tokens"
// @Failure      400 {object} ErrorResponse "Invalid request body or validation failed"
// @Failure      401 {object} ErrorResponse "Invalid credentials"
// @Router       /api/v1/auth/sign-in [post]
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
