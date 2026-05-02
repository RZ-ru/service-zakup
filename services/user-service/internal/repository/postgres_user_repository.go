package repository

import (
	"context"
	"database/sql"
	"user-service/internal/models"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) Create(ctx context.Context, user *models.User) error {
	query := `
        INSERT INTO users (id, email, name)
        VALUES ($1, $2, $3)
    `

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.Name,
	)

	return err
}
