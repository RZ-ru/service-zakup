package database

import (
	"database/sql"
	"log"
	"time"
)

func NewPostgres(dbURL string) *sql.DB {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		if err := db.Ping(); err == nil {
			return db
		}
		time.Sleep(2 * time.Second)
	}

	log.Fatal("cannot connect to db")
	return nil
}
