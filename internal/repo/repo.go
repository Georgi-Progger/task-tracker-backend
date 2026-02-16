package repo

import (
	"context"
	"time"

	"github.com/Georgi-Progger/task-tracker-backend/internal/domain/entity"
	"github.com/Georgi-Progger/task-tracker-backend/internal/domain/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserReposetory interface {
	GetUserById(ctx context.Context, userId string) (entity.User, error)
	CreateUser(ctx context.Context, id uuid.UUID, name, email, hashPassword string) error
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
}

type TaskRepository interface {
	GetUserTasks(ctx context.Context, userId uuid.UUID, limit, offset int) ([]entity.Task, error)
	CreateTask(ctx context.Context, title, text string, status entity.Status, userId uuid.UUID) (string, error)
	UpdateTask(ctx context.Context, taskId uuid.UUID, userId uuid.UUID, task model.TaskRequest) error
	DeleteTask(ctx context.Context, taskId, userId uuid.UUID) error
	CountUsersTasks(ctx context.Context) ([]model.TaskCounter, error)
}

type RefreshTokenRepository interface {
	CreateRefreshToken(ctx context.Context, userID uuid.UUID, ttl time.Duration) (entity.RefreshToken, error)
	GetRefreshToken(ctx context.Context, tokenString string) (entity.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenString string) error
}

type Repository struct {
	UserReposetory
	TaskRepository
	RefreshTokenRepository
}

func NewRepository(db *sqlx.DB) Repository {
	return Repository{
		UserReposetory:         NewUserRepository(db),
		TaskRepository:         NewTaskRepository(db),
		RefreshTokenRepository: NewRefreshTokenRepository(db),
	}
}
