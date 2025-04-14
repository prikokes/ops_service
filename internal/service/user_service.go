package service

import (
	"context"
	"errors"
	"fmt"
	"ops_service/internal/model"
	"ops_service/internal/repository"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, email, password string, role model.UserRole) (int, error) {
	existingUser, err := s.repo.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, errors.New("user not found")) {
		return 0, fmt.Errorf("failed to check existing user: %w", err)
	}

	if existingUser != nil {
		return 0, errors.New("user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create the user
	user := &model.User{
		Email:     email,
		Password:  string(hashedPassword),
		Role:      role,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return user.ID, nil
}

func (s *UserService) AuthenticateUser(ctx context.Context, email, password string) (*model.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	return s.repo.GetByID(ctx, id)
}
