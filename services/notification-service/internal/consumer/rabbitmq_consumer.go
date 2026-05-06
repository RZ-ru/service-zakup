package consumer

import (
	"encoding/json"
	"fmt"
	"log"

	"notification-service/internal/clients"
	"notification-service/internal/email"
	"notification-service/internal/events"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	exchangeName = "task.events"
	queueName    = "notification.task.events"
)

type Consumer struct {
	channel    *amqp.Channel
	userClient *clients.UserClient
	email      *email.Sender
}

func NewConsumer(ch *amqp.Channel, userClient *clients.UserClient, emailSender *email.Sender) *Consumer {
	return &Consumer{
		channel:    ch,
		userClient: userClient,
		email:      emailSender,
	}
}

func (c *Consumer) Start() error {
	err := c.channel.ExchangeDeclare(
		exchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	queue, err := c.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	routingKeys := []string{
		"task.created",
		"task.updated",
		"task.status_changed",
		"task.deleted",
	}

	for _, key := range routingKeys {
		if err := c.channel.QueueBind(queue.Name, key, exchangeName, false, nil); err != nil {
			return err
		}
	}

	messages, err := c.channel.Consume(
		queue.Name,
		"notification-service",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	log.Println("notification-service started")
	log.Printf("listening queue=%s", queueName)

	for msg := range messages {
		if err := c.handleMessage(msg); err != nil {
			log.Printf("handle message error: %v", err)
			_ = msg.Nack(false, true)
			continue
		}

		_ = msg.Ack(false)
	}

	return nil
}

func (c *Consumer) handleMessage(msg amqp.Delivery) error {
	var event events.TaskEvent

	if err := json.Unmarshal(msg.Body, &event); err != nil {
		return err
	}

	if event.UserID == "" {
		return fmt.Errorf("event user_id is empty")
	}

	user, err := c.userClient.GetByID(event.UserID)
	if err != nil {
		return err
	}

	subject, body := buildEmail(msg.RoutingKey, event)

	return c.email.Send(user.Email, subject, body)
}

func buildEmail(routingKey string, event events.TaskEvent) (string, string) {
	switch routingKey {
	case "task.created":
		return "Создана новая задача", fmt.Sprintf("Создана задача: %s", event.Title)

	case "task.updated":
		return "Задача изменена", fmt.Sprintf("Изменена задача: %s", event.Title)

	case "task.status_changed":
		return "Статус задачи изменён", fmt.Sprintf("Задача %s изменила статус: %s → %s", event.Title, event.OldStatus, event.NewStatus)

	case "task.deleted":
		return "Задача удалена", fmt.Sprintf("Удалена задача: %s", event.Title)

	default:
		return "Изменение задачи", fmt.Sprintf("Событие по задаче: %s", event.Title)
	}
}
