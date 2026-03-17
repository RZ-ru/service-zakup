package request // Интерфейс отправки событий

import "context"

// EventPublisher — контракт для отправки событий наружу.
// Сервис будет знать только про этот интерфейс.
type EventPublisher interface {
	PublishApplicationCreated(ctx context.Context, event ApplicationCreatedEvent) error
	PublishApplicationStatusChanged(ctx context.Context, event ApplicationStatusChangedEvent) error
}
