package application // Общие ошибки домена

import "errors"

var (
	ErrInvalidAuthorID         = errors.New("invalid author id")
	ErrInvalidProductID        = errors.New("invalid product id")
	ErrInvalidQuantity         = errors.New("invalid quantity")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrApplicationNotFound     = errors.New("application not found")
)
