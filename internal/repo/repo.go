package repo

import (
	"context"
	"time"

	"github.com/Georgi-Progger/task-tracker-backend/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserReposetory interface {
	GetUserById(ctx context.Context, userId string) (domain.User, error)
	CreateUser(ctx context.Context, name, email, hashPassword string) error
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
}

type TaskRepository interface {
	GetUserTasks(ctx context.Context, userId string) ([]domain.Task, error)
	CreateTask(ctx context.Context, task domain.Task) (domain.Task, error)
	UpdateTask(ctx context.Context, task domain.Task) (domain.Task, error)
}

type RefreshTokenRepository interface {
	CreateRefreshToken(ctx context.Context, userID uuid.UUID, ttl time.Duration) (domain.RefreshToken, error)
	GetRefreshToken(ctx context.Context, tokenString string) (domain.RefreshToken, error)
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
