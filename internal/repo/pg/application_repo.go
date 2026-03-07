package pg

import (
	"context"
	"errors"

	"zakup/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ApplicationRepo struct {
	db *pgx.Conn
}

func NewApplicationRepo(db *pgx.Conn) *ApplicationRepo { return &ApplicationRepo{db: db} }

func (r *ApplicationRepo) ProductExists(ctx context.Context, productID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM products WHERE id=$1)`, productID).Scan(&exists)
	return exists, err
}

func (r *ApplicationRepo) Create(ctx context.Context, app *models.Application) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO applications
			(id, product_id, author_id, department, amount, status, created_at, updated_at, version)
		VALUES
			($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`, app.ID, app.ProductID, app.AuthorID, app.Department, app.Amount, string(app.Status), app.CreatedAt, app.UpdatedAt, app.Version)
	return err
}

func (r *ApplicationRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Application, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, product_id, author_id, department, amount, status, created_at, updated_at, version
		FROM applications
		WHERE id=$1
	`, id)

	var a models.Application
	var status string
	if err := row.Scan(&a.ID, &a.ProductID, &a.AuthorID, &a.Department, &a.Amount, &status, &a.CreatedAt, &a.UpdatedAt, &a.Version); err != nil {
		return nil, err
	}
	a.Status = models.Status(status)
	return &a, nil
}

var ErrConflict = errors.New("conflict")

func (r *ApplicationRepo) UpdateStatus(ctx context.Context, id uuid.UUID, newStatus models.Status, newVersion int) error {
	// optimistic locking: обновляем только если версия совпала
	tag, err := r.db.Exec(ctx, `
		UPDATE applications
		SET status=$2, updated_at=now(), version=$3
		WHERE id=$1 AND version=$4
	`, id, string(newStatus), newVersion, newVersion-1)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrConflict
	}
	return nil
}
