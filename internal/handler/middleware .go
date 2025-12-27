package handler

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
)

func (h *Handler) AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header required")
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization format")
			}

			tokenString := parts[1]

			claims, err := h.service.ValidateToken(tokenString)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
			}

			userIDStr, ok := claims["sub"].(string)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims")
			}

			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user ID in token")
			}

			c.Set("user_id", userID)

			return next(c)
		}
	}
}

func GetUserID(c echo.Context) (uuid.UUID, bool) {
	userID, ok := c.Get("user_id").(uuid.UUID)
	return userID, ok
}
