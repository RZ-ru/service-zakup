package services

import (
	"context"
	"errors"
	"user-service/internal/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	//GetByID(id string) (*models.User, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
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
