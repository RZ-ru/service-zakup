package request

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type OutboxStatus string

const (
	OutboxStatusPending   OutboxStatus = "pending"
	OutboxStatusPublished OutboxStatus = "published"
	OutboxStatusFailed    OutboxStatus = "failed"
)

type OutboxEvent struct {
	ID            uuid.UUID
	AggregateType string
	AggregateID   uuid.UUID
	EventType     string
	RoutingKey    string
	Payload       json.RawMessage
	Status        OutboxStatus
	Attempts      int
	Error         string
	CreatedAt     time.Time
	PublishedAt   *time.Time
}

func NewOutboxEvent(
	aggregateType string,
	aggregateID uuid.UUID,
	eventType string,
	routingKey string,
	payload []byte,
) *OutboxEvent {
	return &OutboxEvent{
		ID:            uuid.New(),
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		EventType:     eventType,
		RoutingKey:    routingKey,
		Payload:       payload,
		Status:        OutboxStatusPending,
		Attempts:      0,
		Error:         "",
		CreatedAt:     time.Now(),
	}
}
