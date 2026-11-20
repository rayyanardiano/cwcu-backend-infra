-- =============================================
-- Migration 002: Academic (SIAKAD) Tables
-- =============================================

-- Academic Years
CREATE TABLE IF NOT EXISTS academic_years (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    year       VARCHAR(9)  NOT NULL, -- e.g. "2024/2025"
    semester   VARCHAR(10) NOT NULL CHECK (semester IN ('Ganjil', 'Genap')),
    is_active  BOOLEAN     DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Rooms
CREATE TABLE IF NOT EXISTS rooms (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       VARCHAR(100) NOT NULL,
    building   VARCHAR(100),
    capacity   INT NOT NULL DEFAULT 30,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Courses (Mata Kuliah)
CREATE TABLE IF NOT EXISTS courses (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    prodi_id    UUID REFERENCES study_programs(id),
    code        VARCHAR(20)  UNIQUE NOT NULL,
    name        VARCHAR(255) NOT NULL,
    sks         INT NOT NULL CHECK (sks >= 1 AND sks <= 6),
    type        VARCHAR(20) NOT NULL CHECK (type IN ('wajib', 'pilihan')),
    description TEXT,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

-- Class Schedules (Jadwal Kuliah)
CREATE TABLE IF NOT EXISTS class_schedules (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_id        UUID REFERENCES courses(id),
    lecturer_id      UUID REFERENCES lecturers(id),
    academic_year_id UUID REFERENCES academic_years(id),
    room             VARCHAR(100),
    day              VARCHAR(10) NOT NULL CHECK (day IN ('Senin','Selasa','Rabu','Kamis','Jumat','Sabtu')),
    start_time       TIME NOT NULL,
    end_time         TIME NOT NULL,
    quota            INT NOT NULL DEFAULT 30,
    created_at       TIMESTAMPTZ DEFAULT NOW(),
    updated_at       TIMESTAMPTZ DEFAULT NOW()
);

-- KRS (Kartu Rencana Studi)
CREATE TABLE IF NOT EXISTS krs (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id       UUID REFERENCES users(id),
    academic_year_id UUID REFERENCES academic_years(id),
    status           VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft','submitted','approved','rejected')),
    total_sks        INT NOT NULL DEFAULT 0,
    created_at       TIMESTAMPTZ DEFAULT NOW(),
    updated_at       TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (student_id, academic_year_id)
);

-- KRS Details (mata kuliah yang diambil dalam KRS)
CREATE TABLE IF NOT EXISTS krs_details (
    id       UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    krs_id   UUID REFERENCES krs(id) ON DELETE CASCADE,
    class_id UUID REFERENCES class_schedules(id)
);

-- Grades (Nilai)
CREATE TABLE IF NOT EXISTS grades (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID REFERENCES users(id),
    class_id   UUID REFERENCES class_schedules(id),
    midterm    NUMERIC(5,2) DEFAULT 0,
    final      NUMERIC(5,2) DEFAULT 0,
    grade      VARCHAR(3),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (student_id, class_id)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_class_schedules_day ON class_schedules(day);
CREATE INDEX IF NOT EXISTS idx_krs_student ON krs(student_id);
CREATE INDEX IF NOT EXISTS idx_grades_student ON grades(student_id);
