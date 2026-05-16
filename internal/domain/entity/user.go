package entity

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name,omitempty"`
	Email    string    `json:"email"`
	Password string    `json:"password,omitempty"`
}
