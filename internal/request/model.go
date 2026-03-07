package application // Описывает, что такое заявка и как она себя ведет.

import (
	"time"

	"github.com/google/uuid"
)

type Application struct {
	ID        uuid.UUID
	Number    string
	AuthorID  uuid.UUID
	ProductID uuid.UUID
	Quantity  int
	Comment   string
	Status    Status
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Создание новой заявки
func NewApplication(authorID, productID uuid.UUID, quantity int, comment string) (*Application, error) {
	if authorID == uuid.Nil {
		return nil, ErrInvalidAuthorID
	}

	if productID == uuid.Nil {
		return nil, ErrInvalidProductID
	}

	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}

	now := time.Now()

	return &Application{
		ID:        uuid.New(),
		AuthorID:  authorID,
		ProductID: productID,
		Quantity:  quantity,
		Comment:   comment,
		Status:    StatusDraft,
		Version:   1,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Проверка перехода статуса
func (a *Application) CanTransitionTo(next Status) bool {
	allowed, ok := allowedTransitions[a.Status]
	if !ok {
		return false
	}

	for _, s := range allowed {
		if s == next {
			return true
		}
	}

	return false
}

func (a *Application) ChangeStatus(next Status) error {
	if !a.CanTransitionTo(next) {
		return ErrInvalidStatusTransition
	}

	a.Status = next
	a.Version++
	a.UpdatedAt = time.Now()

	return nil
}
