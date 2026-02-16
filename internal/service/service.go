package service

import (
	"context"
	"time"

	"github.com/Georgi-Progger/task-tracker-backend/internal/domain/entity"
	"github.com/Georgi-Progger/task-tracker-backend/internal/domain/model"
	"github.com/Georgi-Progger/task-tracker-backend/internal/repo"
	"github.com/Georgi-Progger/task-tracker-common/kafka"
	"github.com/Georgi-Progger/task-tracker-common/logger"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, user model.RegisterRequest, refreshTokenTTL time.Duration) (string, string, error)
	Login(ctx context.Context, user model.LoginRequest, refreshTokenTTL time.Duration) (string, string, error)
	RefreshAccessToken(ctx context.Context, refreshTokenString string) (string, error)
	ValidateToken(tokenString string) (jwt.MapClaims, error)
}

type UserService interface {
	CreateUser(ctx context.Context, id uuid.UUID, name, email, password string) error
	GetUserById(ctx context.Context, userId string) (entity.User, error)
}

type TaskService interface {
	CreateTask(ctx context.Context, userID uuid.UUID, task model.TaskRequest) (string, error)
	GetUserTasks(ctx context.Context, userId uuid.UUID, limit, offset int) ([]entity.Task, error)
	UpdateTask(ctx context.Context, taskId uuid.UUID, userId uuid.UUID, task model.TaskRequest) error
	DeleteTask(ctx context.Context, taskId, userId uuid.UUID) error
}

type EmailService interface {
	SendTaskCountMessage(ctx context.Context) error
}

type Service struct {
	AuthService
	UserService
	TaskService
	EmailService
}

func NewService(repo repo.Repository, jwtSecret string, poducer kafka.Producer, accessTokenTTL time.Duration, logger logger.Logger) Service {
	return Service{
		AuthService:  NewAuthService(repo.UserReposetory, repo.RefreshTokenRepository, *NewEmailService(repo.TaskRepository, poducer, logger), *NewUserService(repo.UserReposetory), jwtSecret, accessTokenTTL, logger),
		UserService:  NewUserService(repo.UserReposetory),
		TaskService:  NewTaskSrvice(repo.TaskRepository),
		EmailService: NewEmailService(repo.TaskRepository, poducer, logger),
	}
}
