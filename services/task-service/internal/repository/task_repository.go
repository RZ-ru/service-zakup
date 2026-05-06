package repository

import (
	"context"
	"database/sql"
	"task-service/internal/models"
)

type TaskRepository interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	CreateTx(ctx context.Context, tx *sql.Tx, task *models.Task) error
	GetByID(ctx context.Context, id string) (*models.Task, error)

	//GetTasks(ctx context.Context) ([]models.Task, error)
	//UpdateByID(ctx context.Context, id string) (*models.Task, error)
	//DeleteByID(ctx context.Context, id string) error
}
