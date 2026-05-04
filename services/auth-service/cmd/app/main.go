package main

import (
	"github.com/gin-gonic/gin"

	"auth-service/internal/handlers"
	"auth-service/internal/services"
)

func main() {
	r := gin.Default()

	svc := services.NewAuthService()
	h := handlers.NewHandler(svc)

	r.POST("/login", h.Login)

	r.Run(":8080")
}
