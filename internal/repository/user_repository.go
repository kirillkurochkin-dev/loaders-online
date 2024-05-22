package repository

import (
	"context"
	"database/sql"
	"loaders-online/internal/entity/dto"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return UserRepository{db: db}
}

func (u *UserRepository) Register(ctx context.Context, userDto *dto.UserSignUpDto) (int, error) {
	var userID int
	err := u.db.QueryRowContext(ctx, "INSERT INTO users (username, password_hash, role) VALUES ($1, $2, $3) RETURNING user_id",
		userDto.Username, userDto.Password, userDto.Role).Scan(&userID)
	return userID, err
}

func (u *UserRepository) GetUserByUsername(ctx context.Context, username string) (dto.UserByUsername, error) {
	var us dto.UserByUsername
	err := u.db.QueryRowContext(ctx, "SELECT user_id, role, password_hash FROM users WHERE username=$1", username).
		Scan(&us.UserID, &us.Role, &us.PasswordHash)
	return us, err
}
