package request // Определяет интерфейс хранилища заявок.

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, app *Application) error
	GetByID(ctx context.Context, id uuid.UUID) (*Application, error)
	List(ctx context.Context) ([]*Application, error)
	Update(ctx context.Context, app *Application) error
	CreateWithOutbox(ctx context.Context, app *Application, event *OutboxEvent) error
}
