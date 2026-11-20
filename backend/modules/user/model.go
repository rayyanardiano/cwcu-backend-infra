package user

import "time"

type User struct {
	ID           string     `db:"id" json:"id"`
	NIM          *string    `db:"nim" json:"nim,omitempty"`
	NIP          *string    `db:"nip" json:"nip,omitempty"`
	Name         string     `db:"name" json:"name"`
	Email        string     `db:"email" json:"email"`
	PasswordHash string     `db:"password_hash" json:"-"`
	Role         string     `db:"role" json:"role"`
	IsActive     bool       `db:"is_active" json:"is_active"`
	AvatarURL    *string    `db:"avatar_url" json:"avatar_url,omitempty"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at" json:"-"`
}
