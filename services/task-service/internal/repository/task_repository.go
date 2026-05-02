package repository // Определяет интерфейс хранилища заявок.

import (
	"context"
	"task-service/internal/models"
)

type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	GetByID(ctx context.Context, id string) (*models.Task, error)
}
