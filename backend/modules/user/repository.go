package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

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
		return nil, err
	}
	return user, nil
}

func (r *Repository) FindAll(ctx context.Context) ([]*User, error) {
	q := `SELECT id, nim, nip, name, email, role, is_active, avatar_url, created_at, updated_at
		  FROM users WHERE deleted_at IS NULL ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		u := &User{}
		if err := rows.Scan(&u.ID, &u.NIM, &u.NIP, &u.Name, &u.Email,
			&u.Role, &u.IsActive, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *Repository) Create(ctx context.Context, u *User) error {
	q := `INSERT INTO users (name, email, password_hash, role, nim, nip, is_active)
		  VALUES ($1, $2, $3, $4, $5, $6, true) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, q, u.Name, u.Email, u.PasswordHash, u.Role, u.NIM, u.NIP).
		Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

func (r *Repository) Update(ctx context.Context, id string, u *User) error {
	q := `UPDATE users SET name = $1, avatar_url = $2, updated_at = NOW()
		  WHERE id = $3 AND deleted_at IS NULL`
	_, err := r.db.Exec(ctx, q, u.Name, u.AvatarURL, id)
	return err
}

func (r *Repository) SoftDelete(ctx context.Context, id string) error {
	q := `UPDATE users SET deleted_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, q, id)
	return err
}
