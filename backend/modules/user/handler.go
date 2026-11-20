package user

import (
	"cwcu/backend/config"
	"cwcu/backend/middleware"
	"cwcu/backend/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

// Handler handles user HTTP requests
type Handler struct {
	service  *Service
	validate *validator.Validate
}

// RegisterRoutes registers user routes
func RegisterRoutes(api *echo.Group, db *pgxpool.Pool, redis *redis.Client, cfg *config.Config) {
	repo := NewRepository(db)
	svc := NewService(repo, redis, cfg)
	h := &Handler{service: svc, validate: validator.New()}

	users := api.Group("/users", middleware.JWTAuth(cfg))

	// Profile routes (all authenticated users)
	users.GET("/profile", h.GetProfile)
	users.PUT("/profile", h.UpdateProfile)

	// Admin only routes
	users.GET("", h.ListUsers, middleware.RequireRoles("admin"))
	users.POST("", h.CreateUser, middleware.RequireRoles("admin"))
	users.GET("/:id", h.GetUser, middleware.RequireRoles("admin"))
	users.PUT("/:id", h.UpdateUser, middleware.RequireRoles("admin"))
	users.DELETE("/:id", h.DeleteUser, middleware.RequireRoles("admin"))
}

func (h *Handler) GetProfile(c echo.Context) error {
	userID := c.Get("user_id").(string)
	user, err := h.service.GetProfile(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(response.StatusNotFound, response.Error("User not found", nil))
	}
	return c.JSON(response.StatusOK, response.Success("Profile retrieved", user))
}

func (h *Handler) UpdateProfile(c echo.Context) error {
	userID := c.Get("user_id").(string)
	var req UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(response.StatusBadRequest, response.Error("Invalid request", err.Error()))
	}
	if err := h.validate.Struct(req); err != nil {
		return c.JSON(response.StatusBadRequest, response.Error("Validation failed", err.Error()))
	}
	user, err := h.service.UpdateProfile(c.Request().Context(), userID, req)
	if err != nil {
		return c.JSON(response.StatusInternalServerError, response.Error("Failed to update profile", err.Error()))
	}
	return c.JSON(response.StatusOK, response.Success("Profile updated", user))
}

func (h *Handler) ListUsers(c echo.Context) error {
	users, err := h.service.ListUsers(c.Request().Context())
	if err != nil {
		return c.JSON(response.StatusInternalServerError, response.Error("Failed to fetch users", err.Error()))
	}
	return c.JSON(response.StatusOK, response.Success("Users retrieved", users))
}

func (h *Handler) CreateUser(c echo.Context) error {
	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(response.StatusBadRequest, response.Error("Invalid request", err.Error()))
	}
	if err := h.validate.Struct(req); err != nil {
		return c.JSON(response.StatusBadRequest, response.Error("Validation failed", err.Error()))
	}
	user, err := h.service.CreateUser(c.Request().Context(), req)
	if err != nil {
		return c.JSON(response.StatusConflict, response.Error(err.Error(), nil))
	}
	return c.JSON(response.StatusCreated, response.Success("User created", user))
}

func (h *Handler) GetUser(c echo.Context) error {
	id := c.Param("id")
	user, err := h.service.GetProfile(c.Request().Context(), id)
	if err != nil {
		return c.JSON(response.StatusNotFound, response.Error("User not found", nil))
	}
	return c.JSON(response.StatusOK, response.Success("User retrieved", user))
}

func (h *Handler) UpdateUser(c echo.Context) error {
	id := c.Param("id")
	var req UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(response.StatusBadRequest, response.Error("Invalid request", err.Error()))
	}
	user, err := h.service.UpdateProfile(c.Request().Context(), id, req)
	if err != nil {
		return c.JSON(response.StatusInternalServerError, response.Error("Update failed", err.Error()))
	}
	return c.JSON(response.StatusOK, response.Success("User updated", user))
}

func (h *Handler) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.DeleteUser(c.Request().Context(), id); err != nil {
		return c.JSON(response.StatusInternalServerError, response.Error("Delete failed", err.Error()))
	}
	return c.JSON(response.StatusOK, response.Success("User deleted", nil))
}
