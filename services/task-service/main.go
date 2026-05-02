package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"task-service/internal/db"
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

	// подключение
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	// ждём БД
	waitForDB(dbConn)

	// миграции
	db.RunMigrations(dbURL)

	if err := dbConn.Ping(); err != nil {
		log.Fatal(err)
	}

	repo := repository.NewPostgresRepo(dbConn)
	service := services.NewTaskService(repo)
	handler := handlers.NewHandler(service)

	r.POST("/tasks", handler.CreateTask)
	r.GET("/tasks/:id", handler.GetTask)

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
