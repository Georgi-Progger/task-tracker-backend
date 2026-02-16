package http

import (
	"fmt"
	"net/http"
	"strings"
	"time"

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

func (h *Handler) RateLimitMiddleware(limit int, window time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL
			key := fmt.Sprintf("rate_limit:%s:%s", c.RealIP(), path)

			allowed, err := h.limiter.Allow(c.Request().Context(), limit, window, key)
			if err != nil {
				h.logger.Error(err, "error get allow status")
				return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
			}

			if !allowed {
				c.Response().Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
				c.Response().Header().Set("X-RateLimit-Remaining", "0")
				c.Response().Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(window).Unix()))

				return echo.NewHTTPError(http.StatusTooManyRequests, "Too many requests")
			}
			return next(c)
		}
	}
}
