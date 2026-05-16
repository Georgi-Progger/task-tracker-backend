package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetUser(c echo.Context) error {
	userID := c.Get("user_id").(uuid.UUID)
	if len(userID) == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	user, err := h.service.GetUserById(c.Request().Context(), userID.String())
	if err != nil {
		h.logger.Error(err, "Error get user")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, user)
}
