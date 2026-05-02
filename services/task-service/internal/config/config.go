package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddress      string
	PostgresDSN      string
	RabbitMQURL      string
	RabbitMQExchange string
}

func MustLoad() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		HTTPAddress:      getEnv("HTTP_ADDRESS", ":8080"),
		PostgresDSN:      getEnv("POSTGRES_DSN", ""),
		RabbitMQURL:      getEnv("RABBITMQ_URL", ""),
		RabbitMQExchange: getEnv("RABBITMQ_EXCHANGE", "application.events"),
	}

	if cfg.PostgresDSN == "" {
		log.Fatal("Postgress_DSN is required")
	}
	/*
		if cfg.RabbitMQURL == "" {
			log.Fatal("RabbitMQ_URL is required")
		}
	*/
	return cfg

}

func getEnv(key string, defaulValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaulValue
	}
	return value
}
