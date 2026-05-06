package models

import "time"

type OutboxEvent struct {
	ID            string
	AggregateType string
	AggregateID   string
	EventType     string
	RoutingKey    string
	Payload       []byte
	Status        string
	Attempts      int
	CreatedAt     time.Time
	ProcessedAt   *time.Time
}
