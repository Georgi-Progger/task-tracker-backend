package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Georgi-Progger/task-tracker-backend/internal/domain/entity"
	"github.com/Georgi-Progger/task-tracker-backend/internal/repo"
	"github.com/google/uuid"
)

type userService struct {
	userRepo repo.UserReposetory
}

func NewUserService(userRepo repo.UserReposetory) *userService {
	return &userService{
		userRepo: userRepo,
	}
}

func (u *userService) CreateUser(ctx context.Context, id uuid.UUID, name, email, password string) error {
	err := u.userRepo.CreateUser(ctx, id, name, email, password)
	if err != nil {
		return fmt.Errorf("error create user")
	}
	return nil
}

func (u *userService) GetUserById(ctx context.Context, userId string) (entity.User, error) {
	user, err := u.userRepo.GetUserById(ctx, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, fmt.Errorf("user id not found")
		}
		return entity.User{}, err
	}
	return user, nil
}
