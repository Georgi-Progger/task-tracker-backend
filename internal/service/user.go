package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Georgi-Progger/task-tracker-backend/internal/domain"
	"github.com/Georgi-Progger/task-tracker-backend/internal/repo"
)

type userService struct {
	userRepo repo.UserReposetory
}

func NewUserService(userRepo repo.UserReposetory) *userService {
	return &userService{
		userRepo: userRepo,
	}
}

func (u *userService) GetUserById(ctx context.Context, userId string) (domain.User, error) {
	user, err := u.userRepo.GetUserById(ctx, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user id not found")
		}
		return domain.User{}, err
	}
	return user, nil
}
