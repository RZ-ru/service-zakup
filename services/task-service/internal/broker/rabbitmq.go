package broker

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const TaskExchange = "task.events"

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	var conn *amqp.Connection
	var err error

	for i := 1; i <= 10; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}

		fmt.Printf("rabbitmq not ready, retry %d/10: %v\n", i, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("connect rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("open rabbitmq channel: %w", err)
	}

	r := &RabbitMQ{
		conn:    conn,
		channel: ch,
	}

	if err := r.declare(); err != nil {
		_ = r.Close()
		return nil, err
	}

	return r, nil
}

func (r *RabbitMQ) declare() error {
	err := r.channel.ExchangeDeclare(
		TaskExchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("declare exchange: %w", err)
	}

	return nil
}

func (r *RabbitMQ) Publish(ctx context.Context, routingKey string, payload []byte) error {
	err := r.channel.PublishWithContext(
		ctx,
		TaskExchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         payload,
		},
	)
	if err != nil {
		return fmt.Errorf("publish message: %w", err)
	}

	return nil
}

func (r *RabbitMQ) Close() error {
	if r.channel != nil {
		_ = r.channel.Close()
	}

	if r.conn != nil {
		return r.conn.Close()
	}

	return nil
}
