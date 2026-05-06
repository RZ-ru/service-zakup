package repository

import (
	"context"
	"database/sql"
	"task-service/internal/models"
)

type OutboxRepository interface {
	CreateTx(ctx context.Context, tx *sql.Tx, event *models.OutboxEvent) error
	GetPending(ctx context.Context, limit int) ([]models.OutboxEvent, error)
	MarkProcessed(ctx context.Context, id string) error
	MarkFailed(ctx context.Context, id string) error
}
