package entity

import "github.com/google/uuid"

type Status string

const (
	CREATE     Status = "CREATE"
	INPROGRESS Status = "IN_PROGRESS"
	COMPLETE   Status = "COMPLETE"
)

type Task struct {
	Id     uuid.UUID `json:"id"`
	UserId uuid.UUID `json:"-"`
	Title  string    `json:"title"`
	Text   string    `json:"text"`
	Status Status    `json:"status"`
}
