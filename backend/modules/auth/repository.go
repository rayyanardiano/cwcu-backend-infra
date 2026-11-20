package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles auth data access
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new auth repository
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// FindByEmail retrieves a user by email
func (r *Repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	q := `SELECT id, nim, nip, name, email, password_hash, role, is_active, avatar_url, created_at, updated_at
		  FROM users WHERE email = $1 AND deleted_at IS NULL`

	user := &User{}
	err := r.db.QueryRow(ctx, q, email).Scan(
		&user.ID, &user.NIM, &user.NIP, &user.Name, &user.Email,
		&user.PasswordHash, &user.Role, &user.IsActive, &user.AvatarURL,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// FindByID retrieves a user by ID
func (r *Repository) FindByID(ctx context.Context, id string) (*User, error) {
	q := `SELECT id, nim, nip, name, email, password_hash, role, is_active, avatar_url, created_at, updated_at
		  FROM users WHERE id = $1 AND deleted_at IS NULL`

	user := &User{}
	err := r.db.QueryRow(ctx, q, id).Scan(
		&user.ID, &user.NIM, &user.NIP, &user.Name, &user.Email,
		&user.PasswordHash, &user.Role, &user.IsActive, &user.AvatarURL,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
