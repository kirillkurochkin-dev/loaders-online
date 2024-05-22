package service

import (
	"context"
	"loaders-online/internal/entity/dto"
	"loaders-online/pkg/util"
)

type UserRepository interface {
	Register(ctx context.Context, user *dto.UserSignUpDto) (int, error)
	GetUserByUsername(ctx context.Context, username string) (dto.UserByUsername, error)
}

type UserService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) Register(ctx context.Context, user *dto.UserSignUpDto) error {
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	_, err = s.userRepository.Register(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) Login(ctx context.Context, user *dto.UserSignInDto) (string, error) {
	tempUser, err := s.userRepository.GetUserByUsername(ctx, user.Username)
	if err != nil {
		return "", err
	}

	err = util.CheckPasswordHash(tempUser.PasswordHash, user.Password)
	if err != nil {
		return "", err
	}

	token, err := util.GenerateJWT(tempUser.UserID, tempUser.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}
