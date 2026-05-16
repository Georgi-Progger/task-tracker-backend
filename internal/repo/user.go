package repo

import (
	"context"

	"github.com/Georgi-Progger/task-tracker-backend/internal/domain/entity"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) CreateUser(ctx context.Context, id uuid.UUID, name, email, hashPassword string) error {
	query := `
			INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4);
	`

	_, err := u.db.ExecContext(ctx, query, id, name, email, hashPassword)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) GetUserById(ctx context.Context, userId string) (entity.User, error) {
	query := `
			SELECT id, email FROM users WHERE id = $1;
	`

	var user entity.User
	err := u.db.QueryRowContext(ctx, query, userId).Scan(
		&user.Id,
		&user.Email,
	)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	query := `
			SELECT id, email, password FROM users WHERE email = $1;
	`

	var user entity.User
	err := u.db.QueryRowContext(ctx, query, email).Scan(
		&user.Id,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}
