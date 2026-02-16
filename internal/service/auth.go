package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Georgi-Progger/task-tracker-backend/internal/domain"
	"github.com/Georgi-Progger/task-tracker-backend/internal/domain/entity"
	"github.com/Georgi-Progger/task-tracker-backend/internal/domain/model"
	"github.com/Georgi-Progger/task-tracker-backend/internal/repo"
	"github.com/Georgi-Progger/task-tracker-backend/pkg/hash"
	"github.com/Georgi-Progger/task-tracker-common/logger"
	"github.com/google/uuid"

	"github.com/golang-jwt/jwt/v5"
)

type authService struct {
	userRepo         repo.UserReposetory
	refreshTokenRepo repo.RefreshTokenRepository
	emailService     emailService
	userService      userService
	jwtSecret        []byte
	accessTokenTTL   time.Duration
	logger           logger.Logger
}

func NewAuthService(userRepo repo.UserReposetory, refreshTokenRepo repo.RefreshTokenRepository,
	emailService emailService, userService userService, jwtSecret string, accessTokenTTL time.Duration, logger logger.Logger) *authService {
	return &authService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		emailService:     emailService,
		jwtSecret:        []byte(jwtSecret),
		accessTokenTTL:   accessTokenTTL,
		userService:      userService,
		logger:           logger,
	}
}

func (a *authService) Register(ctx context.Context, user model.RegisterRequest, refreshTokenTTL time.Duration) (string, string, error) {
	_, err := a.userRepo.GetUserByEmail(ctx, user.Email)
	if err == nil {
		return "", "", domain.ErrEmailInUse
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return "", "", domain.ErrEmailInUse
	}

	hashPassword, err := hash.HashPassword(user.Password)
	if err != nil {
		return "", "", fmt.Errorf("error hashed password")
	}

	userEntity := entity.User{
		Id:       uuid.New(),
		Email:    user.Email,
		Name:     user.Name,
		Password: user.Password,
	}

	err = a.userService.CreateUser(ctx, userEntity.Id, userEntity.Name, userEntity.Email, hashPassword)
	if err != nil {
		return "", "", fmt.Errorf("error create user")
	}

	go func() {
		email := model.Email{
			Recipient: user.Email,
			Subject:   fmt.Sprintf("Приветствуем, %s, в taskcounter", user.Name),
			Body:      "Вы прошли успешную регистрацию в taskcounter!!!",
		}

		err = a.emailService.SendWelcomeMessage(email)
		if err != nil {
			a.logger.Error(err, "Failed to send welcome email")
		}
	}()

	accessToken, err := a.generateAccessToken(userEntity)
	if err != nil {
		return "", "", fmt.Errorf("error generate token")
	}

	token, err := a.refreshTokenRepo.CreateRefreshToken(ctx, userEntity.Id, refreshTokenTTL)
	if err != nil {
		return "", "", err
	}

	return accessToken, token.Token, nil
}

func (a *authService) Login(ctx context.Context, userDomain model.LoginRequest, refreshTokenTTL time.Duration) (string, string, error) {
	user, err := a.userRepo.GetUserByEmail(ctx, userDomain.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", domain.ErrInvalidCredentials
		}
		return "", "", err
	}

	err = hash.VerifyPassword(user.Password, userDomain.Password)
	if err != nil {
		return "", "", domain.ErrInvalidCredentials
	}

	accessToken, err := a.generateAccessToken(user)
	if err != nil {
		return "", "", err
	}

	token, err := a.refreshTokenRepo.CreateRefreshToken(ctx, user.Id, refreshTokenTTL)
	if err != nil {
		return "", "", err
	}

	return accessToken, token.Token, nil
}

func (a *authService) RefreshAccessToken(ctx context.Context, refreshTokenString string) (string, error) {
	token, err := a.refreshTokenRepo.GetRefreshToken(ctx, refreshTokenString)
	if err != nil {
		return "", domain.ErrInvalidToken
	}

	if token.Revoked {
		return "", domain.ErrInvalidToken
	}

	if time.Now().After(token.ExpiresAt) {
		return "", domain.ErrExpiredToken
	}

	user, err := a.userRepo.GetUserById(ctx, token.ID.String())
	if err != nil {
		return "", err
	}

	accessToken, err := a.generateAccessToken(user)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (a *authService) generateAccessToken(user entity.User) (string, error) {
	expirationTime := time.Now().Add(a.accessTokenTTL)

	claims := jwt.MapClaims{
		"sub":      user.Id.String(),
		"username": user.Name,
		"email":    user.Email,
		"exp":      expirationTime.Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(a.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *authService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrInvalidToken
		}
		return a.jwtSecret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, domain.ErrExpiredToken
		}
		return nil, domain.ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, domain.ErrInvalidToken
}
