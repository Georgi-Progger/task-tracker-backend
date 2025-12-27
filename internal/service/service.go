package service

import (
	"context"
	"time"

	"github.com/Georgi-Progger/task-tracker-backend/internal/domain"
	"github.com/Georgi-Progger/task-tracker-backend/internal/repo"
	logger "github.com/Georgi-Progger/task-tracker-backend/pkg/looger"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	Register(ctx context.Context, user domain.User) (string, error)
	Login(ctx context.Context, user domain.User, refreshTokenTTL time.Duration) (string, string, error)
	RefreshAccessToken(ctx context.Context, refreshTokenString string) (string, error)
	ValidateToken(tokenString string) (jwt.MapClaims, error)
}

type UserService interface {
	GetUserById(ctx context.Context, userId string) (domain.User, error)
}

type Service struct {
	AuthService
	UserService
}

func NewService(repo repo.Repository, jwtSecret string, accessTokenTTL time.Duration, logger logger.Logger) Service {
	return Service{
		AuthService: NewAuthService(repo.UserReposetory, repo.RefreshTokenRepository, jwtSecret, accessTokenTTL, logger),
		UserService: NewUserService(repo.UserReposetory),
	}
}
