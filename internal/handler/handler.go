package handler

import (
	"github.com/Georgi-Progger/task-tracker-backend/internal/service"
	"github.com/Georgi-Progger/task-tracker-backend/pkg/logger"
)

type Handler struct {
	service service.Service
	logger  logger.Logger
}

func NewHandler(service service.Service, logger logger.Logger) Handler {
	return Handler{
		service: service,
		logger:  logger,
	}
}
