package http

import (
	"time"

	"github.com/labstack/echo/v4"
)

func (h *Handler) SetupRoutes(e *echo.Echo) {
	api := e.Group("/api")
	{
		auth := api.Group("/users")
		protectedUser := api.Group("/users")
		protectedTask := api.Group("/tasks")
		{
			auth.POST("/sign-up", h.Register)
			auth.POST("/sign-in", h.Login)
		}

		protectedUser.Use(h.AuthMiddleware())
		protectedUser.Use(h.RateLimitMiddleware(10, time.Minute))
		{
			protectedUser.GET("", h.GetUser)
		}

		protectedTask.Use(h.AuthMiddleware())
		protectedTask.Use(h.RateLimitMiddleware(10, time.Minute))
		{
			protectedTask.POST("", h.CreateTask)
			protectedTask.GET("", h.GetUserTasks)
			protectedTask.PUT("/:taskId", h.UpdateTask)
			protectedTask.DELETE("/:taskId", h.DeleteTask)
		}
	}
}
