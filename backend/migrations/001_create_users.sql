-- =============================================
-- Migration 001: Core Users & Auth
-- =============================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Users
CREATE TABLE IF NOT EXISTS users (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nim           VARCHAR(20)  UNIQUE,
    nip           VARCHAR(20)  UNIQUE,
    name          VARCHAR(255) NOT NULL,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT         NOT NULL,
    role          VARCHAR(20)  NOT NULL CHECK (role IN ('admin', 'lecturer', 'student', 'staff')),
    is_active     BOOLEAN      DEFAULT TRUE,
    avatar_url    TEXT,
    created_at    TIMESTAMPTZ  DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);

-- Faculties
CREATE TABLE IF NOT EXISTS faculties (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       VARCHAR(255) NOT NULL,
    code       VARCHAR(20)  UNIQUE NOT NULL,
    created_at TIMESTAMPTZ  DEFAULT NOW()
);

-- Departments
CREATE TABLE IF NOT EXISTS departments (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    faculty_id  UUID         REFERENCES faculties(id),
    name        VARCHAR(255) NOT NULL,
    code        VARCHAR(20)  UNIQUE NOT NULL,
    created_at  TIMESTAMPTZ  DEFAULT NOW()
);

-- Study Programs (Prodi)
CREATE TABLE IF NOT EXISTS study_programs (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    department_id UUID         REFERENCES departments(id),
    name          VARCHAR(255) NOT NULL,
    code          VARCHAR(20)  UNIQUE NOT NULL,
    degree        VARCHAR(10)  NOT NULL CHECK (degree IN ('D3', 'S1', 'S2', 'S3')),
    created_at    TIMESTAMPTZ  DEFAULT NOW()
);

-- Lecturers
CREATE TABLE IF NOT EXISTS lecturers (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id       UUID REFERENCES users(id) ON DELETE CASCADE,
    nidn          VARCHAR(20) UNIQUE,
    department_id UUID REFERENCES departments(id),
    created_at    TIMESTAMPTZ DEFAULT NOW()
);

-- Seed default admin user (password: Admin@123)
INSERT INTO users (name, email, password_hash, role)
VALUES ('Administrator', 'admin@cwcu.ac.id', crypt('Admin@123', gen_salt('bf')), 'admin')
ON CONFLICT (email) DO NOTHING;
