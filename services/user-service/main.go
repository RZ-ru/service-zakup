package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"user-service/internal/db"
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

	// подключение
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	waitForDB(dbConn)

	// миграции
	db.RunMigrations(dbURL)

	if err := dbConn.Ping(); err != nil {
		log.Fatal(err)
	}

	repo := repository.NewPostgresRepo(dbConn)
	service := services.NewUserService(repo)
	handler := handlers.NewHandler(service)

	r.POST("/users", handler.CreateUser)

	r.Run(":8080")
}

func waitForDB(db *sql.DB) {
	for i := 0; i < 10; i++ {
		err := db.Ping()
		if err == nil {
			log.Println("DB connected")
			return
		}

		log.Println("waiting for DB...")
		time.Sleep(2 * time.Second)
	}

	log.Fatal("cannot connect to DB")
}
