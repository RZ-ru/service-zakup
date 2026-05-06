package workers

import (
	"context"
	"log"
	"time"

	"task-service/internal/broker"
	"task-service/internal/repository"
)

type OutboxWorker struct {
	outbox repository.OutboxRepository
	broker *broker.RabbitMQ
}

func NewOutboxWorker(outbox repository.OutboxRepository, broker *broker.RabbitMQ) *OutboxWorker {
	return &OutboxWorker{
		outbox: outbox,
		broker: broker,
	}
}

func (w *OutboxWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	log.Println("outbox worker started")

	for {
		select {
		case <-ctx.Done():
			log.Println("outbox worker stopped")
			return

		case <-ticker.C:
			w.process(ctx)
		}
	}
}

func (w *OutboxWorker) process(ctx context.Context) {
	events, err := w.outbox.GetPending(ctx, 10)
	if err != nil {
		log.Printf("outbox get pending error: %v", err)
		return
	}

	for _, event := range events {
		err := w.broker.Publish(ctx, event.RoutingKey, event.Payload)
		if err != nil {
			log.Printf("publish outbox event error: id=%s err=%v", event.ID, err)

			if markErr := w.outbox.MarkFailed(ctx, event.ID); markErr != nil {
				log.Printf("mark failed error: id=%s err=%v", event.ID, markErr)
			}

			continue
		}

		if err := w.outbox.MarkProcessed(ctx, event.ID); err != nil {
			log.Printf("mark processed error: id=%s err=%v", event.ID, err)
			continue
		}

		log.Printf("outbox event processed: id=%s routing_key=%s", event.ID, event.RoutingKey)
	}
}
