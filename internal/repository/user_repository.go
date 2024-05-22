package repository

import (
	"context"
	"database/sql"
	"loaders-online/internal/entity/dto"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) Register(ctx context.Context, userDto *dto.UserSignInDto) (int, error) {
	var userID int
	err := u.db.QueryRowContext(ctx, "INSERT INTO users (username, password_hash, role) VALUES ($1, $2, $3) RETURNING user_id",
		userDto.Username, userDto.Password, userDto.Role).Scan(&userID)
	return userID, err
}
