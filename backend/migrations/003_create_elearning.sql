-- =============================================
-- Migration 003: E-Learning Tables
-- =============================================

-- E-Learning Courses (linked to class schedules)
CREATE TABLE IF NOT EXISTS courses_elearn (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    class_id    UUID REFERENCES class_schedules(id),
    title       VARCHAR(255) NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (class_id)
);

-- Materials
CREATE TABLE IF NOT EXISTS materials (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_elearn_id UUID REFERENCES courses_elearn(id) ON DELETE CASCADE,
    title            VARCHAR(255) NOT NULL,
    type             VARCHAR(20) NOT NULL CHECK (type IN ('pdf', 'video', 'link', 'text')),
    content_url      TEXT,
    content          TEXT,
    "order"          INT DEFAULT 0,
    created_at       TIMESTAMPTZ DEFAULT NOW()
);

-- Assignments
CREATE TABLE IF NOT EXISTS assignments (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_elearn_id UUID REFERENCES courses_elearn(id) ON DELETE CASCADE,
    title            VARCHAR(255) NOT NULL,
    description      TEXT NOT NULL,
    due_date         TIMESTAMPTZ,
    max_score        NUMERIC(5,2) NOT NULL DEFAULT 100,
    created_at       TIMESTAMPTZ DEFAULT NOW()
);

-- Submissions
CREATE TABLE IF NOT EXISTS submissions (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    assignment_id UUID REFERENCES assignments(id) ON DELETE CASCADE,
    student_id    UUID REFERENCES users(id),
    file_url      TEXT,
    content       TEXT,
    score         NUMERIC(5,2),
    feedback      TEXT,
    submitted_at  TIMESTAMPTZ DEFAULT NOW(),
    graded_at     TIMESTAMPTZ,
    UNIQUE (assignment_id, student_id)
);

-- Quizzes
CREATE TABLE IF NOT EXISTS quizzes (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_elearn_id UUID REFERENCES courses_elearn(id) ON DELETE CASCADE,
    title            VARCHAR(255) NOT NULL,
    duration         INT NOT NULL DEFAULT 60, -- in minutes
    start_at         TIMESTAMPTZ,
    end_at           TIMESTAMPTZ,
    created_at       TIMESTAMPTZ DEFAULT NOW()
);

-- Quiz Questions
CREATE TABLE IF NOT EXISTS quiz_questions (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    quiz_id        UUID REFERENCES quizzes(id) ON DELETE CASCADE,
    question       TEXT NOT NULL,
    options        JSONB NOT NULL, -- [{label: "A", text: "..."}, ...]
    correct_answer VARCHAR(5) NOT NULL,
    points         INT NOT NULL DEFAULT 1
);

-- Quiz Attempts
CREATE TABLE IF NOT EXISTS quiz_attempts (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    quiz_id      UUID REFERENCES quizzes(id),
    student_id   UUID REFERENCES users(id),
    score        NUMERIC(5,2) DEFAULT 0,
    submitted_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (quiz_id, student_id)
);

-- Forums
CREATE TABLE IF NOT EXISTS forums (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_elearn_id UUID REFERENCES courses_elearn(id) ON DELETE CASCADE,
    title            VARCHAR(255) NOT NULL,
    created_by       UUID REFERENCES users(id),
    created_at       TIMESTAMPTZ DEFAULT NOW()
);

-- Forum Posts
CREATE TABLE IF NOT EXISTS forum_posts (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    forum_id   UUID REFERENCES forums(id) ON DELETE CASCADE,
    user_id    UUID REFERENCES users(id),
    content    TEXT NOT NULL,
    parent_id  UUID REFERENCES forum_posts(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Announcements
CREATE TABLE IF NOT EXISTS announcements (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_elearn_id UUID REFERENCES courses_elearn(id) ON DELETE CASCADE,
    title            VARCHAR(255) NOT NULL,
    content          TEXT NOT NULL,
    created_by       UUID REFERENCES users(id),
    created_at       TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_materials_course ON materials(course_elearn_id);
CREATE INDEX IF NOT EXISTS idx_assignments_course ON assignments(course_elearn_id);
CREATE INDEX IF NOT EXISTS idx_submissions_assignment ON submissions(assignment_id);
CREATE INDEX IF NOT EXISTS idx_forum_posts_forum ON forum_posts(forum_id);
CREATE INDEX IF NOT EXISTS idx_announcements_course ON announcements(course_elearn_id);
