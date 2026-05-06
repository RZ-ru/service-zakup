package services

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"task-service/internal/clients"
	"task-service/internal/models"
	"task-service/internal/repository"

	"github.com/google/uuid"
)

type TaskService struct {
	repo   repository.TaskRepository
	outbox repository.OutboxRepository
	perm   *clients.PermissionClient
}

func NewTaskService(repo repository.TaskRepository, outbox repository.OutboxRepository, perm *clients.PermissionClient) *TaskService {
	return &TaskService{
		repo:   repo,
		outbox: outbox,
		perm:   perm,
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

	payload, err := json.Marshal(map[string]string{
		"event_type": "task.created",
		"task_id":    task.ID,
		"user_id":    task.UserID,
		"title":      task.Title,
		"status":     task.Status,
	})
	if err != nil {
		return nil, err
	}

	event := &models.OutboxEvent{
		ID:            uuid.NewString(),
		AggregateType: "task",
		AggregateID:   task.ID,
		EventType:     "task.created",
		RoutingKey:    "task.created",
		Payload:       payload,
		Status:        "pending",
		Attempts:      0,
		CreatedAt:     time.Now(),
		ProcessedAt:   nil,
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	if err := s.repo.CreateTx(ctx, tx, task); err != nil {
		return nil, err
	}

	if err := s.outbox.CreateTx(ctx, tx, event); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	committed = true

	// Выдаем права, связываем пользователя и задачу
	err = s.perm.Create(ctx, task.ID)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) GetByID(ctx context.Context, userID, taskID, role string) (*models.Task, error) {

	if role == "admin" {
		return s.repo.GetByID(ctx, taskID)
	}

	allowed, err := s.perm.Check(ctx, taskID)
	if err != nil {
		return nil, err
	}

	if !allowed {
		return nil, errors.New("forbidden")
	}

	return s.repo.GetByID(ctx, taskID)
}
