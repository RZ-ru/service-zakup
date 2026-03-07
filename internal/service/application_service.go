package service

import (
	"context"
	"errors"

	"zakup/internal/models"
	"zakup/internal/repo"
	"zakup/validation_service"

	"github.com/google/uuid"
)

var ErrProductNotFound = errors.New("product not found")

type ApplicationService struct {
	repo repo.ApplicationRepo
}

func NewApplicationService(r repo.ApplicationRepo) *ApplicationService {
	return &ApplicationService{repo: r}
}

func (s *ApplicationService) Create(ctx context.Context, authorID uuid.UUID, in validation_service.CreateApplicationInput) (*models.Application, error) {
	// 1) валидация входа
	if err := validation_service.ValidateCreateApplication(in); err != nil {
		return nil, err
	}

	// 2) проверка, что продукт существует
	ok, err := s.repo.ProductExists(ctx, in.ProductID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrProductNotFound
	}

	// 3) создаём доменную модель
	app := models.NewApplication(in.ProductID, authorID, in.Department, in.Amount)

	// 4) сохраняем
	if err := s.repo.Create(ctx, app); err != nil {
		return nil, err
	}
	return app, nil
}

func (s *ApplicationService) ChangeStatus(ctx context.Context, id uuid.UUID, next models.Status) (*models.Application, error) {
	app, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := app.SetStatus(next); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateStatus(ctx, app.ID, app.Status, app.Version); err != nil {
		return nil, err
	}

	return app, nil
}
