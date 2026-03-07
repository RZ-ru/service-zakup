package repo

import (
	"context"

	"zakup/internal/models"

	"github.com/google/uuid"
)

type ApplicationRepo interface {
	Create(ctx context.Context, app *models.Application) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Application, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, newStatus models.Status, newVersion int) error
	ProductExists(ctx context.Context, productID uuid.UUID) (bool, error)
}
