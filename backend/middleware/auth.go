package middleware

import (
	"net/http"
	"strings"

	"cwcu/backend/config"
	"cwcu/backend/pkg/jwt"
	"cwcu/backend/pkg/response"

	"github.com/labstack/echo/v4"
)

// JWTAuth is the authentication middleware
func JWTAuth(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, response.Error("Authorization header is required", nil))
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, response.Error("Invalid authorization format", nil))
			}

			claims, err := jwt.ParseToken(parts[1], cfg)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, response.Error("Invalid or expired token", nil))
			}

			// Store claims in context
			c.Set("user_id", claims.UserID)
			c.Set("user_role", claims.Role)
			c.Set("user_email", claims.Email)

			return next(c)
		}
	}
}

// RequireRoles restricts access to specific roles
func RequireRoles(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole := c.Get("user_role").(string)
			for _, r := range roles {
				if r == userRole {
					return next(c)
				}
			}
			return c.JSON(http.StatusForbidden, response.Error("Access forbidden: insufficient permissions", nil))
		}
	}
}
