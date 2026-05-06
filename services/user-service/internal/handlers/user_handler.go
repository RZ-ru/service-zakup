package handlers

import (
	"user-service/internal/services"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *services.UserService
}

func NewHandler(s *services.UserService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Create(c, req.Email, req.Name)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, user)
}

func (h *Handler) GetUser(c *gin.Context) {
	userID := c.Param("id")

	user, err := h.service.GetByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}

	c.JSON(200, user)
}
