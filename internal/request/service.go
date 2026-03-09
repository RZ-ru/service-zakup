package application // Собирает бизнес-сценарии

import (
	"context"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

type CreateInput struct {
	AuthorID  uuid.UUID
	ProductID uuid.UUID
	Quantity  int
	Comment   string
}

func (s *Service) Create(ctx context.Context, in CreateInput) (*Application, error) { //Создание заявки
	app, err := NewApplication(in.AuthorID, in.ProductID, in.Quantity, in.Comment)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, app); err != nil {
		return nil, err
	}

	return app, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Application, error) { //Получить заявку по ID
	if id == uuid.Nil {
		return nil, ErrApplicationNotFound
	}

	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]*Application, error) { //Список заявок
	return s.repo.List(ctx)
}

type ChangeStatusInput struct {
	ID     uuid.UUID
	Status Status
	// new upadte
}

func (s *Service) ChangeStatus(ctx context.Context, in ChangeStatusInput) (*Application, error) { //Смена статуса
	app, err := s.repo.GetByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	if err := app.ChangeStatus(in.Status); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, app); err != nil {
		return nil, err
	}

	return app, nil
}
