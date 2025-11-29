package repo

import (
	"context"

	"github.com/Georgi-Progger/task-tracker-backend/internal/model"
	"github.com/jmoiron/sqlx"
)

type UserReposetory interface {
	GetUserById(ctx context.Context, userId string) (model.User, error)
	CreateUser(ctx context.Context, user model.User) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
}

type TaskRepository interface {
	GetUserTasks(ctx context.Context, userId string) ([]model.Task, error)
	CreateTask(ctx context.Context, task model.Task) (model.Task, error)
	UpdateTask(ctx context.Context, task model.Task) (model.Task, error)
}

type Repository struct {
	UserReposetory
	TaskRepository
}

func NewRepository(db sqlx.DB) Repository {
	return Repository{
		UserReposetory: NewUserRepository(db),
		TaskRepository: NewTaskRepository(db),
	}
}
