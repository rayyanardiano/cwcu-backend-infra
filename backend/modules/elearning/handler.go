package elearning

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

func RegisterRoutes(api *echo.Group, db *pgxpool.Pool, rdb *redis.Client, cfg *config.Config, mq *lavinmq.Client) {
	repo := NewRepository(db)
	svc := NewService(repo, rdb, cfg, mq)
	h := &Handler{service: svc, validate: validator.New()}

	el := api.Group("/elearning", middleware.JWTAuth(cfg))

	// Courses
	el.GET("/courses", h.ListCourses)
	el.GET("/courses/:id", h.GetCourse)

	// Materials
	el.GET("/courses/:id/materials", h.ListMaterials)
	el.POST("/courses/:id/materials", h.CreateMaterial, middleware.RequireRoles("lecturer", "admin"))
	el.GET("/materials/:id", h.GetMaterial)

	// Assignments
	el.GET("/courses/:id/assignments", h.ListAssignments)
	el.POST("/courses/:id/assignments", h.CreateAssignment, middleware.RequireRoles("lecturer", "admin"))
	el.GET("/assignments/:id", h.GetAssignment)
	el.POST("/assignments/:id/submit", h.SubmitAssignment, middleware.RequireRoles("student"))
	el.GET("/assignments/:id/submissions", h.ListSubmissions, middleware.RequireRoles("lecturer", "admin"))
	el.PUT("/submissions/:id/grade", h.GradeSubmission, middleware.RequireRoles("lecturer", "admin"))

	// Quizzes
	el.GET("/courses/:id/quizzes", h.ListQuizzes)
	el.POST("/courses/:id/quizzes", h.CreateQuiz, middleware.RequireRoles("lecturer", "admin"))
	el.GET("/quizzes/:id", h.GetQuiz)
	el.POST("/quizzes/:id/attempt", h.AttemptQuiz, middleware.RequireRoles("student"))

	// Forums
	el.GET("/courses/:id/forums", h.ListForums)
	el.POST("/courses/:id/forums", h.CreateForum, middleware.RequireRoles("lecturer", "admin"))
	el.GET("/forums/:id/posts", h.ListPosts)
	el.POST("/forums/:id/posts", h.CreatePost)

	// Announcements
	el.GET("/courses/:id/announcements", h.ListAnnouncements)
	el.POST("/courses/:id/announcements", h.CreateAnnouncement, middleware.RequireRoles("lecturer", "admin"))
}

func (h *Handler) ListCourses(c echo.Context) error {
	userID := c.Get("user_id").(string)
	courses, err := h.service.ListEnrolledCourses(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(response.StatusInternalServerError, response.Error("Failed to fetch courses", nil))
	}
	return c.JSON(response.StatusOK, response.Success("Courses retrieved", courses))
}

func (h *Handler) GetCourse(c echo.Context) error {
	id := c.Param("id")
	course, err := h.service.GetCourse(c.Request().Context(), id)
	if err != nil {
		return c.JSON(response.StatusNotFound, response.Error("Course not found", nil))
	}
	return c.JSON(response.StatusOK, response.Success("Course retrieved", course))
}

func (h *Handler) ListMaterials(c echo.Context) error {
	courseID := c.Param("id")
	materials, err := h.service.ListMaterials(c.Request().Context(), courseID)
	if err != nil {
		return c.JSON(response.StatusInternalServerError, response.Error("Failed to fetch materials", nil))
	}
	return c.JSON(response.StatusOK, response.Success("Materials retrieved", materials))
}

func (h *Handler) CreateMaterial(c echo.Context) error {
	courseID := c.Param("id")
	var req MaterialRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(response.StatusBadRequest, response.Error("Invalid request", err.Error()))
	}
	req.CourseID = courseID
	mat, err := h.service.CreateMaterial(c.Request().Context(), req)
	if err != nil {
		return c.JSON(response.StatusInternalServerError, response.Error("Failed to create material", err.Error()))
	}
	return c.JSON(response.StatusCreated, response.Success("Material created", mat))
}

