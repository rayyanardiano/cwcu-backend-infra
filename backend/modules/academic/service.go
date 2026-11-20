package academic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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

// ─── Courses ─────────────────────────────────────────────────────────────────

func (s *Service) ListCourses(ctx context.Context) ([]*Course, error) {
	cacheKey := "courses:all"
	cached, err := s.redis.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var courses []*Course
		if json.Unmarshal(cached, &courses) == nil {
			return courses, nil
		}
	}

	courses, err := s.repo.FindAllCourses(ctx)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(courses); err == nil {
		s.redis.Set(ctx, cacheKey, data, 30*time.Minute)
	}

	return courses, nil
}

func (s *Service) GetCourse(ctx context.Context, id string) (*Course, error) {
	return s.repo.FindCourseByID(ctx, id)
}

func (s *Service) CreateCourse(ctx context.Context, req CourseRequest) (*Course, error) {
	course := &Course{
		ProdiID:     req.ProdiID,
		Code:        req.Code,
		Name:        req.Name,
		SKS:         req.SKS,
		Type:        req.Type,
		Description: req.Description,
	}
	if err := s.repo.CreateCourse(ctx, course); err != nil {
		return nil, err
	}
	// Invalidate cache
	s.redis.Del(ctx, "courses:all")
	return course, nil
}

// ─── Schedules ───────────────────────────────────────────────────────────────

func (s *Service) ListSchedules(ctx context.Context) ([]*ClassSchedule, error) {
	cacheKey := "schedules:all"
	cached, err := s.redis.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var schedules []*ClassSchedule
		if json.Unmarshal(cached, &schedules) == nil {
			return schedules, nil
		}
	}

	schedules, err := s.repo.FindAllSchedules(ctx)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(schedules); err == nil {
		s.redis.Set(ctx, cacheKey, data, 1*time.Hour)
	}

	return schedules, nil
}

func (s *Service) CreateSchedule(ctx context.Context, req ScheduleRequest) (*ClassSchedule, error) {
	schedule := &ClassSchedule{
		CourseID:       req.CourseID,
		LecturerID:     req.LecturerID,
		AcademicYearID: req.AcademicYearID,
		Room:           req.Room,
		Day:            req.Day,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		Quota:          req.Quota,
	}
	if err := s.repo.CreateSchedule(ctx, schedule); err != nil {
		return nil, err
	}
	s.redis.Del(ctx, "schedules:all")
	return schedule, nil
}

// ─── KRS ─────────────────────────────────────────────────────────────────────

func (s *Service) GetStudentKRS(ctx context.Context, studentID string) ([]*KRS, error) {
	return s.repo.FindKRSByStudent(ctx, studentID)
}

func (s *Service) SubmitKRS(ctx context.Context, studentID string, req KRSRequest) (*KRS, error) {
	krs := &KRS{
		StudentID:      studentID,
		AcademicYearID: req.AcademicYearID,
		Status:         "submitted",
		// TODO: calculate total SKS from class IDs
	}
	if err := s.repo.CreateKRS(ctx, krs); err != nil {
		return nil, err
	}

	// Publish event to LavinMQ
	payload, _ := json.Marshal(map[string]string{
		"student_id": studentID,
		"krs_id":     krs.ID,
		"event":      "krs.submitted",
	})
	s.mq.Publish(ctx, "notification.push", payload)

	return krs, nil
}

func (s *Service) ApproveKRS(ctx context.Context, krsID string) error {
	if err := s.repo.UpdateKRSStatus(ctx, krsID, "approved"); err != nil {
		return err
	}

	// Publish approval notification
	payload, _ := json.Marshal(map[string]string{
		"krs_id": krsID,
		"event":  "krs.approved",
	})
	s.mq.Publish(ctx, "notification.push", payload)
	return nil
}

// ─── Grades ──────────────────────────────────────────────────────────────────

func (s *Service) GetGrades(ctx context.Context, studentID string) ([]*Grade, error) {
	cacheKey := fmt.Sprintf("grades:%s", studentID)
	cached, err := s.redis.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var grades []*Grade
		if json.Unmarshal(cached, &grades) == nil {
			return grades, nil
		}
	}
	grades, err := s.repo.FindGradesByStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}
	if data, err := json.Marshal(grades); err == nil {
		s.redis.Set(ctx, cacheKey, data, 15*time.Minute)
	}
	return grades, nil
}

func (s *Service) InputGrade(ctx context.Context, req GradeRequest) (*Grade, error) {
	grade := &Grade{
		StudentID: req.StudentID,
		ClassID:   req.ClassID,
		Midterm:   req.Midterm,
		Final:     req.Final,
		Grade:     calculateLetterGrade(req.Midterm, req.Final),
	}
	if err := s.repo.UpsertGrade(ctx, grade); err != nil {
		return nil, err
	}
	// Invalidate grade cache
	s.redis.Del(ctx, fmt.Sprintf("grades:%s", req.StudentID))

	// Publish grade updated event
	payload, _ := json.Marshal(map[string]string{
		"student_id": req.StudentID,
		"class_id":   req.ClassID,
		"event":      "grade.updated",
	})
	s.mq.Publish(ctx, "grade.updated", payload)
	return grade, nil
}

func (s *Service) GetTranscript(ctx context.Context, studentID string) ([]*Grade, error) {
	return s.repo.FindGradesByStudent(ctx, studentID)
}

// calculateLetterGrade converts numeric score to letter grade
func calculateLetterGrade(midterm, final float64) string {
	total := midterm*0.4 + final*0.6
	switch {
	case total >= 85:
		return "A"
	case total >= 75:
		return "B+"
	case total >= 70:
		return "B"
	case total >= 65:
		return "C+"
	case total >= 60:
		return "C"
	case total >= 50:
		return "D"
	default:
		return "E"
	}
}
