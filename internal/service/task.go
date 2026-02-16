package service

import (
	"context"
	"fmt"

	"github.com/Georgi-Progger/task-tracker-backend/internal/domain/entity"
	"github.com/Georgi-Progger/task-tracker-backend/internal/domain/model"
	"github.com/Georgi-Progger/task-tracker-backend/internal/repo"

	"github.com/google/uuid"
)

type taskService struct {
	taskRepo repo.TaskRepository
}

func NewTaskSrvice(taskRepo repo.TaskRepository) *taskService {
	return &taskService{
		taskRepo: taskRepo,
	}
}

func (t *taskService) CreateTask(ctx context.Context, userID uuid.UUID, task model.TaskRequest) (string, error) {
	taskId, err := t.taskRepo.CreateTask(ctx, task.Title, task.Text, task.Status, userID)
	if err != nil {
		return "", fmt.Errorf("error created task: %v", err)
	}

	return taskId, nil
}

func (t *taskService) GetUserTasks(ctx context.Context, userId uuid.UUID, limit, offset int) ([]entity.Task, error) {
	tasks, err := t.taskRepo.GetUserTasks(ctx, userId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error get users tasks: %v", err)
	}

	return tasks, nil
}

func (t *taskService) UpdateTask(ctx context.Context, taskId uuid.UUID, userId uuid.UUID, task model.TaskRequest) error {
	err := t.taskRepo.UpdateTask(ctx, taskId, userId, task)
	if err != nil {
		return fmt.Errorf("update task error: %v", err)
	}
	return nil
}

func (t *taskService) DeleteTask(ctx context.Context, taskId, userId uuid.UUID) error {
	err := t.taskRepo.DeleteTask(ctx, taskId, userId)
	if err != nil {
		return fmt.Errorf("delete task error: %v", err)
	}
	return nil
}
