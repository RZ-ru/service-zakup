package models

import (
	"time"
)

type Status string

const (
	StatusNew       Status = "new"
	StatusInProcess Status = "in_progress"
	StatusApproved  Status = "approved"
	StatusRejected  Status = "rejected"
	StatusCanceled  Status = "canceled"
)

type Application struct {
	ID          int    `json:"id"`
	ProductName string `json:"product_name"`
	//AuthorID    uuid.UUID `json:"author_id"`
	Department string    `json:"department"`
	Amount     float64   `json:"amount"`
	Status     Status    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Version    int       `json:"version"`
}

func NewApplication(
	productName string,
	//authorID uuid.UUID,
	department string,
	amount float64,
) *Application {
	now := time.Now()

	return &Application{
		//ID:          1, //uuid.New(),
		ProductName: productName,
		//AuthorID:    authorID,
		Department: department,
		Amount:     amount,
		Status:     StatusNew,
		CreatedAt:  now,
		UpdatedAt:  now,
		Version:    1,
	}
}
