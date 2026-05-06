package main

import (
	"log"
	"os"

	"user-service/internal/database"
	"user-service/internal/handlers"
	"user-service/internal/middleware"
	"user-service/internal/repository"
	"user-service/internal/services"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
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
	service := services.NewUserService(repo)
	handler := handlers.NewHandler(service)

	r.POST("/users", handler.CreateUser)
	r.GET("/users/:id", handler.GetUser)

	r.Run(":8080")
}
