package services

import (
	"context"
	"errors"
	"time"

	"task-service/internal/clients"
	"task-service/internal/models"
	"task-service/internal/repository"

	"github.com/google/uuid"
)

type TaskService struct {
	repo repository.TaskRepository
	perm *clients.PermissionClient
}

func NewTaskService(repo repository.TaskRepository, perm *clients.PermissionClient) *TaskService {
	return &TaskService{
		repo: repo,
		perm: perm,
	}
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

	// Выдаем права, связываем пользователя и задачу
	err := s.perm.Create(userID, task.ID)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) GetByID(ctx context.Context, userID, taskID string) (*models.Task, error) {

	allowed, err := s.perm.Check(userID, taskID)
	if err != nil {
		return nil, err
	}

	if !allowed {
		return nil, errors.New("forbidden")
	}

	return s.repo.GetByID(ctx, taskID)
}
