package dto

type UserSignUpDto struct {
	Username string `json:"username" validate:"required,max=50"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"required,oneof=loader customer"`
}

type UserSignInDto struct {
	Username string `json:"username" validate:"required,max=50"`
	Password string `json:"password" validate:"required"`
}

type UserByUsername struct {
	UserID       int    `json:"user_id" validate:"required"`
	PasswordHash string `json:"password" validate:"required"`
	Role         string `json:"role" validate:"required,oneof=loader customer"`
}
