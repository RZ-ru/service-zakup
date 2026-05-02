package models

import "time"

type Task struct {
	ID          string
	Title       string
	Description string
	Status      string // new, in_progress, done
	UserID      string // owner
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
