package model

import "github.com/Georgi-Progger/task-tracker-backend/internal/domain/entity"

type TaskRequest struct {
	Title  string        `json:"title"`
	Text   string        `json:"text"`
	Status entity.Status `json:"status"`
}

type TaskCounter struct {
	CompleteTaskCount int    `json:"complete_task"`
	PendingTaskCount  int    `json:"pending_task"`
	Email             string `json:"email"`
}
