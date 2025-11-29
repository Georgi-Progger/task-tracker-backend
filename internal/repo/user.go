package repo

import "github.com/jmoiron/sqlx"

type userRepository struct {
	db sqlx.DB
}

func NewUserRepository(db sqlx.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}
