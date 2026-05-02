package repository

import (
	"context"
	"database/sql"
	"task-service/internal/models"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) Create(ctx context.Context, task *models.Task) error {
	query := `
	INSERT INTO tasks (id, title, description, status, user_id, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx,
		query,
		task.ID,
		task.Title,
		task.Description,
		task.Status,
		task.UserID,
		task.CreatedAt,
		task.UpdatedAt,
	)
	return err
}

func (r *PostgresRepo) GetByID(ctx context.Context, id string) (*models.Task, error) {
	query := `
        SELECT id, title, description, status, user_id
        FROM tasks
        WHERE id = $1
    `

	row := r.db.QueryRowContext(ctx, query, id)

	var t models.Task
	err := row.Scan(
		&t.ID,
		&t.Title,
		&t.Description,
		&t.Status,
		&t.UserID,
	)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
