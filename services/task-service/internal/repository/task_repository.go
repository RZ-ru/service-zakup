package repository

import (
	"context"
	"database/sql"
	"task-service/internal/models"
)

type TaskRepository interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)

	CreateTx(ctx context.Context, tx *sql.Tx, task *models.Task) error
	UpdateTx(ctx context.Context, tx *sql.Tx, task *models.Task) (*models.Task, error)
	DeleteTx(ctx context.Context, tx *sql.Tx, id string) error

	GetByID(ctx context.Context, id string) (*models.Task, error)
}
