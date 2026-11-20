package auth

import (
	"time"

	"github.com/google/uuid"
)

// User represents the user entity
type User struct {
	ID           uuid.UUID  `db:"id"`
	NIM          *string    `db:"nim"`
	NIP          *string    `db:"nip"`
	Name         string     `db:"name"`
	Email        string     `db:"email"`
	PasswordHash string     `db:"password_hash"`
	Role         string     `db:"role"` // admin, lecturer, student, staff
	IsActive     bool       `db:"is_active"`
	AvatarURL    *string    `db:"avatar_url"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
}
