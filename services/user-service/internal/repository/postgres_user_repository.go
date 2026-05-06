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
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO users (id, email, name)
		VALUES ($1, $2, $3)
	`, user.ID, user.Email, user.Name)

	return err
}

func (r *PostgresRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, email, name
		FROM users
		WHERE id = $1
	`, id)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Name,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
