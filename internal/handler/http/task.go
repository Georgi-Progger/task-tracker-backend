package http

import (
	"net/http"
	"strconv"

	"github.com/Georgi-Progger/task-tracker-backend/internal/domain/model"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateTask(c echo.Context) error {
	userID := c.Get("user_id").(uuid.UUID)
	if len(userID) == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	var task model.TaskRequest
	if err := c.Bind(&task); err != nil {
		h.logger.Error(err, "Error body task")
		return echo.NewHTTPError(http.StatusInternalServerError, "Error body task")
	}

	if err := h.validator.Struct(task); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	taskId, err := h.service.TaskService.CreateTask(c.Request().Context(), userID, task)
	if err != nil {
		h.logger.Error(err, "Error create task")
		return echo.NewHTTPError(http.StatusInternalServerError, "Error create task")
	}
	return c.JSON(http.StatusCreated, taskId)
}

func (h *Handler) GetUserTasks(c echo.Context) error {
	userID := c.Get("user_id").(uuid.UUID)
	if len(userID) == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	limit := c.QueryParam("limit")
	offset := c.QueryParam("offset")

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt <= 0 {
		limitInt = 0
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil || offsetInt < 0 {
		offsetInt = 0
	}

	tasks, err := h.service.TaskService.GetUserTasks(c.Request().Context(), userID, limitInt, offsetInt)
	if err != nil {
		h.logger.Error(err, "Error get tasks")
		return echo.NewHTTPError(http.StatusInternalServerError, "Error get tasks")
	}
	return c.JSON(http.StatusCreated, tasks)
}

func (h *Handler) UpdateTask(c echo.Context) error {
	userID := c.Get("user_id").(uuid.UUID)
	if len(userID) == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	taskId := c.Param("taskId")

	var task model.TaskRequest
	if err := c.Bind(&task); err != nil {
		h.logger.Error(err, "Failed get task body")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	taskUUID, err := uuid.Parse(taskId)
	if err != nil {
		h.logger.Error(err, "Error parse task id")
		return echo.NewHTTPError(http.StatusInternalServerError, "Error get task id")
	}

	err = h.service.TaskService.UpdateTask(c.Request().Context(), taskUUID, userID, task)
	if err != nil {
		h.logger.Error(err, "Error update task")
		return echo.NewHTTPError(http.StatusInternalServerError, "Error update task")
	}
	return c.JSON(http.StatusCreated, Response{
		Message: "task is update",
	})
}

func (h *Handler) DeleteTask(c echo.Context) error {
	userID := c.Get("user_id").(uuid.UUID)
	if len(userID) == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	taskId := c.Param("taskId")
	taskUUID, err := uuid.Parse(taskId)
	if err != nil {
		h.logger.Error(err, "Error get task id")
		return echo.NewHTTPError(http.StatusInternalServerError, "Error get task id")
	}

	err = h.service.TaskService.DeleteTask(c.Request().Context(), taskUUID, userID)
	if err != nil {
		h.logger.Error(err, "Error delete task")
		return echo.NewHTTPError(http.StatusInternalServerError, "Error delete task")
	}
	return c.JSON(http.StatusCreated, Response{
		Message: "task is delete",
	})
}
