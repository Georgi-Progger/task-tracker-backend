package repo

import (
	"context"
	"time"

	"github.com/Georgi-Progger/task-tracker-backend/internal/domain/entity"
	"github.com/Georgi-Progger/task-tracker-backend/internal/domain/model"

	"github.com/google/uuid"
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

func (t *taskRepository) GetUserTasks(ctx context.Context, userId uuid.UUID, limit, offset int) ([]entity.Task, error) {
	query := `
		SELECT id, title, text, task_status FROM tasks
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3;
	`

	var tasks []entity.Task
	rows, err := t.db.QueryContext(ctx, query, userId, limit, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var task entity.Task
		err := rows.Scan(&task.Id, &task.Title, &task.Text, &task.Status)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}
func (t *taskRepository) CreateTask(ctx context.Context, title, text string, status entity.Status, userId uuid.UUID) (string, error) {
	query := `
		INSERT INTO tasks (id, title, text, user_id, created_at, task_status) VALUES ($1, $2, $3, $4, $5, $6);
	`

	taskId := uuid.New()
	_, err := t.db.ExecContext(ctx, query, taskId, title, text, userId, time.Now(), status)
	if err != nil {
		return "", err
	}

	return taskId.String(), nil
}

func (t *taskRepository) UpdateTask(ctx context.Context, taskId uuid.UUID, userId uuid.UUID, task model.TaskRequest) error {
	query := `
		UPDATE tasks  
		SET text = $1, title = $2, task_status = $3
		WHERE user_id = $4 AND id = $5;
	`

	_, err := t.db.ExecContext(ctx, query, task.Text, task.Title, task.Status, userId, taskId)
	return err
}

func (t *taskRepository) DeleteTask(ctx context.Context, taskId, userId uuid.UUID) error {
	query := `
		DELETE FROM tasks 
		WHERE id = $1 AND user_id = $2;
	`

	_, err := t.db.ExecContext(ctx, query, taskId, userId)
	return err
}

func (t *taskRepository) CountUsersTasks(ctx context.Context) ([]model.TaskCounter, error) {
	query := `
		SELECT
			u.email,
			COUNT(CASE WHEN t.task_status = 'COMPLETE' THEN 1 END) as completed_count,
			COUNT(CASE WHEN t.task_status IN ('CREATE', 'IN_PROGRESS') THEN 1 END) as pending_count
		FROM users u
		LEFT JOIN tasks t ON u.id = t.user_id
		GROUP BY u.id, u.email
		ORDER BY u.email;
	`

	var tasksCounters []model.TaskCounter
	rows, err := t.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var taskCounter model.TaskCounter
		err := rows.Scan(&taskCounter.Email, &taskCounter.CompleteTaskCount, &taskCounter.PendingTaskCount)
		if err != nil {
			return nil, err
		}

		tasksCounters = append(tasksCounters, taskCounter)
	}

	return tasksCounters, nil
}
