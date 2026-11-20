package elearning

import (
	"context"
	"encoding/json"

	"cwcu/backend/config"
	"cwcu/backend/pkg/lavinmq"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	repo  *Repository
	redis *redis.Client
	cfg   *config.Config
	mq    *lavinmq.Client
}

func NewService(repo *Repository, rdb *redis.Client, cfg *config.Config, mq *lavinmq.Client) *Service {
	return &Service{repo: repo, redis: rdb, cfg: cfg, mq: mq}
}

func (s *Service) ListEnrolledCourses(ctx context.Context, userID string) ([]*CourseElearn, error) {
	return s.repo.FindEnrolledCourses(ctx, userID)
}

func (s *Service) GetCourse(ctx context.Context, id string) (*CourseElearn, error) {
	return s.repo.FindCourseByID(ctx, id)
}

func (s *Service) ListMaterials(ctx context.Context, courseID string) ([]*Material, error) {
	return s.repo.FindMaterialsByCourse(ctx, courseID)
}

func (s *Service) CreateMaterial(ctx context.Context, req MaterialRequest) (*Material, error) {
	m := &Material{
		CourseID:   req.CourseID,
		Title:      req.Title,
		Type:       req.Type,
		ContentURL: req.ContentURL,
		Content:    req.Content,
		Order:      req.Order,
	}
	if err := s.repo.CreateMaterial(ctx, m); err != nil {
		return nil, err
	}

	// Notify students about new material
	payload, _ := json.Marshal(map[string]string{
		"course_id": req.CourseID,
		"title":     req.Title,
		"event":     "material.created",
	})
	s.mq.Publish(ctx, "notification.push", payload)

	return m, nil
}

func (s *Service) GetMaterial(ctx context.Context, id string) (*Material, error) {
	return s.repo.FindMaterialByID(ctx, id)
}

func (s *Service) ListAssignments(ctx context.Context, courseID string) ([]*Assignment, error) {
	return s.repo.FindAssignmentsByCourse(ctx, courseID)
}

func (s *Service) CreateAssignment(ctx context.Context, req AssignmentRequest) (*Assignment, error) {
	a := &Assignment{
		CourseID:    req.CourseID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		MaxScore:    req.MaxScore,
	}
	if err := s.repo.CreateAssignment(ctx, a); err != nil {
		return nil, err
	}

	// Notify students about new assignment
	payload, _ := json.Marshal(map[string]string{
		"course_id": req.CourseID,
		"title":     req.Title,
		"event":     "assignment.created",
	})
	s.mq.Publish(ctx, "notification.push", payload)

	return a, nil
}

func (s *Service) GetAssignment(ctx context.Context, id string) (*Assignment, error) {
	return s.repo.FindAssignmentByID(ctx, id)
}

func (s *Service) SubmitAssignment(ctx context.Context, asgID, studentID string, req SubmissionRequest) (*Submission, error) {
	sub := &Submission{
		AssignmentID: asgID,
		StudentID:    studentID,
		FileURL:      req.FileURL,
		Content:      req.Content,
	}
	if err := s.repo.CreateSubmission(ctx, sub); err != nil {
		return nil, err
	}

	// Notify lecturer about submission
	payload, _ := json.Marshal(map[string]string{
		"assignment_id": asgID,
		"student_id":    studentID,
		"event":         "assignment.submitted",
	})
	s.mq.Publish(ctx, "assignment.submitted", payload)

	return sub, nil
}

func (s *Service) ListSubmissions(ctx context.Context, asgID string) ([]*Submission, error) {
	return s.repo.FindSubmissionsByAssignment(ctx, asgID)
}

func (s *Service) GradeSubmission(ctx context.Context, id string, req GradeSubmissionRequest) (*Submission, error) {
	return s.repo.GradeSubmission(ctx, id, req.Score, req.Feedback)
}

func (s *Service) ListQuizzes(ctx context.Context, courseID string) ([]*Quiz, error) {
	return s.repo.FindQuizzesByCourse(ctx, courseID)
}

func (s *Service) CreateQuiz(ctx context.Context, req QuizRequest) (*Quiz, error) {
	qz := &Quiz{
		CourseID: req.CourseID,
		Title:    req.Title,
		Duration: req.Duration,
		StartAt:  req.StartAt,
		EndAt:    req.EndAt,
	}
	if err := s.repo.CreateQuiz(ctx, qz); err != nil {
		return nil, err
	}
	return qz, nil
}

func (s *Service) GetQuiz(ctx context.Context, id string) (*Quiz, error) {
	return s.repo.FindQuizByID(ctx, id)
}

func (s *Service) AttemptQuiz(ctx context.Context, quizID, studentID string, req QuizAttemptRequest) (*QuizAttempt, error) {
	// TODO: validate answers against correct answers and calculate score
	attempt := &QuizAttempt{
		QuizID:    quizID,
		StudentID: studentID,
		Score:     0, // will be calculated from answers
	}
	if err := s.repo.CreateAttempt(ctx, attempt); err != nil {
		return nil, err
	}
	return attempt, nil
}

func (s *Service) ListForums(ctx context.Context, courseID string) ([]*Forum, error) {
	return s.repo.FindForumsByCourse(ctx, courseID)
}

func (s *Service) CreateForum(ctx context.Context, req ForumRequest) (*Forum, error) {
	f := &Forum{
		CourseID:  req.CourseID,
		Title:     req.Title,
		CreatedBy: req.CreatedBy,
	}
	if err := s.repo.CreateForum(ctx, f); err != nil {
		return nil, err
	}
	return f, nil
}

func (s *Service) ListPosts(ctx context.Context, forumID string) ([]*ForumPost, error) {
	return s.repo.FindPostsByForum(ctx, forumID)
}

func (s *Service) CreatePost(ctx context.Context, req PostRequest) (*ForumPost, error) {
	p := &ForumPost{
		ForumID:  req.ForumID,
		UserID:   req.UserID,
		Content:  req.Content,
		ParentID: req.ParentID,
	}
	if err := s.repo.CreatePost(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) ListAnnouncements(ctx context.Context, courseID string) ([]*Announcement, error) {
	return s.repo.FindAnnouncementsByCourse(ctx, courseID)
}

func (s *Service) CreateAnnouncement(ctx context.Context, req AnnouncementRequest) (*Announcement, error) {
	a := &Announcement{
		CourseID:  req.CourseID,
		Title:     req.Title,
		Content:   req.Content,
		CreatedBy: req.CreatedBy,
	}
	if err := s.repo.CreateAnnouncement(ctx, a); err != nil {
		return nil, err
	}

	// Notify students
	payload, _ := json.Marshal(map[string]string{
		"course_id": req.CourseID,
		"title":     req.Title,
		"event":     "announcement.created",
	})
	s.mq.Publish(ctx, "notification.push", payload)

	return a, nil
}
