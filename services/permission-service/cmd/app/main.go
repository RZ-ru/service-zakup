package main

import (
	"context"
	"log"
	"net"
	"os"
	"permission-service/internal/cache"
	"permission-service/internal/database"
	"permission-service/internal/grpcserver"
	"permission-service/internal/handlers"
	"permission-service/internal/middleware"
	"permission-service/internal/repository"
	"permission-service/internal/services"
	permissionpb "permission-service/proto"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
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

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}

	db := database.NewPostgres(dbURL)

	redisClient := cache.NewRedis(redisAddr)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx); err != nil {
		log.Printf("redis unavailable: %v", err)
	} else {
		log.Println("redis connected")
	}

	repo := repository.NewPostgresRepo(db)
	service := services.NewPermissionService(repo, redisClient)
	handler := handlers.NewHandler(service)

	grpcAddr := os.Getenv("GRPC_ADDR")
	if grpcAddr == "" {
		grpcAddr = ":9090"
	}

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	permissionpb.RegisterPermissionServiceServer(
		grpcServer,
		grpcserver.NewPermissionServer(service),
	)

	go func() {
		log.Printf("permission gRPC server started on %s", grpcAddr)

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	defer grpcServer.GracefulStop()

	r.POST("/permissions", handler.Create)
	r.GET("/permissions/check", handler.Check)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
