package main

import (
	"log"
	"os"
	"permission-service/internal/database"
	"permission-service/internal/handlers"
	"permission-service/internal/middleware"
	"permission-service/internal/repository"
	"permission-service/internal/services"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
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

	db := database.NewPostgres(os.Getenv("DB_URL"))

	repo := repository.NewPostgresRepo(db)
	service := services.NewPermissionService(repo)
	handler := handlers.NewHandler(service)

	r.POST("/permissions", handler.Create)
	r.GET("/permissions/check", handler.Check)

	r.Run(":8080")

}
