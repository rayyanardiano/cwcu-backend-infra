package academic

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

func (r *Repository) FindAllCourses(ctx context.Context) ([]*Course, error) {
	q := `SELECT id, prodi_id, code, name, sks, type, description, created_at, updated_at
		  FROM courses WHERE deleted_at IS NULL ORDER BY name`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*Course
	for rows.Next() {
		c := &Course{}
		if err := rows.Scan(&c.ID, &c.ProdiID, &c.Code, &c.Name, &c.SKS, &c.Type, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}
	return courses, nil
}

func (r *Repository) FindCourseByID(ctx context.Context, id string) (*Course, error) {
	q := `SELECT id, prodi_id, code, name, sks, type, description, created_at, updated_at
		  FROM courses WHERE id = $1 AND deleted_at IS NULL`
	c := &Course{}
	err := r.db.QueryRow(ctx, q, id).Scan(&c.ID, &c.ProdiID, &c.Code, &c.Name, &c.SKS, &c.Type, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (r *Repository) CreateCourse(ctx context.Context, c *Course) error {
	q := `INSERT INTO courses (prodi_id, code, name, sks, type, description)
		  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, q, c.ProdiID, c.Code, c.Name, c.SKS, c.Type, c.Description).
		Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

// ─── Schedules ───────────────────────────────────────────────────────────────

func (r *Repository) FindAllSchedules(ctx context.Context) ([]*ClassSchedule, error) {
	q := `SELECT id, course_id, lecturer_id, academic_year_id, room, day, start_time, end_time, quota, created_at
		  FROM class_schedules ORDER BY day, start_time`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []*ClassSchedule
	for rows.Next() {
		s := &ClassSchedule{}
		if err := rows.Scan(&s.ID, &s.CourseID, &s.LecturerID, &s.AcademicYearID, &s.Room, &s.Day, &s.StartTime, &s.EndTime, &s.Quota, &s.CreatedAt); err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}

func (r *Repository) CreateSchedule(ctx context.Context, s *ClassSchedule) error {
	q := `INSERT INTO class_schedules (course_id, lecturer_id, academic_year_id, room, day, start_time, end_time, quota)
		  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at`
	return r.db.QueryRow(ctx, q, s.CourseID, s.LecturerID, s.AcademicYearID, s.Room, s.Day, s.StartTime, s.EndTime, s.Quota).
		Scan(&s.ID, &s.CreatedAt)
}

// ─── KRS ─────────────────────────────────────────────────────────────────────

func (r *Repository) FindKRSByStudent(ctx context.Context, studentID string) ([]*KRS, error) {
	q := `SELECT id, student_id, academic_year_id, status, total_sks, created_at, updated_at
		  FROM krs WHERE student_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, q, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*KRS
	for rows.Next() {
		k := &KRS{}
		rows.Scan(&k.ID, &k.StudentID, &k.AcademicYearID, &k.Status, &k.TotalSKS, &k.CreatedAt, &k.UpdatedAt)
		list = append(list, k)
	}
	return list, nil
}

func (r *Repository) CreateKRS(ctx context.Context, k *KRS) error {
	q := `INSERT INTO krs (student_id, academic_year_id, status, total_sks)
		  VALUES ($1, $2, 'submitted', $3) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, q, k.StudentID, k.AcademicYearID, k.TotalSKS).
		Scan(&k.ID, &k.CreatedAt, &k.UpdatedAt)
}

func (r *Repository) UpdateKRSStatus(ctx context.Context, id, status string) error {
	q := `UPDATE krs SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, q, status, id)
	return err
}

// ─── Grades ──────────────────────────────────────────────────────────────────

func (r *Repository) FindGradesByStudent(ctx context.Context, studentID string) ([]*Grade, error) {
	q := `SELECT id, student_id, class_id, midterm, final, grade, updated_at
		  FROM grades WHERE student_id = $1`
	rows, err := r.db.Query(ctx, q, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var grades []*Grade
	for rows.Next() {
		g := &Grade{}
		rows.Scan(&g.ID, &g.StudentID, &g.ClassID, &g.Midterm, &g.Final, &g.Grade, &g.UpdatedAt)
		grades = append(grades, g)
	}
	return grades, nil
}

func (r *Repository) UpsertGrade(ctx context.Context, g *Grade) error {
	q := `INSERT INTO grades (student_id, class_id, midterm, final, grade)
		  VALUES ($1, $2, $3, $4, $5)
		  ON CONFLICT (student_id, class_id)
		  DO UPDATE SET midterm = $3, final = $4, grade = $5, updated_at = NOW()
		  RETURNING id, updated_at`
	return r.db.QueryRow(ctx, q, g.StudentID, g.ClassID, g.Midterm, g.Final, g.Grade).
		Scan(&g.ID, &g.UpdatedAt)
}
