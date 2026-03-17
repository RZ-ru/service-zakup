package request // Структуры событий

import (
	"time"

	"github.com/google/uuid"
)

// ApplicationCreatedEvent — событие о создании заявки.
type ApplicationCreatedEvent struct {
	ApplicationID uuid.UUID `json:"application_id"`
	AuthorID      uuid.UUID `json:"author_id"`
	ProductID     uuid.UUID `json:"product_id"`
	Quantity      int32     `json:"quantity"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

// ApplicationStatusChangedEvent — событие о смене статуса заявки.
type ApplicationStatusChangedEvent struct {
	ApplicationID uuid.UUID `json:"application_id"`
	OldStatus     string    `json:"old_status"`
	NewStatus     string    `json:"new_status"`
	ChangedAt     time.Time `json:"changed_at"`
	Version       int       `json:"version"`
}
