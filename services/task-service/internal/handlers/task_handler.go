// internal/handlers/task_handler.go
package handlers

import (
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

	task, err := h.service.Create(c, req.Title, req.Description, userID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *Handler) GetTask(c *gin.Context) {
	id := c.Param("id")

	task, err := h.service.GetByID(c, id)
	if err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}

	c.JSON(200, task)
}
