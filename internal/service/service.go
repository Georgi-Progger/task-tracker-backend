package service

import (
	"context"
	"time"

	"github.com/Georgi-Progger/task-tracker-backend/internal/domain/entity"
	"github.com/Georgi-Progger/task-tracker-backend/internal/repo"
	"github.com/Georgi-Progger/task-tracker-common/kafka"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, user entity.User) (string, error)
	Login(ctx context.Context, user entity.User, refreshTokenTTL time.Duration) (string, string, error)
	RefreshAccessToken(ctx context.Context, refreshTokenString string) (string, error)
	ValidateToken(tokenString string) (jwt.MapClaims, error)
}

type UserService interface {
	GetUserById(ctx context.Context, userId string) (entity.User, error)
}

type TaskService interface {
	CreateTask(ctx context.Context, task entity.Task) (string, error)
	GetUserTasks(ctx context.Context, userId uuid.UUID, limit, offset int) ([]entity.Task, error)
	UpdateTask(ctx context.Context, task entity.Task) error
	DeleteTask(ctx context.Context, taskId, userId uuid.UUID) error
}

type EmailService struct {
}

type Service struct {
	AuthService
	UserService
	TaskService
}

func NewService(repo repo.Repository, jwtSecret string, poducer kafka.Producer, accessTokenTTL time.Duration) Service {
	return Service{
		AuthService: NewAuthService(repo.UserReposetory, repo.RefreshTokenRepository, NewEmailService(poducer), jwtSecret, accessTokenTTL),
		UserService: NewUserService(repo.UserReposetory),
		TaskService: NewTaskSrvice(repo.TaskRepository),
	}
}