func (h *Handler) GetMaterial(c echo.Context) error {
	id := c.Param("id")
	mat, err := h.service.GetMaterial(c.Request().Context(), id)
	if err != nil {
		return c.JSON(response.StatusNotFound, response.Error("Material not found", nil))
	}
	return c.JSON(response.StatusOK, response.Success("Material retrieved", mat))
}

func (h *Handler) ListAssignments(c echo.Context) error {
	courseID := c.Param("id")
	assignments, err := h.service.ListAssignments(c.Request().Context(), courseID)
	if err != nil {
		return c.JSON(response.StatusInternalServerError, response.Error("Failed to fetch assignments", nil))
	}
	return c.JSON(response.StatusOK, response.Success("Assignments retrieved", assignments))
}

func (h *Handler) CreateAssignment(c echo.Context) error {
	courseID := c.Param("id")
	var req AssignmentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(response.StatusBadRequest, response.Error("Invalid request", err.Error()))
	}
	req.CourseID = courseID
	asg, err := h.service.CreateAssignment(c.Request().Context(), req)
	if err != nil {
		return c.JSON(response.StatusInternalServerError, response.Error("Failed to create assignment", err.Error()))
	}
	return c.JSON(response.StatusCreated, response.Success("Assignment created", asg))
}

func (h *Handler) GetAssignment(c echo.Context) error {
	id := c.Param("id")
	asg, err := h.service.GetAssignment(c.Request().Context(), id)
	if err != nil {
		return c.JSON(response.StatusNotFound, response.Error("Assignment not found", nil))
	}
	return c.JSON(response.StatusOK, response.Success("Assignment retrieved", asg))
}

func (h *Handler) SubmitAssignment(c echo.Context) error {
	asgID := c.Param("id")
	studentID := c.Get("user_id").(string)
	var req SubmissionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(response.StatusBadRequest, response.Error("Invalid request", err.Error()))
	}
	sub, err := h.service.SubmitAssignment(c.Request().Context(), asgID, studentID, req)
	if err != nil {
		return c.JSON(response.StatusInternalServerError, response.Error("Failed to submit assignment", err.Error()))
	}
	return c.JSON(response.StatusCreated, response.Success("Assignment submitted", sub))
}

func (h *Handler) ListSubmissions(c echo.Context) error {
	asgID := c.Param("id")
	subs, err := h.service.ListSubmissions(c.Request().Context(), asgID)
	if err != nil {
		return c.JSON(response.StatusInternalServerError, response.Error("Failed to fetch submissions", nil))
	}
	return c.JSON(response.StatusOK, response.Success("Submissions retrieved", subs))
}

func (h *Handler) GradeSubmission(c echo.Context) error {
	id := c.Param("id")
	var req GradeSubmissionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(response.StatusBadRequest, response.Error("Invalid request", err.Error()))
	}
	sub, err := h.service.GradeSubmission(c.Request().Context(), id, req)
	if err != nil {
		return c.JSON(response.StatusInternalServerError, response.Error("Failed to grade submission", err.Error()))
	}
	return c.JSON(response.StatusOK, response.Success("Submission graded", sub))
}

func (h *Handler) ListQuizzes(c echo.Context) error {
	courseID := c.Param("id")
	quizzes, err := h.service.ListQuizzes(c.Request().Context(), courseID)
	if err != nil {
		return c.JSON(response.StatusInternalServerError, response.Error("Failed to fetch quizzes", nil))
	}
	return c.JSON(response.StatusOK, response.Success("Quizzes retrieved", quizzes))
}

func (h *Handler) CreateQuiz(c echo.Context) error {
	courseID := c.Param("id")
	var req QuizRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(response.StatusBadRequest, response.Error("Invalid request", err.Error()))
	}
	req.CourseID = courseID
	quiz, err := h.service.CreateQuiz(c.Request().Context(), req)
	if err != nil {
		return c.JSON(response.StatusInternalServerError, response.Error("Failed to create quiz", err.Error()))
	}
	return c.JSON(response.StatusCreated, response.Success("Quiz created", quiz))
}

