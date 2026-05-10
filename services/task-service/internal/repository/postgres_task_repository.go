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

func (r *PostgresRepo) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

func (r *PostgresRepo) CreateTx(ctx context.Context, tx *sql.Tx, task *models.Task) error {
	query := `
		INSERT INTO tasks (id, title, description, status, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := tx.ExecContext(
		ctx,
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

func (r *PostgresRepo) UpdateTx(ctx context.Context, tx *sql.Tx, task *models.Task) (*models.Task, error) {
	query := `
		UPDATE tasks
		SET title = $2,
		    description = $3,
		    status = $4,
		    updated_at = $5
		WHERE id = $1
		RETURNING id, title, description, status, user_id, created_at, updated_at
	`

	row := tx.QueryRowContext(
		ctx,
		query,
		task.ID,
		task.Title,
		task.Description,
		task.Status,
		task.UpdatedAt,
	)

	var updated models.Task
	err := row.Scan(
		&updated.ID,
		&updated.Title,
		&updated.Description,
		&updated.Status,
		&updated.UserID,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *PostgresRepo) DeleteTx(ctx context.Context, tx *sql.Tx, id string) error {
	query := `
		DELETE FROM tasks
		WHERE id = $1
	`

	_, err := tx.ExecContext(ctx, query, id)

	return err
}
