package handlers

import (
	"auth-service/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *services.AuthService
}

func NewHandler(s *services.AuthService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Login(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.GenerateToken(req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
