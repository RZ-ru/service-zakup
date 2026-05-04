package handlers

import (
	"permission-service/internal/services"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *services.PermissionService
}

func NewHandler(s *services.PermissionService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Create(c *gin.Context) {
	var req struct {
		TaskID string `json:"task_id"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if req.TaskID == "" {
		c.JSON(400, gin.H{"error": "task_id required"})
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	err := h.service.Create(userID, req.TaskID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Status(201)
}

func (h *Handler) Check(c *gin.Context) {
	userID := c.GetString("user_id")
	taskID := c.Query("task_id")

	if userID == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	if taskID == "" {
		c.JSON(400, gin.H{"error": "task_id required"})
		return
	}

	ok, err := h.service.Check(userID, taskID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"allowed": ok})
}
