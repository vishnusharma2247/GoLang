package service

import (
	"context"
	"errors"
	"strings"

	"http-learning/internal/domain"
	"http-learning/internal/repository/postgres"
	"http-learning/internal/security"
)

type AuthService struct {
	userRepo *postgres.UserRepository
}

var ErrInvalidInput = errors.New("invalid input")
var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrInvalidCredentials = errors.New("invalid credentials")

func NewAuthService(userRepo *postgres.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, email, plainPassword string) (*domain.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || len(plainPassword) < 6 {
		return nil, ErrInvalidInput
	}

	existingUser, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	passwordHash, err := security.HashPassword(plainPassword)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:        email,
		PasswordHash: passwordHash,
		Role:         "user",
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, plainPassword string) (*domain.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || plainPassword == "" {
		return nil, ErrInvalidInput
	}

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := security.CheckPassword(user.PasswordHash, plainPassword); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
