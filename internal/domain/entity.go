package domain

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	Id     uuid.UUID
	UserId string
	Title  string
	Text   string
}

type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
	Revoked   bool
}

type User struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"user_name,omitempty"`
	Email    string    `json:"email"`
	Password string    `json:"password,omitempty"`
}

type Response struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type ErrorResponse struct {
	Message string
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	Token string `json:"access_token"`
}
