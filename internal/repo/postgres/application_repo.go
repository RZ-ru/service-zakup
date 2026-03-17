package postgres

import (
	"context"
	"errors"

	"zakup/internal/request"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ request.Repository = (*ApplicationRepository)(nil)
var ErrApplicationNotFound = errors.New("application not found")

type ApplicationRepository struct {
	pool *pgxpool.Pool
}

func NewApplicationRepository(pool *pgxpool.Pool) *ApplicationRepository {
	return &ApplicationRepository{pool: pool}
}

func (r *ApplicationRepository) Create(ctx context.Context, app *request.Application) error {
	query := `
		INSERT INTO applications (
			id,
			author_id,
			product_id,
			quantity,
			comment,
			status,
			version,
			created_at,
			updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`

	_, err := r.pool.Exec(
		ctx,
		query,
		app.ID,
		app.AuthorID,
		app.ProductID,
		app.Quantity,
		app.Comment,
		string(app.Status),
		app.Version,
		app.CreatedAt,
		app.UpdatedAt,
	)

	return err
}

func (r *ApplicationRepository) GetByID(ctx context.Context, id uuid.UUID) (*request.Application, error) {
	query := `
		SELECT
			id,
			author_id,
			product_id,
			quantity,
			comment,
			status,
			version,
			created_at,
			updated_at
		FROM applications
		WHERE id = $1
	`

	var app request.Application
	var status string

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&app.ID,
		&app.AuthorID,
		&app.ProductID,
		&app.Quantity,
		&app.Comment,
		&status,
		&app.Version,
		&app.CreatedAt,
		&app.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrApplicationNotFound
		}
		return nil, err
	}

	app.Status = request.Status(status)

	return &app, nil
}

func (r *ApplicationRepository) List(ctx context.Context) ([]*request.Application, error) {
	query := `
		SELECT
			id,
			author_id,
			product_id,
			quantity,
			comment,
			status,
			version,
			created_at,
			updated_at
		FROM applications
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*request.Application

	for rows.Next() {
		var app request.Application
		var status string

		err := rows.Scan(
			&app.ID,
			&app.AuthorID,
			&app.ProductID,
			&app.Quantity,
			&app.Comment,
			&status,
			&app.Version,
			&app.CreatedAt,
			&app.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		app.Status = request.Status(status)
		result = append(result, &app)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *ApplicationRepository) Update(ctx context.Context, app *request.Application) error {
	query := `
		UPDATE applications
		SET
			author_id = $2,
			product_id = $3,
			quantity = $4,
			comment = $5,
			status = $6,
			version = $7,
			updated_at = $8
		WHERE id = $1
	`

	cmdTag, err := r.pool.Exec(
		ctx,
		query,
		app.ID,
		app.AuthorID,
		app.ProductID,
		app.Quantity,
		app.Comment,
		string(app.Status),
		app.Version,
		app.UpdatedAt,
	)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrApplicationNotFound
	}

	return nil
}
func (r *ApplicationRepository) CreateWithOutbox(
	ctx context.Context,
	app *request.Application,
	event *request.OutboxEvent,
) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	appQuery := `
		INSERT INTO applications (
			id,
			author_id,
			product_id,
			quantity,
			comment,
			status,
			version,
			created_at,
			updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`

	_, err = tx.Exec(
		ctx,
		appQuery,
		app.ID,
		app.AuthorID,
		app.ProductID,
		app.Quantity,
		app.Comment,
		string(app.Status),
		app.Version,
		app.CreatedAt,
		app.UpdatedAt,
	)
	if err != nil {
		return err
	}

	outboxQuery := `
		INSERT INTO outbox_events (
			id,
			aggregate_type,
			aggregate_id,
			event_type,
			routing_key,
			payload,
			status,
			attempts,
			error,
			created_at,
			published_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`

	_, err = tx.Exec(
		ctx,
		outboxQuery,
		event.ID,
		event.AggregateType,
		event.AggregateID,
		event.EventType,
		event.RoutingKey,
		event.Payload,
		string(event.Status),
		event.Attempts,
		event.Error,
		event.CreatedAt,
		event.PublishedAt,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
