package academic

import "time"

// Course represents a mata kuliah
type Course struct {
	ID          string    `db:"id" json:"id"`
	ProdiID     string    `db:"prodi_id" json:"prodi_id"`
	Code        string    `db:"code" json:"code"`
	Name        string    `db:"name" json:"name"`
	SKS         int       `db:"sks" json:"sks"`
	Type        string    `db:"type" json:"type"` // wajib, pilihan
	Description *string   `db:"description" json:"description,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// ClassSchedule represents a jadwal kelas
type ClassSchedule struct {
	ID             string    `db:"id" json:"id"`
	CourseID       string    `db:"course_id" json:"course_id"`
	LecturerID     string    `db:"lecturer_id" json:"lecturer_id"`
	AcademicYearID string    `db:"academic_year_id" json:"academic_year_id"`
	Room           string    `db:"room" json:"room"`
	Day            string    `db:"day" json:"day"`
	StartTime      string    `db:"start_time" json:"start_time"`
	EndTime        string    `db:"end_time" json:"end_time"`
	Quota          int       `db:"quota" json:"quota"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}

// KRS represents Kartu Rencana Studi
type KRS struct {
	ID             string    `db:"id" json:"id"`
	StudentID      string    `db:"student_id" json:"student_id"`
	AcademicYearID string    `db:"academic_year_id" json:"academic_year_id"`
	Status         string    `db:"status" json:"status"` // draft, submitted, approved, rejected
	TotalSKS       int       `db:"total_sks" json:"total_sks"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

// KRSDetail represents one class in a KRS
type KRSDetail struct {
	ID      string `db:"id" json:"id"`
	KRSID   string `db:"krs_id" json:"krs_id"`
	ClassID string `db:"class_id" json:"class_id"`
}

// Grade represents nilai mahasiswa
type Grade struct {
	ID        string    `db:"id" json:"id"`
	StudentID string    `db:"student_id" json:"student_id"`
	ClassID   string    `db:"class_id" json:"class_id"`
	Midterm   float64   `db:"midterm" json:"midterm"`
	Final     float64   `db:"final" json:"final"`
	Grade     string    `db:"grade" json:"grade"` // A, B, C, D, E
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
