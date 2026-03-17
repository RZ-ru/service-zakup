package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"zakup/internal/config"
	"zakup/internal/handler"
	"zakup/internal/migrator"
	"zakup/internal/repo/postgres"
	"zakup/internal/request"
)

func main() {
	//	Загружаем конфиг
	cfg := config.MustLoad()

	// 2. Подключаем миграции
	if err := migrator.Up(cfg.PostgresDSN, "./internal/migrations"); err != nil {
		log.Fatalf("run migrations: %v", err)
	}
	log.Println("migrations applied")

	// 3. Контекст для инициализации
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 4. Подключение к Postgres
	pgPool, err := pgxpool.New(ctx, cfg.PostgresDSN)
	if err != nil {
		log.Fatalf("postgres connect error: %v", err)
	}
	defer pgPool.Close()

	if err := pgPool.Ping(ctx); err != nil {
		log.Fatalf("postgres ping error: %v", err)
	}

	// 5. Подключение к RabbitMQ
	/* rmqConn, err := amqp.Dial(cfg.RabbitMQURL)
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
		log.Fatalf("rabbitmq exchange declare error: %v", err)
	} */

	// 6. Создаем repository
	applicationRepo := postgres.NewApplicationRepository(pgPool)

	// 7. Создаем service
	applicationService := request.NewService(applicationRepo)

	// 8. Создаем handler
	applicationHandler := handler.NewApplicationHandler(applicationService)

	// 9. Роутер
	r := gin.Default()
	applicationHandler.Register(r)

	// 10. HTTP сервер
	server := &http.Server{
		Addr:         cfg.HTTPAddress,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	// 11. Запуск сервера в горутине
	go func() {
		log.Printf("HTTP server started on %s", cfg.HTTPAddress)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server error: %v", err)
		}
	}()

	// 12. Ожидание сигнала завершения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("shutting down server...")

	// 13. Корректное завершение
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	log.Println("server stopped")
}
