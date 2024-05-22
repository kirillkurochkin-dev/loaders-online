package service

import (
	"context"
	"loaders-online/internal/entity/dto"
	"loaders-online/internal/repository"
	"loaders-online/pkg/game"
	"loaders-online/pkg/util"
	"math/rand"
)

type UserRepository interface {
	Register(ctx context.Context, userDto *dto.UserSignUpDto) (int, error)
	GetUserByUsername(ctx context.Context, username string) (dto.UserByUsername, error)
}

type CustomerRepository interface {
	CreateCustomer(ctx context.Context, customer *dto.CustomerSignUpDto) error
	GetCustomerById(ctx context.Context, id int) (dto.CustomerOutputDto, error)
	UpdateCustomer(ctx context.Context, outputDto dto.CustomerUpdateDto) error
}

type LoaderRepository interface {
	GetAssignedLoaders(ctx context.Context, id int) ([]dto.LoaderOutputDto, error)
	GetLoaderById(ctx context.Context, id int) (*dto.LoaderOutputDto, error)
	CreateLoader(ctx context.Context, loader *dto.LoaderOutputDto) error
	UpdateLoader(ctx context.Context, loader *dto.LoaderOutputDto) error
}

type UserService struct {
	userRepository     UserRepository
	customerRepository CustomerRepository
	loadersRepository  LoaderRepository
}

func NewUserService(userRepository UserRepository, customerRepository CustomerRepository, loaderRepository *repository.LoaderRepository) *UserService {
	return &UserService{
		userRepository:     userRepository,
		customerRepository: customerRepository,
		loadersRepository:  loaderRepository,
	}
}

func (s *UserService) Register(ctx context.Context, user *dto.UserSignUpDto) error {
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	id, err := s.userRepository.Register(ctx, user)
	if err != nil {
		return err
	}

	switch user.Role {
	case "customer":
		err := s.customerRepository.CreateCustomer(ctx, generateRandomCustomer(id))
		if err != nil {
			return err
		}
	case "loader":
		err := s.loadersRepository.CreateLoader(ctx, generateRandomLoader(id))
		if err != nil {
			return err
		}
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
func (s *UserService) CreateCustomer(ctx context.Context, customer *dto.CustomerSignUpDto) error {
	return s.customerRepository.CreateCustomer(ctx, customer)
}
func (s *UserService) GetCustomerById(ctx context.Context, id int) (dto.CustomerOutputDto, error) {
	return s.customerRepository.GetCustomerById(ctx, id)
}
func (s *UserService) GetAssignedLoaders(ctx context.Context, id int) ([]dto.LoaderOutputDto, error) {
	return s.loadersRepository.GetAssignedLoaders(ctx, id)
}
func (s *UserService) GetLoaderById(ctx context.Context, id int) (*dto.LoaderOutputDto, error) {
	return s.loadersRepository.GetLoaderById(ctx, id)
}
func (s *UserService) UpdateLoader(ctx context.Context, loader *dto.LoaderOutputDto) error {
	return s.loadersRepository.UpdateLoader(ctx, loader)
}
func (s *UserService) UpdateCustomer(ctx context.Context, outputDto dto.CustomerUpdateDto) error {
	return s.customerRepository.UpdateCustomer(ctx, outputDto)
}
func generateRandomCustomer(id int) *dto.CustomerSignUpDto {
	return &dto.CustomerSignUpDto{
		CustomerID:      id,
		StartingCapital: float64(rand.Intn(90001) + 10000),
	}
}
func generateRandomLoader(id int) *dto.LoaderOutputDto {
	maxWeight := float64(rand.Intn(11) + 20)
	drunkenness := rand.Intn(2) == 1
	Fatigue := float64(rand.Intn(71))
	Salary := float64(rand.Intn(20001) + 10000)

	loader := &dto.LoaderOutputDto{
		LoaderID:    id,
		MaxWeight:   maxWeight,
		Drunkenness: drunkenness,
		Fatigue:     Fatigue,
		Salary:      Salary,
	}

	game.Recalculate(loader)

	return loader
}
