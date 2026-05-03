package services

import (
	"permission-service/internal/models"
	"permission-service/internal/repository"
)

type PermissionService struct {
	repo repository.PermissionRepository
}

func NewPermissionService(r repository.PermissionRepository) *PermissionService {
	return &PermissionService{repo: r}
}

func (s *PermissionService) Create(userID, taskID string) error {
	return s.repo.Create(&models.Permission{
		UserID: userID,
		TaskID: taskID,
		Role:   "owner",
	})
}

func (s *PermissionService) Check(userID, taskID string) (bool, error) {
	return s.repo.Exists(userID, taskID)
}
