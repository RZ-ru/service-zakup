package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusNew       Status = "new"
	StatusInProcess Status = "in_progress"
	StatusApproved  Status = "approved"
	StatusRejected  Status = "rejected"
	StatusCanceled  Status = "canceled"
)

func (s Status) Valid() bool {
	switch s {
	case StatusNew, StatusInProcess, StatusApproved, StatusRejected, StatusCanceled:
		return true
	default:
		return false
	}
}

type Application struct {
	ID         uuid.UUID `json:"id"`
	ProductID  uuid.UUID `json:"product_id"`
	AuthorID   uuid.UUID `json:"author_id"`
	Department string    `json:"department"`
	Amount     float64   `json:"amount"`
	Status     Status    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Version    int       `json:"version"`
}

// Создание доменной сущности — сервер сам задаёт ID/Status/время/version
func NewApplication(productID, authorID uuid.UUID, department string, amount float64) *Application {
	now := time.Now().UTC()

	return &Application{
		ID:         uuid.New(),
		ProductID:  productID,
		AuthorID:   authorID,
		Department: department,
		Amount:     amount,
		Status:     StatusNew,
		CreatedAt:  now,
		UpdatedAt:  now,
		Version:    1,
	}
}

var ErrInvalidStatusTransition = errors.New("invalid status transition")

func (a *Application) CanTransitionTo(next Status) bool {
	if !next.Valid() {
		return false
	}

	switch a.Status {
	case StatusNew:
		return next == StatusInProcess || next == StatusCanceled
	case StatusInProcess:
		return next == StatusApproved || next == StatusRejected || next == StatusCanceled
	case StatusApproved, StatusRejected, StatusCanceled:
		return false // финальные
	default:
		return false
	}
}

func (a *Application) SetStatus(next Status) error {
	if !a.CanTransitionTo(next) {
		return ErrInvalidStatusTransition
	}
	a.Status = next
	a.UpdatedAt = time.Now().UTC()
	a.Version++
	return nil
}
