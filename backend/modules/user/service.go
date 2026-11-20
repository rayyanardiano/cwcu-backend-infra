package user

import (
	"context"
	"errors"

	"cwcu/backend/config"
	"cwcu/backend/modules/auth"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	repo  *Repository
	redis *redis.Client
	cfg   *config.Config
}

func NewService(repo *Repository, redis *redis.Client, cfg *config.Config) *Service {
	return &Service{repo: repo, redis: redis, cfg: cfg}
}

func (s *Service) GetProfile(ctx context.Context, id string) (*User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) ListUsers(ctx context.Context) ([]*User, error) {
	return s.repo.FindAll(ctx)
}

func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	u := &User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hash,
		Role:         req.Role,
		NIM:          req.NIM,
		NIP:          req.NIP,
	}

	if err := s.repo.Create(ctx, u); err != nil {
		return nil, errors.New("email already exists or database error")
	}

	return u, nil
}

func (s *Service) UpdateProfile(ctx context.Context, id string, req UpdateProfileRequest) (*User, error) {
	u := &User{
		Name:      req.Name,
		AvatarURL: req.AvatarURL,
	}
	if err := s.repo.Update(ctx, id, u); err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, id)
}

func (s *Service) DeleteUser(ctx context.Context, id string) error {
	return s.repo.SoftDelete(ctx, id)
}
