package request // Собирает бизнес-сценарии

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

type CreateInput struct {
	AuthorID  uuid.UUID
	ProductID uuid.UUID
	Quantity  int32
	Comment   string
}

// Создание заявки
func (s *Service) Create(ctx context.Context, in CreateInput) (*Application, error) {
	app, err := NewApplication(in.AuthorID, in.ProductID, in.Quantity, in.Comment)
	if err != nil {
		return nil, err
	}

	eventPayload, err := json.Marshal(ApplicationCreatedEvent{
		ApplicationID: app.ID,
		AuthorID:      app.AuthorID,
		ProductID:     app.ProductID,
		Quantity:      app.Quantity,
		Status:        string(app.Status),
		CreatedAt:     app.CreatedAt,
	})
	if err != nil {
		return nil, err
	}

	outboxEvent := NewOutboxEvent(
		"application",
		app.ID,
		"application.created",
		"application.created",
		eventPayload,
	)

	if err := s.repo.CreateWithOutbox(ctx, app, outboxEvent); err != nil {
		return nil, err
	}

	return app, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Application, error) { //Получить заявку по ID
	if id == uuid.Nil {
		return nil, ErrInvalidApplicationID
	}

	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]*Application, error) { //Список заявок
	return s.repo.List(ctx)
}

type ChangeStatusInput struct {
	ID     uuid.UUID
	Status Status
}

func (s *Service) ChangeStatus(ctx context.Context, in ChangeStatusInput) (*Application, error) {
	if in.ID == uuid.Nil {
		return nil, ErrInvalidApplicationID
	}

	if !in.Status.Valid() {
		return nil, ErrInvalidStatus
	}

	app, err := s.repo.GetByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	oldStatus := app.Status

	if err := app.ChangeStatus(in.Status); err != nil {
		return nil, err
	}

	eventPayload, err := json.Marshal(ApplicationStatusChangedEvent{
		ApplicationID: app.ID,
		OldStatus:     string(oldStatus),
		NewStatus:     string(app.Status),
		ChangedAt:     app.UpdatedAt,
		Version:       app.Version,
	})
	if err != nil {
		return nil, err
	}

	outboxEvent := NewOutboxEvent(
		"application",
		app.ID,
		"application.status_changed",
		"application.status_changed",
		eventPayload,
	)

	if err := s.repo.UpdateWithOutbox(ctx, app, outboxEvent); err != nil {
		return nil, err
	}

	return app, nil
}
