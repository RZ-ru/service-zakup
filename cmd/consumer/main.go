package main

import (
	"context"
	"log"
	"time"

	"zakup/internal/broker"
	"zakup/internal/broker/rabbitmq"
	"zakup/internal/config"
	"zakup/internal/repo/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	cfg := config.MustLoad()
	if cfg.RabbitMQURL == "" {
		log.Fatal("RABBITMQ_URL is required for consumer")
	}

	ctx := context.Background()

	pgPool, err := pgxpool.New(ctx, cfg.PostgresDSN)
	if err != nil {
		log.Fatalf("postgres connect error: %v", err)
	}
	defer pgPool.Close()

	rmqConn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("rabbitmq connect error: %v", err)
	}
	defer rmqConn.Close()

	rmqChannel, err := rmqConn.Channel()
	if err != nil {
		log.Fatalf("rabbitmq channel error: %v", err)
	}
	defer rmqChannel.Close()

	if err := rmqChannel.ExchangeDeclare(
		cfg.RabbitMQExchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		log.Fatalf("exchange declare error: %v", err)
	}

	outboxRepo := postgres.NewOutboxRepository(pgPool)
	publisher := rabbitmq.NewPublisher(rmqChannel, cfg.RabbitMQExchange)
	relay := broker.NewOutboxRelay(outboxRepo, publisher)

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	log.Println("outbox relay started")

	for {
		runCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		if err := relay.RunOnce(runCtx, 50); err != nil {
			log.Printf("relay run error: %v", err)
		}
		cancel()

		<-ticker.C
	}
}
