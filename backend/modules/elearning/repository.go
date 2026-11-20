package elearning

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// ─── Courses ─────────────────────────────────────────────────────────────────

func (r *Repository) FindEnrolledCourses(ctx context.Context, userID string) ([]*CourseElearn, error) {
	q := `SELECT ce.id, ce.class_id, ce.title, ce.description, ce.created_at, ce.updated_at
		  FROM courses_elearn ce
		  INNER JOIN krs_details kd ON kd.class_id = ce.class_id
		  INNER JOIN krs k ON k.id = kd.krs_id
		  WHERE k.student_id = $1 AND k.status = 'approved'`
	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*CourseElearn
	for rows.Next() {
		c := &CourseElearn{}
		rows.Scan(&c.ID, &c.ClassID, &c.Title, &c.Description, &c.CreatedAt, &c.UpdatedAt)
		courses = append(courses, c)
	}
	return courses, nil
}

func (r *Repository) FindCourseByID(ctx context.Context, id string) (*CourseElearn, error) {
	q := `SELECT id, class_id, title, description, created_at, updated_at FROM courses_elearn WHERE id = $1`
	c := &CourseElearn{}
	err := r.db.QueryRow(ctx, q, id).Scan(&c.ID, &c.ClassID, &c.Title, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

// ─── Materials ───────────────────────────────────────────────────────────────

func (r *Repository) FindMaterialsByCourse(ctx context.Context, courseID string) ([]*Material, error) {
	q := `SELECT id, course_elearn_id, title, type, content_url, content, "order", created_at
		  FROM materials WHERE course_elearn_id = $1 ORDER BY "order"`
	rows, err := r.db.Query(ctx, q, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var materials []*Material
	for rows.Next() {
		m := &Material{}
		rows.Scan(&m.ID, &m.CourseID, &m.Title, &m.Type, &m.ContentURL, &m.Content, &m.Order, &m.CreatedAt)
		materials = append(materials, m)
	}
	return materials, nil
}

func (r *Repository) CreateMaterial(ctx context.Context, m *Material) error {
	q := `INSERT INTO materials (course_elearn_id, title, type, content_url, content, "order")
		  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`
	return r.db.QueryRow(ctx, q, m.CourseID, m.Title, m.Type, m.ContentURL, m.Content, m.Order).
		Scan(&m.ID, &m.CreatedAt)
}

func (r *Repository) FindMaterialByID(ctx context.Context, id string) (*Material, error) {
	q := `SELECT id, course_elearn_id, title, type, content_url, content, "order", created_at FROM materials WHERE id = $1`
	m := &Material{}
	err := r.db.QueryRow(ctx, q, id).Scan(&m.ID, &m.CourseID, &m.Title, &m.Type, &m.ContentURL, &m.Content, &m.Order, &m.CreatedAt)
	return m, err
}

// ─── Assignments ─────────────────────────────────────────────────────────────

func (r *Repository) FindAssignmentsByCourse(ctx context.Context, courseID string) ([]*Assignment, error) {
	q := `SELECT id, course_elearn_id, title, description, due_date, max_score, created_at
		  FROM assignments WHERE course_elearn_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, q, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Assignment
	for rows.Next() {
		a := &Assignment{}
		rows.Scan(&a.ID, &a.CourseID, &a.Title, &a.Description, &a.DueDate, &a.MaxScore, &a.CreatedAt)
		list = append(list, a)
	}
	return list, nil
}

func (r *Repository) CreateAssignment(ctx context.Context, a *Assignment) error {
	q := `INSERT INTO assignments (course_elearn_id, title, description, due_date, max_score)
		  VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	return r.db.QueryRow(ctx, q, a.CourseID, a.Title, a.Description, a.DueDate, a.MaxScore).
		Scan(&a.ID, &a.CreatedAt)
}

func (r *Repository) FindAssignmentByID(ctx context.Context, id string) (*Assignment, error) {
	q := `SELECT id, course_elearn_id, title, description, due_date, max_score, created_at FROM assignments WHERE id = $1`
	a := &Assignment{}
	err := r.db.QueryRow(ctx, q, id).Scan(&a.ID, &a.CourseID, &a.Title, &a.Description, &a.DueDate, &a.MaxScore, &a.CreatedAt)
	return a, err
}

// ─── Submissions ─────────────────────────────────────────────────────────────

func (r *Repository) CreateSubmission(ctx context.Context, s *Submission) error {
	q := `INSERT INTO submissions (assignment_id, student_id, file_url, content)
		  VALUES ($1, $2, $3, $4) RETURNING id, submitted_at`
	return r.db.QueryRow(ctx, q, s.AssignmentID, s.StudentID, s.FileURL, s.Content).
		Scan(&s.ID, &s.SubmittedAt)
}

func (r *Repository) FindSubmissionsByAssignment(ctx context.Context, assignmentID string) ([]*Submission, error) {
	q := `SELECT id, assignment_id, student_id, file_url, content, score, feedback, submitted_at, graded_at
		  FROM submissions WHERE assignment_id = $1`
	rows, err := r.db.Query(ctx, q, assignmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Submission
	for rows.Next() {
		s := &Submission{}
		rows.Scan(&s.ID, &s.AssignmentID, &s.StudentID, &s.FileURL, &s.Content, &s.Score, &s.Feedback, &s.SubmittedAt, &s.GradedAt)
		list = append(list, s)
	}
	return list, nil
}

func (r *Repository) GradeSubmission(ctx context.Context, id string, score float64, feedback *string) (*Submission, error) {
	q := `UPDATE submissions SET score = $1, feedback = $2, graded_at = NOW()
		  WHERE id = $3 RETURNING id, assignment_id, student_id, score, feedback, submitted_at, graded_at`
	s := &Submission{}
	err := r.db.QueryRow(ctx, q, score, feedback, id).
		Scan(&s.ID, &s.AssignmentID, &s.StudentID, &s.Score, &s.Feedback, &s.SubmittedAt, &s.GradedAt)
	return s, err
}

// ─── Quizzes ─────────────────────────────────────────────────────────────────

func (r *Repository) FindQuizzesByCourse(ctx context.Context, courseID string) ([]*Quiz, error) {
	q := `SELECT id, course_elearn_id, title, duration, start_at, end_at, created_at FROM quizzes WHERE course_elearn_id = $1`
	rows, err := r.db.Query(ctx, q, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Quiz
	for rows.Next() {
		qz := &Quiz{}
		rows.Scan(&qz.ID, &qz.CourseID, &qz.Title, &qz.Duration, &qz.StartAt, &qz.EndAt, &qz.CreatedAt)
		list = append(list, qz)
	}
	return list, nil
}

func (r *Repository) CreateQuiz(ctx context.Context, qz *Quiz) error {
	q := `INSERT INTO quizzes (course_elearn_id, title, duration, start_at, end_at)
		  VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	return r.db.QueryRow(ctx, q, qz.CourseID, qz.Title, qz.Duration, qz.StartAt, qz.EndAt).
		Scan(&qz.ID, &qz.CreatedAt)
}

func (r *Repository) FindQuizByID(ctx context.Context, id string) (*Quiz, error) {
	q := `SELECT id, course_elearn_id, title, duration, start_at, end_at, created_at FROM quizzes WHERE id = $1`
	qz := &Quiz{}
	err := r.db.QueryRow(ctx, q, id).Scan(&qz.ID, &qz.CourseID, &qz.Title, &qz.Duration, &qz.StartAt, &qz.EndAt, &qz.CreatedAt)
	return qz, err
}

func (r *Repository) CreateAttempt(ctx context.Context, a *QuizAttempt) error {
	q := `INSERT INTO quiz_attempts (quiz_id, student_id, score) VALUES ($1, $2, $3) RETURNING id, submitted_at`
	return r.db.QueryRow(ctx, q, a.QuizID, a.StudentID, a.Score).Scan(&a.ID, &a.SubmittedAt)
}

// ─── Forums ──────────────────────────────────────────────────────────────────

func (r *Repository) FindForumsByCourse(ctx context.Context, courseID string) ([]*Forum, error) {
	q := `SELECT id, course_elearn_id, title, created_by, created_at FROM forums WHERE course_elearn_id = $1`
	rows, err := r.db.Query(ctx, q, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Forum
	for rows.Next() {
		f := &Forum{}
		rows.Scan(&f.ID, &f.CourseID, &f.Title, &f.CreatedBy, &f.CreatedAt)
		list = append(list, f)
	}
	return list, nil
}

func (r *Repository) CreateForum(ctx context.Context, f *Forum) error {
	q := `INSERT INTO forums (course_elearn_id, title, created_by) VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.db.QueryRow(ctx, q, f.CourseID, f.Title, f.CreatedBy).Scan(&f.ID, &f.CreatedAt)
}

func (r *Repository) FindPostsByForum(ctx context.Context, forumID string) ([]*ForumPost, error) {
	q := `SELECT id, forum_id, user_id, content, parent_id, created_at, updated_at
		  FROM forum_posts WHERE forum_id = $1 AND deleted_at IS NULL ORDER BY created_at`
	rows, err := r.db.Query(ctx, q, forumID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*ForumPost
	for rows.Next() {
		p := &ForumPost{}
		rows.Scan(&p.ID, &p.ForumID, &p.UserID, &p.Content, &p.ParentID, &p.CreatedAt, &p.UpdatedAt)
		list = append(list, p)
	}
	return list, nil
}

func (r *Repository) CreatePost(ctx context.Context, p *ForumPost) error {
	q := `INSERT INTO forum_posts (forum_id, user_id, content, parent_id)
		  VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, q, p.ForumID, p.UserID, p.Content, p.ParentID).
		Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

// ─── Announcements ───────────────────────────────────────────────────────────

func (r *Repository) FindAnnouncementsByCourse(ctx context.Context, courseID string) ([]*Announcement, error) {
	q := `SELECT id, course_elearn_id, title, content, created_by, created_at FROM announcements WHERE course_elearn_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, q, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Announcement
	for rows.Next() {
		a := &Announcement{}
		rows.Scan(&a.ID, &a.CourseID, &a.Title, &a.Content, &a.CreatedBy, &a.CreatedAt)
		list = append(list, a)
	}
	return list, nil
}

func (r *Repository) CreateAnnouncement(ctx context.Context, a *Announcement) error {
	q := `INSERT INTO announcements (course_elearn_id, title, content, created_by)
		  VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	return r.db.QueryRow(ctx, q, a.CourseID, a.Title, a.Content, a.CreatedBy).Scan(&a.ID, &a.CreatedAt)
}
