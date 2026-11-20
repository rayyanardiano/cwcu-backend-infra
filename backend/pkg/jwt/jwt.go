package jwt

import (
	"errors"
	"time"

	"cwcu/backend/config"

	gojwt "github.com/golang-jwt/jwt/v5"
)

// Claims is the JWT custom claims structure
type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	gojwt.RegisteredClaims
}

// GenerateAccessToken creates a short-lived access token
func GenerateAccessToken(userID, role, email string, cfg *config.Config) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role,
		Email:  email,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(cfg.JWTAccessTTL)),
			IssuedAt:  gojwt.NewNumericDate(time.Now()),
			Issuer:    cfg.AppName,
		},
	}
	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

// GenerateRefreshToken creates a long-lived refresh token
func GenerateRefreshToken(userID string, cfg *config.Config) (string, error) {
	claims := &gojwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: gojwt.NewNumericDate(time.Now().Add(cfg.JWTRefreshTTL)),
		IssuedAt:  gojwt.NewNumericDate(time.Now()),
		Issuer:    cfg.AppName,
	}
	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

// ParseToken validates and parses a JWT token
func ParseToken(tokenStr string, cfg *config.Config) (*Claims, error) {
	token, err := gojwt.ParseWithClaims(tokenStr, &Claims{}, func(t *gojwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*gojwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
