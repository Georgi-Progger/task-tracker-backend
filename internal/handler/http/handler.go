package http

import (
	"github.com/Georgi-Progger/task-tracker-backend/internal/service"
	"github.com/Georgi-Progger/task-tracker-common/logger"
	"github.com/Georgi-Progger/task-tracker-rate-limiter/limiter"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Message string `json:"message"`
}

type Handler struct {
	limiter   limiter.Limiter
	service   service.Service
	logger    logger.Logger
	validator validator.Validate
}

func NewHandler(service service.Service, limiter limiter.Limiter, logger logger.Logger, validator validator.Validate) Handler {
	return Handler{
		limiter:   limiter,
		service:   service,
		logger:    logger,
		validator: validator,
	}
}
