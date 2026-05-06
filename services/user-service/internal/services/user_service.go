package services

import (
	"context"
	"errors"
	"user-service/internal/models"
	"user-service/internal/repository"

	"github.com/google/uuid"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Create(ctx context.Context, email, name string) (*models.User, error) {
	if email == "" {
		return nil, errors.New("email required")
	}

	if name == "" {
		return nil, errors.New("name required")
	}

	user := &models.User{
		ID:    uuid.NewString(),
		Email: email,
		Name:  name,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id string) (*models.User, error) {
	if id == "" {
		return nil, errors.New("user_id required")
	}

	return s.repo.GetByID(ctx, id)
}
