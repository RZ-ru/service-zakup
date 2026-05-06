package main

import (
	"log"
	"os"
	"time"

	"notification-service/internal/clients"
	"notification-service/internal/consumer"
	"notification-service/internal/email"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		log.Fatal("RABBITMQ_URL not set")
	}

	userServiceURL := os.Getenv("USER_SERVICE_URL")
	if userServiceURL == "" {
		log.Fatal("USER_SERVICE_URL not set")
	}

	conn, err := connectWithRetry(rabbitURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	userClient := clients.NewUserClient(userServiceURL)
	emailSender := email.NewSender()

	c := consumer.NewConsumer(ch, userClient, emailSender)

	if err := c.Start(); err != nil {
		log.Fatal(err)
	}
}

func connectWithRetry(url string) (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error

	for i := 1; i <= 10; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			return conn, nil
		}

		log.Printf("rabbitmq not ready, retry %d/10: %v", i, err)
		time.Sleep(2 * time.Second)
	}

	return nil, err
}
