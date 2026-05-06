package repository

import (
	"context"
	"database/sql"
	"task-service/internal/models"
)

type PostgresOutboxRepo struct {
	db *sql.DB
}

func NewPostgresOutboxRepo(db *sql.DB) *PostgresOutboxRepo {
	return &PostgresOutboxRepo{db: db}
}

func (r *PostgresOutboxRepo) CreateTx(ctx context.Context, tx *sql.Tx, event *models.OutboxEvent) error {
	query := `
		INSERT INTO outbox_events (
			id,
			aggregate_type,
			aggregate_id,
			event_type,
			routing_key,
			payload,
			status,
			attempts,
			created_at,
			processed_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		event.ID,
		event.AggregateType,
		event.AggregateID,
		event.EventType,
		event.RoutingKey,
		event.Payload,
		event.Status,
		event.Attempts,
		event.CreatedAt,
		event.ProcessedAt,
	)

	return err
}

func (r *PostgresOutboxRepo) GetPending(ctx context.Context, limit int) ([]models.OutboxEvent, error) {
	query := `
		SELECT
			id,
			aggregate_type,
			aggregate_id,
			event_type,
			routing_key,
			payload,
			status,
			attempts,
			created_at,
			processed_at
		FROM outbox_events
		WHERE status = 'pending'
		ORDER BY created_at
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]models.OutboxEvent, 0)

	for rows.Next() {
		var event models.OutboxEvent

		err := rows.Scan(
			&event.ID,
			&event.AggregateType,
			&event.AggregateID,
			&event.EventType,
			&event.RoutingKey,
			&event.Payload,
			&event.Status,
			&event.Attempts,
			&event.CreatedAt,
			&event.ProcessedAt,
		)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (r *PostgresOutboxRepo) MarkProcessed(ctx context.Context, id string) error {
	query := `
		UPDATE outbox_events
		SET status = 'processed',
		    processed_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)

	return err
}

func (r *PostgresOutboxRepo) MarkFailed(ctx context.Context, id string) error {
	query := `
		UPDATE outbox_events
		SET
    		attempts = attempts + 1,
    		status = CASE
        		WHEN attempts + 1 >= 5 THEN 'failed'
        		ELSE 'pending'
    		END
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)

	return err
}
