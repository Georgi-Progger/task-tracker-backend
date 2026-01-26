package handler

import (
	"github.com/Georgi-Progger/task-tracker-backend/internal/service"
	"github.com/Georgi-Progger/task-tracker-common/logger"
	"github.com/Georgi-Progger/task-tracker-rate-limiter/limiter"
)

type Handler struct {
	limiter limiter.Limiter
	service service.Service
	logger  logger.Logger
}

func NewHandler(service service.Service, limiter limiter.Limiter, logger logger.Logger) Handler {
	return Handler{
		limiter: limiter,
		service: service,
		logger:  logger,
	}
}
