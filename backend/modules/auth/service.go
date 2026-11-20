package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cwcu/backend/config"
	pkgjwt "cwcu/backend/pkg/jwt"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

// Service handles auth business logic
type Service struct {
	repo  *Repository
	redis *redis.Client
	cfg   *config.Config
}

// NewService creates a new auth service
func NewService(repo *Repository, redis *redis.Client, cfg *config.Config) *Service {
	return &Service{repo: repo, redis: redis, cfg: cfg}
}

// Login validates credentials and returns tokens
func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !user.IsActive {
		return nil, errors.New("account is inactive")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	userID := user.ID.String()

	accessToken, err := pkgjwt.GenerateAccessToken(userID, user.Role, user.Email, s.cfg)
	if err != nil {
		return nil, err
	}

	refreshToken, err := pkgjwt.GenerateRefreshToken(userID, s.cfg)
	if err != nil {
		return nil, err
	}

	// Store refresh token in Redis
	redisKey := fmt.Sprintf("refresh_token:%s", userID)
	s.redis.Set(ctx, redisKey, refreshToken, s.cfg.JWTRefreshTTL)

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: UserInfo{
			ID:        userID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			NIM:       user.NIM,
			NIP:       user.NIP,
			AvatarURL: user.AvatarURL,
		},
	}, nil
}

// RefreshToken generates new access token from valid refresh token
func (s *Service) RefreshToken(ctx context.Context, req RefreshRequest) (*LoginResponse, error) {
	claims, err := pkgjwt.ParseToken(req.RefreshToken, s.cfg)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	userID := claims.UserID
	redisKey := fmt.Sprintf("refresh_token:%s", userID)

	stored, err := s.redis.Get(ctx, redisKey).Result()
	if err != nil || stored != req.RefreshToken {
		return nil, errors.New("refresh token expired or invalid")
	}

	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	newAccessToken, err := pkgjwt.GenerateAccessToken(userID, user.Role, user.Email, s.cfg)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: req.RefreshToken,
		User: UserInfo{
			ID:    userID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}

// Logout invalidates the refresh token
func (s *Service) Logout(ctx context.Context, userID string, accessToken string) error {
	// Remove refresh token
	redisKey := fmt.Sprintf("refresh_token:%s", userID)
	s.redis.Del(ctx, redisKey)

	// Blacklist access token until it expires
	blacklistKey := fmt.Sprintf("blacklist:%s", accessToken)
	s.redis.Set(ctx, blacklistKey, "1", s.cfg.JWTAccessTTL)

	return nil
}

// HashPassword hashes a plain password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// IsTokenBlacklisted checks if a token has been invalidated
func (s *Service) IsTokenBlacklisted(ctx context.Context, token string) bool {
	key := fmt.Sprintf("blacklist:%s", token)
	result := s.redis.Exists(ctx, key)
	return result.Val() > 0
}

// generatePasswordResetToken creates a time-limited reset token (stored in Redis)
func (s *Service) generatePasswordResetToken(ctx context.Context, email string) (string, error) {
	token := fmt.Sprintf("%s-%d", email, time.Now().UnixNano())
	hashed, err := HashPassword(token)
	if err != nil {
		return "", err
	}
	key := fmt.Sprintf("reset_password:%s", hashed)
	s.redis.Set(ctx, key, email, 30*time.Minute)
	return hashed, nil
}
