// internal/handlers/task_handler.go
package handlers

import (
	"context"
	"net/http"
	"task-service/internal/services"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *services.TaskService
}

func NewHandler(s *services.TaskService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CreateTask(c *gin.Context) {

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, gin.H{"error": "missing Authorization header"})
		return
	}

	ctx := context.WithValue(c.Request.Context(), "auth_header", authHeader)

	task, err := h.service.Create(ctx, req.Title, req.Description, userID)
	if err != nil {
		if err.Error() == "forbidden" {
			c.JSON(403, gin.H{"error": "forbidden"})
			return
		}

		c.JSON(404, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *Handler) GetTask(c *gin.Context) {
	taskID := c.Param("id")
	userID := c.GetString("user_id")
	role := c.GetString("role")

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, gin.H{"error": "missing Authorization header"})
		return
	}
	ctx := context.WithValue(c.Request.Context(), "auth_header", authHeader)

	task, err := h.service.GetByID(ctx, userID, taskID, role)
	if err != nil {
		if err.Error() == "forbidden" {
			c.JSON(403, gin.H{"error": "forbidden"})
			return
		}

		c.JSON(404, gin.H{"error": "not found"})
		return
	}

	c.JSON(200, task)
}
