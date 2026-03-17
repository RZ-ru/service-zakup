package postgres

import (
	"context"
	"encoding/json"
	"time"

	"zakup/internal/request"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OutboxRepository struct {
	pool *pgxpool.Pool
}

func NewOutboxRepository(pool *pgxpool.Pool) *OutboxRepository {
	return &OutboxRepository{pool: pool}
}

func (r *OutboxRepository) FetchPendingEvents(ctx context.Context, limit int) ([]*request.OutboxEvent, error) {
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
			error,
			created_at,
			published_at
		FROM outbox_events
		WHERE status = 'pending'
		ORDER BY created_at ASC
		LIMIT $1
	`

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*request.OutboxEvent

	for rows.Next() {
		var e request.OutboxEvent
		var payload []byte
		var status string

		err := rows.Scan(
			&e.ID,
			&e.AggregateType,
			&e.AggregateID,
			&e.EventType,
			&e.RoutingKey,
			&payload,
			&status,
			&e.Attempts,
			&e.Error,
			&e.CreatedAt,
			&e.PublishedAt,
		)
		if err != nil {
			return nil, err
		}

		e.Payload = json.RawMessage(payload)
		e.Status = request.OutboxStatus(status)

		events = append(events, &e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (r *OutboxRepository) MarkPublished(ctx context.Context, id uuid.UUID, publishedAt time.Time) error {
	query := `
		UPDATE outbox_events
		SET
			status = 'published',
			published_at = $2,
			error = ''
		WHERE id = $1
	`
	_, err := r.pool.Exec(ctx, query, id, publishedAt)
	return err
}

func (r *OutboxRepository) MarkFailed(ctx context.Context, id uuid.UUID, errText string) error {
	query := `
		UPDATE outbox_events
		SET
			attempts = attempts + 1,
			error = $2
		WHERE id = $1
	`
	_, err := r.pool.Exec(ctx, query, id, errText)
	return err
}
