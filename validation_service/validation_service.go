package validation_service

import (
	"github.com/google/uuid"
)

type CreateApplicationInput struct {
	ProductID uuid.UUID `json:"product_id"`
	Comment   string    `json:"comment"`
	Quantity  int32     `json:"quantity"`
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
	// if strings.TrimSpace(in.Comment) == "" {
	// 	errs = append(errs, FieldError{Field: "department", Message: "must not be empty"})
	// }
	if in.Quantity <= 0 {
		errs = append(errs, FieldError{Field: "quantity", Message: "must be greater than 0"})
	}

	if len(errs) > 0 {
		return ValidationError{Errors: errs}
	}
	return nil
}
