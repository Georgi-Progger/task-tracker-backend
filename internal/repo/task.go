package repo

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type taskRepository struct {
	db sqlx.DB
}

func NewTaskRepository(db sqlx.DB) *taskRepository {
	return &taskRepository{
		db: db,
	}
}

func (t *taskRepository) GetTasks(ctx context.Context) {

}
