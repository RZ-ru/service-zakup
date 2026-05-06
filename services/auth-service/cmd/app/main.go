package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"auth-service/internal/handlers"
	"auth-service/internal/services"
)

func main() {
	r := gin.Default()

	svc, err := services.NewAuthService()
	if err != nil {
		log.Fatal(err)
	}
	h := handlers.NewHandler(svc)

	r.POST("/login", h.Login)

	r.Run(":8080")
}
