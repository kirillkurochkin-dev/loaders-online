package service

import (
	"context"
	"loaders-online/internal/entity/dto"
	"loaders-online/pkg/util"
)

type UserRepository interface {
	Register(ctx context.Context, user *dto.UserSignInDto) (int, error)
}

type UserService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) Register(ctx context.Context, user *dto.UserSignInDto) error {
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
