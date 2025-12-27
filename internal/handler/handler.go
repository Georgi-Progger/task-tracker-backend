package handler

import "github.com/Georgi-Progger/task-tracker-backend/internal/service"

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) Handler {
	return Handler{service: service}
}
