package auth

import (
	"net/http"
	"strings"

	"cwcu/backend/config"
	"cwcu/backend/middleware"
	"cwcu/backend/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

// Handler handles auth HTTP requests
type Handler struct {
	service  *Service
	validate *validator.Validate
}

// RegisterRoutes registers auth routes to the Echo group
func RegisterRoutes(api *echo.Group, db *pgxpool.Pool, redis *redis.Client, cfg *config.Config) {
	repo := NewRepository(db)
	svc := NewService(repo, redis, cfg)
	h := &Handler{service: svc, validate: validator.New()}

	auth := api.Group("/auth")
	auth.POST("/login", h.Login)
	auth.POST("/refresh", h.RefreshToken)
	auth.POST("/logout", h.Logout, middleware.JWTAuth(cfg))
	auth.POST("/forgot-password", h.ForgotPassword)
	auth.POST("/reset-password", h.ResetPassword)
}

// Login handles POST /api/v1/auth/login
func (h *Handler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(response.StatusBadRequest, response.Error("Invalid request body", err.Error()))
	}
	if err := h.validate.Struct(req); err != nil {
		return c.JSON(response.StatusBadRequest, response.Error("Validation failed", err.Error()))
	}

	result, err := h.service.Login(c.Request().Context(), req)
	if err != nil {
		return c.JSON(response.StatusUnauthorized, response.Error(err.Error(), nil))
	}

	return c.JSON(response.StatusOK, response.Success("Login successful", result))
}

// RefreshToken handles POST /api/v1/auth/refresh
func (h *Handler) RefreshToken(c echo.Context) error {
	var req RefreshRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(response.StatusBadRequest, response.Error("Invalid request body", err.Error()))
	}
	if err := h.validate.Struct(req); err != nil {
		return c.JSON(response.StatusBadRequest, response.Error("Validation failed", err.Error()))
	}

	result, err := h.service.RefreshToken(c.Request().Context(), req)
	if err != nil {
		return c.JSON(response.StatusUnauthorized, response.Error(err.Error(), nil))
	}

	return c.JSON(response.StatusOK, response.Success("Token refreshed", result))
}

// Logout handles POST /api/v1/auth/logout
func (h *Handler) Logout(c echo.Context) error {
	userID := c.Get("user_id").(string)
	authHeader := c.Request().Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	if err := h.service.Logout(c.Request().Context(), userID, token); err != nil {
		return c.JSON(response.StatusInternalServerError, response.Error("Logout failed", err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success("Logged out successfully", nil))
}

// ForgotPassword handles POST /api/v1/auth/forgot-password
func (h *Handler) ForgotPassword(c echo.Context) error {
	var req ForgotPasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(response.StatusBadRequest, response.Error("Invalid request body", err.Error()))
	}
	// TODO: implement email notification via LavinMQ
	return c.JSON(http.StatusOK, response.Success("If the email exists, a reset link has been sent", nil))
}

// ResetPassword handles POST /api/v1/auth/reset-password
func (h *Handler) ResetPassword(c echo.Context) error {
	var req ResetPasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(response.StatusBadRequest, response.Error("Invalid request body", err.Error()))
	}
	// TODO: implement token validation and password update
	return c.JSON(http.StatusOK, response.Success("Password reset successful", nil))
}
