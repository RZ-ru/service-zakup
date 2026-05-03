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
		UserID string `json:"user_id"`
		TaskID string `json:"task_id"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.service.Create(req.UserID, req.TaskID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Status(201)
}

func (h *Handler) Check(c *gin.Context) {
	userID := c.Query("user_id")
	taskID := c.Query("task_id")

	ok, err := h.service.Check(userID, taskID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"allowed": ok})
}
