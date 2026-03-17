package broker

import (
	"context"
	"time"

	"zakup/internal/request"

	"github.com/google/uuid"
)

type OutboxRepo interface {
	FetchPendingEvents(ctx context.Context, limit int) ([]*request.OutboxEvent, error)
	MarkPublished(ctx context.Context, id uuid.UUID, publishedAt time.Time) error
	MarkFailed(ctx context.Context, id uuid.UUID, errText string) error
}

type Publisher interface {
	Publish(ctx context.Context, routingKey string, body []byte) error
}

type OutboxRelay struct {
	repo      OutboxRepo
	publisher Publisher
}

func NewOutboxRelay(repo OutboxRepo, publisher Publisher) *OutboxRelay {
	return &OutboxRelay{
		repo:      repo,
		publisher: publisher,
	}
}

func (r *OutboxRelay) RunOnce(ctx context.Context, batchSize int) error {
	events, err := r.repo.FetchPendingEvents(ctx, batchSize)
	if err != nil {
		return err
	}

	for _, event := range events {
		err := r.publisher.Publish(ctx, event.RoutingKey, event.Payload)
		if err != nil {
			_ = r.repo.MarkFailed(ctx, event.ID, err.Error())
			continue
		}

		if err := r.repo.MarkPublished(ctx, event.ID, time.Now()); err != nil {
			return err
		}
	}

	return nil
}
