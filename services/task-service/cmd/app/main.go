package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"task-service/internal/database"
	"task-service/internal/handlers"
	"task-service/internal/middleware"
	"task-service/internal/repository"
	"task-service/internal/services"
)

func main() {

	r := gin.New()

	r.Use(middleware.RequestID())
	r.Use(middleware.UserContext())
	r.Use(middleware.Logger())

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not set")
	}

	db := database.NewPostgres(dbURL)

	repo := repository.NewPostgresRepo(db)
	service := services.NewTaskService(repo)
	handler := handlers.NewHandler(service)

	r.POST("/tasks", handler.CreateTask)
	r.GET("/tasks/:id", handler.GetTask)

	r.Run(":8080")
}
