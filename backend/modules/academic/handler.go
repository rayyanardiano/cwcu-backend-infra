package academic

import (
	"cwcu/backend/config"
	"cwcu/backend/middleware"
	"cwcu/backend/pkg/lavinmq"
	"cwcu/backend/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	service  *Service
	validate *validator.Validate
}

// RegisterRoutes registers academic module routes
func RegisterRoutes(api *echo.Group, db *pgxpool.Pool, rdb *redis.Client, cfg *config.Config, mq *lavinmq.Client) {
	repo := NewRepository(db)
	svc := NewService(repo, rdb, cfg, mq)
	h := &Handler{service: svc, validate: validator.New()}

	ac := api.Group("/academic", middleware.JWTAuth(cfg))

	// Courses
	ac.GET("/courses", h.ListCourses)
	ac.POST("/courses", h.CreateCourse, middleware.RequireRoles("admin", "staff"))
	ac.GET("/courses/:id", h.GetCourse)
	ac.PUT("/courses/:id", h.UpdateCourse, middleware.RequireRoles("admin", "staff"))
	ac.DELETE("/courses/:id", h.DeleteCourse, middleware.RequireRoles("admin"))

	// Class Schedules
	ac.GET("/schedules", h.ListSchedules)
	ac.POST("/schedules", h.CreateSchedule, middleware.RequireRoles("admin", "staff"))
	ac.GET("/schedules/:id", h.GetSchedule)
	ac.PUT("/schedules/:id", h.UpdateSchedule, middleware.RequireRoles("admin", "staff"))

	// KRS (Course Registration)
	ac.GET("/krs", h.GetKRS)
	ac.POST("/krs", h.SubmitKRS, middleware.RequireRoles("student"))
	ac.PUT("/krs/:id/approve", h.ApproveKRS, middleware.RequireRoles("admin", "staff"))

	// Grades
	ac.GET("/grades", h.GetGrades)
	ac.POST("/grades", h.InputGrade, middleware.RequireRoles("lecturer"))
	ac.PUT("/grades/:id", h.UpdateGrade, middleware.RequireRoles("lecturer"))

	// Transcript
	ac.GET("/transcript", h.GetTranscript, middleware.RequireRoles("student"))
}

func (h *Handler) ListCourses(c echo.Context) error {
	courses, err := h.service.ListCourses(c.Request().Context())
	if err != nil {
		return c.JSON(response.StatusInternalServerError, response.Error("Failed to fetch courses", nil))
	}
	return c.JSON(response.StatusOK, response.Success("Courses retrieved", courses))
}

func (h *Handler) CreateCourse(c echo.Context) error {
	var req CourseRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(response.StatusBadRequest, response.Error("Invalid request", err.Error()))
	}
	if err := h.validate.Struct(req); err != nil {
		return c.JSON(response.StatusBadRequest, response.Error("Validation failed", err.Error()))
	}
	course, err := h.service.CreateCourse(c.Request().Context(), req)
	if err != nil {
		return c.JSON(response.StatusInternalServerError, response.Error("Failed to create course", err.Error()))
	}
	return c.JSON(response.StatusCreated, response.Success("Course created", course))
}

func (h *Handler) GetCourse(c echo.Context) error {
	id := c.Param("id")
	course, err := h.service.GetCourse(c.Request().Context(), id)
	if err != nil {
		return c.JSON(response.StatusNotFound, response.Error("Course not found", nil))
	}
	return c.JSON(response.StatusOK, response.Success("Course retrieved", course))
}

func (h *Handler) UpdateCourse(c echo.Context) error {
	// TODO: implement
	return c.JSON(response.StatusOK, response.Success("Course updated", nil))
}

func (h *Handler) DeleteCourse(c echo.Context) error {
	// TODO: implement
	return c.JSON(response.StatusOK, response.Success("Course deleted", nil))
}

func (h *Handler) ListSchedules(c echo.Context) error {
	schedules, err := h.service.ListSchedules(c.Request().Context())
	if err != nil {
		return c.JSON(response.StatusInternalServerError, response.Error("Failed to fetch schedules", nil))
	}
	return c.JSON(response.StatusOK, response.Success("Schedules retrieved", schedules))
}

func (h *Handler) CreateSchedule(c echo.Context) error {
	var req ScheduleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(response.StatusBadRequest, response.Error("Invalid request", err.Error()))
	}
	schedule, err := h.service.CreateSchedule(c.Request().Context(), req)
	if err != nil {
		return c.JSON(response.StatusInternalServerError, response.Error("Failed to create schedule", err.Error()))
	}
	return c.JSON(response.StatusCreated, response.Success("Schedule created", schedule))
}

func (h *Handler) GetSchedule(c echo.Context) error {
	// TODO: implement
	return c.JSON(response.StatusOK, response.Success("Schedule retrieved", nil))
}

func (h *Handler) UpdateSchedule(c echo.Context) error {
	// TODO: implement
	return c.JSON(response.StatusOK, response.Success("Schedule updated", nil))
}

func (h *Handler) GetKRS(c echo.Context) error {
	userID := c.Get("user_id").(string)
	krs, err := h.service.GetStudentKRS(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(response.StatusInternalServerError, response.Error("Failed to fetch KRS", nil))
	}
	return c.JSON(response.StatusOK, response.Success("KRS retrieved", krs))
}

func (h *Handler) SubmitKRS(c echo.Context) error {
	userID := c.Get("user_id").(string)
	var req KRSRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(response.StatusBadRequest, response.Error("Invalid request", err.Error()))
	}
	krs, err := h.service.SubmitKRS(c.Request().Context(), userID, req)
	if err != nil {
		return c.JSON(response.StatusBadRequest, response.Error(err.Error(), nil))
	}
	return c.JSON(response.StatusCreated, response.Success("KRS submitted", krs))
}

func (h *Handler) ApproveKRS(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.ApproveKRS(c.Request().Context(), id); err != nil {
		return c.JSON(response.StatusInternalServerError, response.Error("Failed to approve KRS", err.Error()))
	}
	return c.JSON(response.StatusOK, response.Success("KRS approved", nil))
}

func (h *Handler) GetGrades(c echo.Context) error {
	userID := c.Get("user_id").(string)
	grades, err := h.service.GetGrades(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(response.StatusInternalServerError, response.Error("Failed to fetch grades", nil))
	}
	return c.JSON(response.StatusOK, response.Success("Grades retrieved", grades))
}

func (h *Handler) InputGrade(c echo.Context) error {
	var req GradeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(response.StatusBadRequest, response.Error("Invalid request", err.Error()))
	}
	grade, err := h.service.InputGrade(c.Request().Context(), req)
	if err != nil {
		return c.JSON(response.StatusInternalServerError, response.Error("Failed to input grade", err.Error()))
	}
	return c.JSON(response.StatusCreated, response.Success("Grade inputted", grade))
}

func (h *Handler) UpdateGrade(c echo.Context) error {
	// TODO: implement
	return c.JSON(response.StatusOK, response.Success("Grade updated", nil))
}

func (h *Handler) GetTranscript(c echo.Context) error {
	userID := c.Get("user_id").(string)
	transcript, err := h.service.GetTranscript(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(response.StatusInternalServerError, response.Error("Failed to fetch transcript", nil))
	}
	return c.JSON(response.StatusOK, response.Success("Transcript retrieved", transcript))
}
