package handler

import (
	"net/http"

	"zakup/internal/request"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ApplicationHandler struct {
	svc *request.Service
}

func NewApplicationHandler(svc *request.Service) *ApplicationHandler {
	return &ApplicationHandler{svc: svc}
}

func (h *ApplicationHandler) Register(r *gin.Engine) {
	// Создать заявку
	r.POST("/applications", h.PostApplications)
	// Изменить заявку
	r.PATCH("/applications/:id/status", h.PatchApplicationStatus)
	// Удалить заявку по id
	//r.DELETE("/applications/:id", h.DeleteApplication)
	// Получить заявки по id пользователя
	//r.GET("/applications/:id", h.GetApplication)
}

func (h *ApplicationHandler) PostApplications(c *gin.Context) {
	var req CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid json"})
		return
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ValidationErrorResponse{
			Error: "validation failed",
			Fields: map[string]string{
				"product_id": "must be a valid UUID",
			},
		})
		return
	}

	// пока заглушка: authorID берём фиксированный
	//authorID := uuid.New()
	authorID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	app, err := h.svc.Create(c.Request.Context(), request.CreateInput{
		AuthorID:  authorID,
		ProductID: productID,
		Quantity:  req.Quantity,
		Comment:   req.Comment,
	})
	if err != nil {
		if err == request.ErrInvalidProductID {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid product_id"})
			return
		}
		if err == request.ErrInvalidQuantity {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid quantity"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, toApplicationResponse(app))
}

func (h *ApplicationHandler) PatchApplicationStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id must be uuid"})
		return
	}

	var req ChangeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid json"})
		return
	}

	if !req.Status.Valid() {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid status"})
		return
	}

	app, err := h.svc.ChangeStatus(c.Request.Context(), request.ChangeStatusInput{
		ID:     id,
		Status: req.Status,
	})
	if err != nil {
		if err == request.ErrInvalidStatusTransition {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid status transition"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, toApplicationResponse(app))
}
