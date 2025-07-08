package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/afreedicp/zolaris-backend-app/internal/transport/dto"
)

// ValidationErrorItem represents a validation error for a specific field
type ValidationErrorItem struct {
	Field string
	Error string
}

// Validate validates a struct using validator tags
func Validate(s any) []ValidationErrorItem {
	validate := validator.New()
	err := validate.Struct(s)

	if err == nil {
		return nil
	}

	var validationErrors []ValidationErrorItem
	validationErrs := err.(validator.ValidationErrors)

	for _, e := range validationErrs {
		var err ValidationErrorItem
		err.Field = e.Field()
		err.Error = formatValidationError(e)
		validationErrors = append(validationErrors, err)
	}

	return validationErrors
}

// ValidationErrorsToString converts validation errors to a string
func ValidationErrorsToString(errors []ValidationErrorItem) string {
	if len(errors) == 0 {
		return ""
	}

	messages := make([]string, len(errors))
	for i, err := range errors {
		messages[i] = fmt.Sprintf("%s: %s", err.Field, err.Error)
	}

	return strings.Join(messages, "; ")
}

// CreateValidationError returns a formatted error message for validation errors
func CreateValidationError(errs []ValidationErrorItem) error {
	if len(errs) == 0 {
		return nil
	}
	return errors.New(ValidationErrorsToString(errs))
}

// CreateDtoValidationErrors converts internal validation errors to DTO validation errors
func CreateDtoValidationErrors(errs []ValidationErrorItem) []dto.ValidationError {
	result := make([]dto.ValidationError, len(errs))
	for i, err := range errs {
		result[i] = dto.ValidationError{
			Field:   err.Field,
			Message: err.Error,
		}
	}
	return result
}

// formatValidationError formats a validation error into a readable message
func formatValidationError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "required field"
	case "min":
		return fmt.Sprintf("must be at least %s characters long", err.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters long", err.Param())
	case "email":
		return "must be a valid email address"
	case "oneof":
		return fmt.Sprintf("must be one of: %s", err.Param())
	}
	return fmt.Sprintf("failed validation for '%s'", err.Tag())
}
