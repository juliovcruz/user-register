package users

import (
	"context"
	"fmt"
)

type repository interface {
	Create(ctx context.Context, user User) (User, error)
	Update(ctx context.Context, email, password string) error
	GetByEMail(ctx context.Context, email string) (User, error)
	GetAll(ctx context.Context, limit, offset int) ([]User, error)
}

type tokenService interface {
	Create(user User) (string, error)
	IsValid(tokenStr string) (bool, error)
}

type mailValidationService interface {
	Create(ctx context.Context, email string) error
	Validate(ctx context.Context, email string, code int) error
}

type hashService interface {
	Create(password string) (string, error)
	IsValid(inputPassword, password string) bool
}

type zipCodeService interface {
	GetAddressByZipCode(zipCode string) (Address, error)
}

type Service struct {
	repo                  repository
	tokenService          tokenService
	hashService           hashService
	zipCodeService        zipCodeService
	mailValidationService mailValidationService
}

func NewService(
	repo repository, tokenService tokenService,
	zipCodeService zipCodeService, hashService hashService,
	mailValidationService mailValidationService,
) *Service {
	return &Service{
		repo: repo, tokenService: tokenService,
		zipCodeService: zipCodeService, hashService: hashService,
		mailValidationService: mailValidationService,
	}
}

func (s *Service) Create(ctx context.Context, request CreateUser) (User, error) {
	if request.Password != request.ConfirmPassword {
		return User{}, ErrPasswordMismatch
	}

	address, err := s.zipCodeService.GetAddressByZipCode(request.ZipCode)
	if err != nil {
		return User{}, fmt.Errorf("failed to fetch address: %w", err)
	}

	hashPassword, err := s.hashService.Create(request.Password)
	if err != nil {
		return User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.repo.Create(ctx, User{
		Email:    request.Email,
		Name:     request.Name,
		Address:  address,
		Password: hashPassword,
	})
	if err != nil {
		return User{}, fmt.Errorf("failed to create user: %w", err)
	}

	user.Password = ""

	return user, nil
}

func (s *Service) UpdatePassword(ctx context.Context, request UpdatePassword) error {
	if request.Password != request.ConfirmPassword {
		return ErrPasswordMismatch
	}

	err := s.mailValidationService.Validate(ctx, request.Email, request.Code)
	if err != nil {
		return err
	}

	password, err := s.hashService.Create(request.Password)
	if err != nil {
		return err
	}

	err = s.repo.Update(ctx, request.Email, password)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEMail(ctx, email)
	if err != nil {
		return "", ErrUserNotFound
	}

	if !s.hashService.IsValid(password, user.Password) {
		return "", ErrInvalidLogin
	}

	token, err := s.tokenService.Create(user)
	if err != nil {
		return "", fmt.Errorf("failed to create token: %w", err)
	}

	return token, nil
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]User, error) {
	users, err := s.repo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

func (s *Service) ForgotPassword(ctx context.Context, email string) error {
	_, err := s.repo.GetByEMail(ctx, email)
	if err != nil {
		return err
	}

	if err := s.mailValidationService.Create(ctx, email); err != nil {
		return err
	}

	return nil
}
