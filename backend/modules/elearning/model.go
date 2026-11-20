package elearning

import "time"

// CourseElearn represents an e-learning course linked to a class schedule
type CourseElearn struct {
	ID          string    `db:"id" json:"id"`
	ClassID     string    `db:"class_id" json:"class_id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// Material represents materi kuliah
type Material struct {
	ID         string    `db:"id" json:"id"`
	CourseID   string    `db:"course_elearn_id" json:"course_id"`
	Title      string    `db:"title" json:"title"`
	Type       string    `db:"type" json:"type"` // pdf, video, link, text
	ContentURL *string   `db:"content_url" json:"content_url,omitempty"`
	Content    *string   `db:"content" json:"content,omitempty"`
	Order      int       `db:"order" json:"order"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

// Assignment represents tugas
type Assignment struct {
	ID          string     `db:"id" json:"id"`
	CourseID    string     `db:"course_elearn_id" json:"course_id"`
	Title       string     `db:"title" json:"title"`
	Description string     `db:"description" json:"description"`
	DueDate     *time.Time `db:"due_date" json:"due_date,omitempty"`
	MaxScore    float64    `db:"max_score" json:"max_score"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
}

// Submission represents submisi tugas oleh mahasiswa
type Submission struct {
	ID           string     `db:"id" json:"id"`
	AssignmentID string     `db:"assignment_id" json:"assignment_id"`
	StudentID    string     `db:"student_id" json:"student_id"`
	FileURL      *string    `db:"file_url" json:"file_url,omitempty"`
	Content      *string    `db:"content" json:"content,omitempty"`
	Score        *float64   `db:"score" json:"score,omitempty"`
	Feedback     *string    `db:"feedback" json:"feedback,omitempty"`
	SubmittedAt  time.Time  `db:"submitted_at" json:"submitted_at"`
	GradedAt     *time.Time `db:"graded_at" json:"graded_at,omitempty"`
}

// Quiz represents kuis
type Quiz struct {
	ID        string     `db:"id" json:"id"`
	CourseID  string     `db:"course_elearn_id" json:"course_id"`
	Title     string     `db:"title" json:"title"`
	Duration  int        `db:"duration" json:"duration"` // in minutes
	StartAt   *time.Time `db:"start_at" json:"start_at,omitempty"`
	EndAt     *time.Time `db:"end_at" json:"end_at,omitempty"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
}

// QuizQuestion represents a quiz question
type QuizQuestion struct {
	ID            string `db:"id" json:"id"`
	QuizID        string `db:"quiz_id" json:"quiz_id"`
	Question      string `db:"question" json:"question"`
	Options       []byte `db:"options" json:"options"` // JSONB: [{label, text}]
	CorrectAnswer string `db:"correct_answer" json:"correct_answer,omitempty"`
	Points        int    `db:"points" json:"points"`
}

// QuizAttempt represents a student's quiz attempt
type QuizAttempt struct {
	ID          string    `db:"id" json:"id"`
	QuizID      string    `db:"quiz_id" json:"quiz_id"`
	StudentID   string    `db:"student_id" json:"student_id"`
	Score       float64   `db:"score" json:"score"`
	SubmittedAt time.Time `db:"submitted_at" json:"submitted_at"`
}

// Forum represents a discussion forum
type Forum struct {
	ID        string    `db:"id" json:"id"`
	CourseID  string    `db:"course_elearn_id" json:"course_id"`
	Title     string    `db:"title" json:"title"`
	CreatedBy string    `db:"created_by" json:"created_by"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// ForumPost represents a post in a forum
type ForumPost struct {
	ID        string     `db:"id" json:"id"`
	ForumID   string     `db:"forum_id" json:"forum_id"`
	UserID    string     `db:"user_id" json:"user_id"`
	Content   string     `db:"content" json:"content"`
	ParentID  *string    `db:"parent_id" json:"parent_id,omitempty"` // for replies
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"-"`
}

// Announcement represents pengumuman di course
type Announcement struct {
	ID        string    `db:"id" json:"id"`
	CourseID  string    `db:"course_elearn_id" json:"course_id"`
	Title     string    `db:"title" json:"title"`
	Content   string    `db:"content" json:"content"`
	CreatedBy string    `db:"created_by" json:"created_by"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
