package entity

type User struct {
	UserID       int    `json:"user_id" validate:"required"`
	Username     string `json:"username" validate:"required,max=50"`
	PasswordHash string `json:"password_hash" validate:"required"`
	Role         string `json:"role" validate:"required,oneof=loader customer"`
}
