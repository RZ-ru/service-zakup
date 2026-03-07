package validation_service

import (
	"strings"

	"github.com/google/uuid"
)

type CreateApplicationInput struct {
	ProductID  uuid.UUID `json:"product_id"`
	Department string    `json:"department"`
	Amount     float64   `json:"amount"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationError struct {
	Errors []FieldError `json:"errors"`
}

func (e ValidationError) Error() string { return "validation error" }

func ValidateCreateApplication(in CreateApplicationInput) error {
	var errs []FieldError

	if in.ProductID == uuid.Nil {
		errs = append(errs, FieldError{Field: "product_id", Message: "must be a valid uuid"})
	}
	if strings.TrimSpace(in.Department) == "" {
		errs = append(errs, FieldError{Field: "department", Message: "must not be empty"})
	}
	if in.Amount <= 0 {
		errs = append(errs, FieldError{Field: "amount", Message: "must be greater than 0"})
	}

	if len(errs) > 0 {
		return ValidationError{Errors: errs}
	}
	return nil
}
