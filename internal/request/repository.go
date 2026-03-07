package application // Говорит, какие операции хранения нужны.

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, app *Application) error
	GetByID(ctx context.Context, id uuid.UUID) (*Application, error)
	List(ctx context.Context) ([]*Application, error)
	Update(ctx context.Context, app *Application) error
}
