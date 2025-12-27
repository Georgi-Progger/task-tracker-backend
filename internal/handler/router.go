package handler

import "github.com/labstack/echo/v4"

func (h *Handler) SetupRoutes(e *echo.Echo) {
	api := e.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/sign-up", h.Register)
			auth.POST("/sign-in", h.Login)
		}

		protected := api.Group("")
		protected.Use(h.AuthMiddleware())
		{
			protected.GET("/user", h.GetUser)
		}
	}

}
