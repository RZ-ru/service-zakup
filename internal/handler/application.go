package handler

import (
	"net/http"

	"zakup/internal/models"
	"zakup/internal/service"
	"zakup/validation_service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ApplicationHandler struct {
	svc *service.ApplicationService
}

func NewApplicationHandler(svc *service.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{svc: svc}
}

func (h *ApplicationHandler) PostApplications(c *gin.Context) {
	var in validation_service.CreateApplicationInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	// пока заглушка: authorID берём фиксированный
	authorID := uuid.New()

	app, err := h.svc.Create(c.Request.Context(), authorID, in)
	if err != nil {
		// красиво отдаём ошибки валидации
		if ve, ok := err.(validation_service.ValidationError); ok {
			c.JSON(http.StatusBadRequest, ve)
			return
		}
		if err == service.ErrProductNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "product_id not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, app)
}

type ChangeStatusInput struct {
	Status models.Status `json:"status"`
}

func (h *ApplicationHandler) PatchApplicationStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be uuid"})
		return
	}

	var in ChangeStatusInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	if !in.Status.Valid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	app, err := h.svc.ChangeStatus(c.Request.Context(), id, in.Status)
	if err != nil {
		if err == models.ErrInvalidStatusTransition {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status transition"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, app)
}
