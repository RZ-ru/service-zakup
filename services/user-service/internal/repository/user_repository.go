package repository

import (
	"context"
	"user-service/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
}

type InMemoryRepo struct{}

func (r *InMemoryRepo) Create(ctx context.Context, user *models.User) error {
	return nil
}
