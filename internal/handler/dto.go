package handler

import (
	"time"

	"zakup/internal/request"
)

// CreateApplicationRequest — входной JSON на создание заявки.
type CreateApplicationRequest struct {
	ProductID string `json:"product_id" binding:"required,uuid"`
	Quantity  int32  `json:"quantity" binding:"required,min=1"`
	Comment   string `json:"comment" binding:"max=1000"`
}

// ChangeStatusRequest — входной JSON на смену статуса.
type ChangeStatusRequest struct {
	Status request.Status `json:"status" binding:"required"`
}

// ApplicationResponse — ответ клиенту по заявке.
type ApplicationResponse struct {
	ID        string         `json:"id"`
	AuthorID  string         `json:"author_id"`
	ProductID string         `json:"product_id"`
	Quantity  int32          `json:"quantity"`
	Comment   string         `json:"comment"`
	Status    request.Status `json:"status"`
	Version   int            `json:"version"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// ErrorResponse — обычная ошибка.
type ErrorResponse struct {
	Error string `json:"error"`
}

// ValidationErrorResponse — ошибки по полям.
type ValidationErrorResponse struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

// toApplicationResponse преобразует доменную модель в DTO ответа.
func toApplicationResponse(app *request.Application) ApplicationResponse {
	return ApplicationResponse{
		ID:        app.ID.String(),
		AuthorID:  app.AuthorID.String(),
		ProductID: app.ProductID.String(),
		Quantity:  app.Quantity,
		Comment:   app.Comment,
		Status:    app.Status,
		Version:   app.Version,
		CreatedAt: app.CreatedAt,
		UpdatedAt: app.UpdatedAt,
	}
}

// toApplicationListResponse преобразует список заявок в DTO.
func toApplicationListResponse(apps []*request.Application) []ApplicationResponse {
	out := make([]ApplicationResponse, 0, len(apps))
	for _, app := range apps {
		out = append(out, toApplicationResponse(app))
	}
	return out
}