func (h *Handler) GetQuiz(c echo.Context) error {
	id := c.Param("id")
	quiz, err := h.service.GetQuiz(c.Request().Context(), id)
	if err != nil {
		return c.JSON(response.StatusNotFound, response.Error("Quiz not found", nil))
	}
	return c.JSON(response.StatusOK, response.Success("Quiz retrieved", quiz))
}

func (h *Handler) AttemptQuiz(c echo.Context) error {
	quizID := c.Param("id")
	studentID := c.Get("user_id").(string)
	var req QuizAttemptRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(response.StatusBadRequest, response.Error("Invalid request", err.Error()))
	}
	attempt, err := h.service.AttemptQuiz(c.Request().Context(), quizID, studentID, req)
	if err != nil {
		return c.JSON(response.StatusBadRequest, response.Error(err.Error(), nil))
	}
	return c.JSON(response.StatusCreated, response.Success("Quiz submitted", attempt))
}

func (h *Handler) ListForums(c echo.Context) error {
	courseID := c.Param("id")
	forums, err := h.service.ListForums(c.Request().Context(), courseID)
	if err != nil {
		return c.JSON(response.StatusInternalServerError, response.Error("Failed to fetch forums", nil))
	}
	return c.JSON(response.StatusOK, response.Success("Forums retrieved", forums))
}

func (h *Handler) CreateForum(c echo.Context) error {
	courseID := c.Param("id")
	userID := c.Get("user_id").(string)
	var req ForumRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(response.StatusBadRequest, response.Error("Invalid request", err.Error()))
	}
	req.CourseID = courseID
	req.CreatedBy = userID
	forum, err := h.service.CreateForum(c.Request().Context(), req)
	if err != nil {
		return c.JSON(response.StatusInternalServerError, response.Error("Failed to create forum", err.Error()))
	}
	return c.JSON(response.StatusCreated, response.Success("Forum created", forum))
}

func (h *Handler) ListPosts(c echo.Context) error {
	forumID := c.Param("id")
	posts, err := h.service.ListPosts(c.Request().Context(), forumID)
	if err != nil {
		return c.JSON(response.StatusInternalServerError, response.Error("Failed to fetch posts", nil))
	}
	return c.JSON(response.StatusOK, response.Success("Posts retrieved", posts))
}

func (h *Handler) CreatePost(c echo.Context) error {
	forumID := c.Param("id")
	userID := c.Get("user_id").(string)
	var req PostRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(response.StatusBadRequest, response.Error("Invalid request", err.Error()))
	}
	req.ForumID = forumID
	req.UserID = userID
	post, err := h.service.CreatePost(c.Request().Context(), req)
	if err != nil {
		return c.JSON(response.StatusInternalServerError, response.Error("Failed to create post", err.Error()))
	}
	return c.JSON(response.StatusCreated, response.Success("Post created", post))
}

func (h *Handler) ListAnnouncements(c echo.Context) error {
	courseID := c.Param("id")
	announcements, err := h.service.ListAnnouncements(c.Request().Context(), courseID)
	if err != nil {
		return c.JSON(response.StatusInternalServerError, response.Error("Failed to fetch announcements", nil))
	}
	return c.JSON(response.StatusOK, response.Success("Announcements retrieved", announcements))
}

func (h *Handler) CreateAnnouncement(c echo.Context) error {
	courseID := c.Param("id")
	userID := c.Get("user_id").(string)
	var req AnnouncementRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(response.StatusBadRequest, response.Error("Invalid request", err.Error()))
	}
	req.CourseID = courseID
	req.CreatedBy = userID
	ann, err := h.service.CreateAnnouncement(c.Request().Context(), req)
	if err != nil {
		return c.JSON(response.StatusInternalServerError, response.Error("Failed to create announcement", err.Error()))
	}
	return c.JSON(response.StatusCreated, response.Success("Announcement created", ann))
}
