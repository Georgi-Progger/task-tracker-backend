package repo

import (
	"context"

	"github.com/Georgi-Progger/task-tracker-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type taskRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) *taskRepository {
	return &taskRepository{
		db: db,
	}
}

func (t *taskRepository) GetUserTasks(ctx context.Context, userId string) ([]domain.Task, error) {
	// query := `
	// `

	return nil, nil
}
func (t *taskRepository) CreateTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	// query := `

	// `
	return domain.Task{}, nil
}
func (t *taskRepository) UpdateTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	return domain.Task{}, nil
}
