package repository

import (
	"database/sql"
	"permission-service/internal/models"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) Create(p *models.Permission) error {
	_, err := r.db.Exec(`
		INSERT INTO permissions (user_id, task_id, role)
		VALUES ($1, $2, $3)
	`, p.UserID, p.TaskID, p.Role)
	return err
}

func (r *PostgresRepo) GetRole(userID, taskID string) (string, error) {
	var role string

	err := r.db.QueryRow(`
		SELECT role
		FROM permissions
		WHERE user_id = $1 AND task_id = $2
	`, userID, taskID).Scan(&role)

	return role, err
}
