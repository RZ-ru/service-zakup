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

func NewUserService(r repository.UserRepository) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) Create(ctx context.Context, email, name string) (*models.User, error) {

	if email == "" {
		return nil, errors.New("email required")
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
