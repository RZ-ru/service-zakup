package request // Доменная модель заявки

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// Application описывает сущность заявки в системе.
// Это доменная модель, которая хранит состояние заявки
// и содержит базовые правила изменения этого состояния.
type Application struct {
	ID uuid.UUID // Уникальный идентификатор заявки

	// Number string // Номер заявки для отображения пользователю.
	// Можно генерировать позже (например: REQ-2026-000001).
	// Пока закомментировано, потому что логика генерации номера не реализована.

	AuthorID  uuid.UUID // Пользователь, который создал заявку
	ProductID uuid.UUID // Товар, на который создается заявка

	Quantity int32  // Количество товара
	Comment  string // Комментарий к заявке

	Status Status // Текущий статус заявки

	Version int // Версия заявки (используется для контроля изменений)

	CreatedAt time.Time // Время создания заявки
	UpdatedAt time.Time // Время последнего изменения
}

// validate проверяет инварианты сущности.
func (a *Application) validate() error {
	if a == nil {
		return ErrNilApplication
	}

	if a.AuthorID == uuid.Nil {
		return ErrInvalidAuthorID
	}

	if a.ProductID == uuid.Nil {
		return ErrInvalidProductID
	}

	if a.Quantity <= 0 {
		return ErrInvalidQuantity
	}

	if !a.Status.Valid() {
		return ErrInvalidStatus
	}

	return nil
}

// NewApplication создаёт новую заявку.
func NewApplication(authorID, productID uuid.UUID, quantity int32, comment string) (*Application, error) {
	comment = strings.TrimSpace(comment)
	now := time.Now()

	app := &Application{
		ID:        uuid.New(),
		AuthorID:  authorID,
		ProductID: productID,
		Quantity:  quantity,
		Comment:   comment,
		Status:    StatusDraft,
		Version:   1,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := app.validate(); err != nil {
		return nil, err
	}

	return app, nil
}

// CanTransitionTo проверяет, можно ли изменить текущий статус
// заявки на следующий.
func (a *Application) CanTransitionTo(next Status) bool {

	if a == nil {
		return false
	}

	// Получаем список разрешенных переходов для текущего статуса
	allowed, ok := allowedTransitions[a.Status]
	if !ok {
		return false
	}

	// Проверяем есть ли нужный статус среди разрешенных
	for _, s := range allowed {
		if s == next {
			return true
		}
	}

	return false
}

// ChangeStatus изменяет статус заявки.
// При изменении статуса увеличивается версия и обновляется UpdatedAt.
func (a *Application) ChangeStatus(next Status) error {

	if a == nil {
		return ErrNilApplication
	}

	// Проверяем допустимость перехода
	if !a.CanTransitionTo(next) {
		return ErrInvalidStatusTransition
	}

	// Меняем статус
	a.Status = next

	// Увеличиваем версию
	a.Version++

	// Обновляем время изменения
	a.UpdatedAt = time.Now()

	return nil
}
