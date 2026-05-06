package services

import (
	"context"
	"database/sql"
	"errors"
	"permission-service/internal/cache"
	"permission-service/internal/models"
	"permission-service/internal/repository"
	"strconv"
	"time"
)

const permissionCacheTTL = time.Minute * 5

type PermissionService struct {
	repo  repository.PermissionRepository
	cache *cache.Redis
}

func NewPermissionService(r repository.PermissionRepository, c *cache.Redis) *PermissionService {
	return &PermissionService{
		repo:  r,
		cache: c,
	}
}

func (s *PermissionService) Create(userID, taskID string) error {
	err := s.repo.Create(&models.Permission{
		UserID: userID,
		TaskID: taskID,
		Role:   "owner",
	})
	if err != nil {
		return err
	}

	ctx := context.Background()
	key := permissionCacheKey(userID, taskID)

	_ = s.cache.Set(ctx, key, "true", permissionCacheTTL)

	return nil
}

func (s *PermissionService) Check(userID, taskID string) (bool, error) {
	ctx := context.Background()
	key := permissionCacheKey(userID, taskID)

	val, err := s.cache.Get(ctx, key)
	if err == nil {
		return val == "true", nil
	}

	role, err := s.repo.GetRole(userID, taskID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = s.cache.Set(ctx, key, "false", permissionCacheTTL)
			return false, nil
		}

		return false, err
	}

	allowed := role == "owner"

	_ = s.cache.Set(ctx, key, strconv.FormatBool(allowed), permissionCacheTTL)

	return allowed, nil
}

func permissionCacheKey(userID, taskID string) string {
	return "perm:" + userID + ":" + taskID
}
