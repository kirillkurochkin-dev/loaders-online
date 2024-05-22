package dto

type UserSignInDto struct {
	Username string `json:"username" validate:"required,max=50"`
	Password string `json:"password_hash" validate:"required"`
	Role     string `json:"role" validate:"required,oneof=loader customer"`
}
