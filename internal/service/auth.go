package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Georgi-Progger/task-tracker-backend/internal/domain"
	"github.com/Georgi-Progger/task-tracker-backend/internal/repo"
	"github.com/Georgi-Progger/task-tracker-backend/pkg/hash"
	logger "github.com/Georgi-Progger/task-tracker-backend/pkg/looger"
	"github.com/golang-jwt/jwt/v5"
)

type authService struct {
	userRepo         repo.UserReposetory
	refreshTokenRepo repo.RefreshTokenRepository
	log              logger.Logger
	jwtSecret        []byte
	accessTokenTTL   time.Duration
}

func NewAuthService(userRepo repo.UserReposetory, refreshTokenRepo repo.RefreshTokenRepository, jwtSecret string, accessTokenTTL time.Duration, log logger.Logger) *authService {
	return &authService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		log:              log,
		jwtSecret:        []byte(jwtSecret),
		accessTokenTTL:   accessTokenTTL,
	}
}

func (a *authService) Register(ctx context.Context, user domain.User) (string, error) {
	_, err := a.userRepo.GetUserByEmail(ctx, user.Email)
	if err == nil {
		return "", domain.ErrEmailInUse
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("user is register")
	}

	hashPassword, err := hash.HashPassword(user.Password)
	if err != nil {
		return "", fmt.Errorf("error hashed password")
	}

	err = a.userRepo.CreateUser(ctx, user.Name, user.Email, hashPassword)
	if err != nil {
		return "", fmt.Errorf("error create user")
	}

	accessToken, err := a.generateAccessToken(user)
	if err != nil {
		return "", fmt.Errorf("error generate token")
	}

	return accessToken, nil
}

func (a *authService) Login(ctx context.Context, userDomain domain.User, refreshTokenTTL time.Duration) (string, string, error) {
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

func (a *authService) generateAccessToken(user domain.User) (string, error) {
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
