package repo

import (
	"context"
	"time"

	"github.com/Georgi-Progger/task-tracker-backend/internal/domain/entity"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type refreshTokenRepository struct {
	db *sqlx.DB
}

func NewRefreshTokenRepository(db *sqlx.DB) *refreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) CreateRefreshToken(ctx context.Context, userID uuid.UUID, ttl time.Duration) (entity.RefreshToken, error) {
	tokenID := uuid.New()
	expiresAt := time.Now().Add(ttl)

	token := entity.RefreshToken{
		ID:        tokenID,
		UserID:    userID,
		Token:     tokenID.String(),
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
		Revoked:   false,
	}
	query := `
        INSERT INTO refresh_tokens (id, user_id, token, expires_at, created_at, revoked)
        VALUES ($1, $2, $3, $4, $5, $6);
    `
	_, err := r.db.ExecContext(ctx, query, &token.ID, &token.UserID, &token.Token, &token.ExpiresAt, &token.CreatedAt, &token.Revoked)
	if err != nil {
		return entity.RefreshToken{}, err
	}
	return token, nil
}

func (r *refreshTokenRepository) GetRefreshToken(ctx context.Context, tokenString string) (entity.RefreshToken, error) {
	query := `
        SELECT id, user_id, token, expires_at, created_at, revoked
        FROM refresh_tokens
        WHERE token = $1;
    `
	var token entity.RefreshToken
	err := r.db.QueryRowContext(ctx, query, tokenString).Scan(
		&token.ID,
		&token.UserID,
		&token.Token,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.Revoked,
	)
	if err != nil {
		return entity.RefreshToken{}, err
	}
	return token, nil
}

func (r *refreshTokenRepository) RevokeRefreshToken(ctx context.Context, tokenString string) error {
	query := `
        UPDATE refresh_tokens
        SET revoked = true
        WHERE token = $1;
    `
	_, err := r.db.ExecContext(ctx, query, tokenString)
	return err
}
