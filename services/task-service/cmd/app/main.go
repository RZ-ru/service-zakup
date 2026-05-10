package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"task-service/internal/broker"
	"task-service/internal/clients"
	"task-service/internal/database"
	"task-service/internal/handlers"
	"task-service/internal/middleware"
	"task-service/internal/repository"
	"task-service/internal/services"
	"task-service/internal/workers"
)

func main() {

	r := gin.New()

	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.Auth())

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not set")
	}

	db := database.NewPostgres(dbURL)

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		log.Fatal("RABBITMQ_URL not set")
	}

	rabbit, err := broker.NewRabbitMQ(rabbitURL)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbit.Close()
	//_ = rabbit

	permissionGRPCAddr := os.Getenv("PERMISSION_GRPC_ADDR")
	if permissionGRPCAddr == "" {
		permissionGRPCAddr = "permission-service:9090"
	}

	permClient, err := clients.NewPermissionClient(permissionGRPCAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer permClient.Close()

	taskRepo := repository.NewPostgresRepo(db)
	outboxRepo := repository.NewPostgresOutboxRepo(db)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	outboxWorker := workers.NewOutboxWorker(outboxRepo, rabbit)
	go outboxWorker.Start(ctx)

	service := services.NewTaskService(taskRepo, outboxRepo, permClient)
	handler := handlers.NewHandler(service)

	r.POST("/tasks", handler.CreateTask)
	r.GET("/tasks/:id", handler.GetTask)
	r.PATCH("/tasks/:id", handler.UpdateTask)
	r.DELETE("/tasks/:id", handler.DeleteTask)

	//8080
	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}
