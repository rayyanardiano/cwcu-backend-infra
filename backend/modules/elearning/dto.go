package elearning

import "time"

type MaterialRequest struct {
	CourseID   string  `json:"course_id"`
	Title      string  `json:"title" validate:"required"`
	Type       string  `json:"type" validate:"required,oneof=pdf video link text"`
	ContentURL *string `json:"content_url,omitempty"`
	Content    *string `json:"content,omitempty"`
	Order      int     `json:"order"`
}

type AssignmentRequest struct {
	CourseID    string     `json:"course_id"`
	Title       string     `json:"title" validate:"required"`
	Description string     `json:"description" validate:"required"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	MaxScore    float64    `json:"max_score" validate:"required,min=1"`
}

type SubmissionRequest struct {
	FileURL *string `json:"file_url,omitempty"`
	Content *string `json:"content,omitempty"`
}

type GradeSubmissionRequest struct {
	Score    float64 `json:"score" validate:"required,min=0"`
	Feedback *string `json:"feedback,omitempty"`
}

type QuizRequest struct {
	CourseID  string     `json:"course_id"`
	Title     string     `json:"title" validate:"required"`
	Duration  int        `json:"duration" validate:"required,min=1"`
	StartAt   *time.Time `json:"start_at,omitempty"`
	EndAt     *time.Time `json:"end_at,omitempty"`
	Questions []struct {
		Question      string   `json:"question" validate:"required"`
		Options       []string `json:"options" validate:"required,min=2"`
		CorrectAnswer string   `json:"correct_answer" validate:"required"`
		Points        int      `json:"points" validate:"required,min=1"`
	} `json:"questions"`
}

type QuizAttemptRequest struct {
	Answers []struct {
		QuestionID string `json:"question_id" validate:"required"`
		Answer     string `json:"answer" validate:"required"`
	} `json:"answers" validate:"required"`
}

type ForumRequest struct {
	CourseID  string `json:"course_id"`
	CreatedBy string `json:"created_by"`
	Title     string `json:"title" validate:"required"`
}

type PostRequest struct {
	ForumID  string  `json:"forum_id"`
	UserID   string  `json:"user_id"`
	Content  string  `json:"content" validate:"required"`
	ParentID *string `json:"parent_id,omitempty"`
}

type AnnouncementRequest struct {
	CourseID  string `json:"course_id"`
	CreatedBy string `json:"created_by"`
	Title     string `json:"title" validate:"required"`
	Content   string `json:"content" validate:"required"`
}
