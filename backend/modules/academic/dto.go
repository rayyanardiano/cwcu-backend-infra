package academic

type CourseRequest struct {
	ProdiID     string  `json:"prodi_id" validate:"required"`
	Code        string  `json:"code" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	SKS         int     `json:"sks" validate:"required,min=1,max=6"`
	Type        string  `json:"type" validate:"required,oneof=wajib pilihan"`
	Description *string `json:"description,omitempty"`
}

type ScheduleRequest struct {
	CourseID       string `json:"course_id" validate:"required"`
	LecturerID     string `json:"lecturer_id" validate:"required"`
	AcademicYearID string `json:"academic_year_id" validate:"required"`
	Room           string `json:"room" validate:"required"`
	Day            string `json:"day" validate:"required,oneof=Senin Selasa Rabu Kamis Jumat Sabtu"`
	StartTime      string `json:"start_time" validate:"required"`
	EndTime        string `json:"end_time" validate:"required"`
	Quota          int    `json:"quota" validate:"required,min=1"`
}

type KRSRequest struct {
	AcademicYearID string   `json:"academic_year_id" validate:"required"`
	ClassIDs       []string `json:"class_ids" validate:"required,min=1"`
}

type GradeRequest struct {
	StudentID string  `json:"student_id" validate:"required"`
	ClassID   string  `json:"class_id" validate:"required"`
	Midterm   float64 `json:"midterm" validate:"required,min=0,max=100"`
	Final     float64 `json:"final" validate:"required,min=0,max=100"`
}
