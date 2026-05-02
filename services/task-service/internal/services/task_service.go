package services

import (
	"context"
	"errors"
	"time"

	"task-service/internal/models"
	"task-service/internal/repository"

	"github.com/google/uuid"
)

type TaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(r repository.TaskRepository) *TaskService {
	return &TaskService{repo: r}
}

func (s *TaskService) Create(ctx context.Context, title, description, userID string) (*models.Task, error) {

	if title == "" {
		return nil, errors.New("title required")
	}

	task := &models.Task{
		ID:          uuid.NewString(),
		Title:       title,
		Description: description,
		Status:      "new",
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) GetByID(ctx context.Context, id string) (*models.Task, error) {
	return s.repo.GetByID(ctx, id)
}
