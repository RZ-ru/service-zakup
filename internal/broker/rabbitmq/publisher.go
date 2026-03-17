package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"zakup/internal/request"
)

// Publisher отправляет события в RabbitMQ.
type Publisher struct {
	ch       *amqp.Channel
	exchange string
}

// NewPublisher создаёт publisher.
func NewPublisher(ch *amqp.Channel, exchange string) *Publisher {
	return &Publisher{
		ch:       ch,
		exchange: exchange,
	}
}

// PublishApplicationCreated публикует событие о создании заявки.
func (p *Publisher) PublishApplicationCreated(ctx context.Context, event request.ApplicationCreatedEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal application created event: %w", err)
	}

	err = p.ch.PublishWithContext(
		ctx,
		p.exchange,            // exchange
		"application.created", // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Type:         "application.created",
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("publish application created event: %w", err)
	}

	return nil
}
