package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/Georgi-Progger/task-tracker-backend/internal/domain"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Register(c echo.Context) error {
	var user domain.User

	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if user.Email == "" || user.Name == "" || user.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Email, name and password are required")
	}

	accessToken, err := h.service.AuthService.Register(c.Request().Context(), user)
	if err != nil {
		if errors.Is(err, domain.ErrEmailInUse) {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusConflict, "Email already in use")
		}
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Error creating user")
	}

	response := domain.Response{
		AccessToken: accessToken,
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *Handler) Login(c echo.Context) error {
	var user domain.User

	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if user.Email == "" || user.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Email and password are required")
	}

	accessToken, refreshToken, err := h.service.AuthService.Login(c.Request().Context(), user, 15*time.Minute)
	if err != nil {
		c.Logger().Error(err)
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return echo.NewHTTPError(http.StatusNotFound, "Email or password incorrect")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Error login user")
	}

	response := domain.Response{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *Handler) RefreshToken(c echo.Context) error {
	var req domain.RefreshRequest

	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	token, err := h.service.RefreshAccessToken(c.Request().Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidToken) || errors.Is(err, domain.ErrExpiredToken) {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired refresh token")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")

	}

	response := domain.RefreshResponse{Token: token}
	return c.JSON(http.StatusCreated, response)
}
