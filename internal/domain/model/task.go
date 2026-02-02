package model

type TaskCounter struct {
	CompleteTaskCount int    `json:"complete_task"`
	PendingTaskCount  int    `json:"pending_task"`
	Email             string `json:"email"`
}
