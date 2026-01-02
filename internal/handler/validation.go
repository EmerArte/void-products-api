package handler

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidationError represents a field-level validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrorResponse represents the response for validation errors
type ValidationErrorResponse struct {
	Success bool              `json:"success"`
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Details []ValidationError `json:"details,omitempty"`
}

// FormatValidationErrors converts validator errors into user-friendly messages
func FormatValidationErrors(err error) (string, []ValidationError) {
	var details []ValidationError

	// Check if it's a validation error
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		// Not a validation error, return generic message
		return err.Error(), nil
	}

	// Process each validation error
	for _, e := range validationErrors {
		field := formatFieldName(e.Field())
		message := getValidationMessage(e)

		details = append(details, ValidationError{
			Field:   field,
			Message: message,
		})
	}

	// Create summary message
	if len(details) == 1 {
		return fmt.Sprintf("Validation failed for field '%s': %s", details[0].Field, details[0].Message), details
	}

	return fmt.Sprintf("Validation failed for %d field(s)", len(details)), details
}

// formatFieldName converts field name to snake_case for consistency
func formatFieldName(field string) string {
	// Convert from PascalCase to snake_case
	var result strings.Builder
	for i, r := range field {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// getValidationMessage returns a user-friendly message based on the validation tag
func getValidationMessage(e validator.FieldError) string {
	field := formatFieldName(e.Field())
	tag := e.Tag()
	param := e.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("'%s' is required", field)

	case "email":
		return fmt.Sprintf("'%s' must be a valid email address", field)

	case "min":
		switch e.Kind().String() {
		case "string":
			return fmt.Sprintf("'%s' must be at least %s characters long", field, param)
		case "slice", "array":
			return fmt.Sprintf("'%s' must contain at least %s items", field, param)
		default:
			return fmt.Sprintf("'%s' must be at least %s", field, param)
		}

	case "max":
		switch e.Kind().String() {
		case "string":
			return fmt.Sprintf("'%s' must be at most %s characters long", field, param)
		case "slice", "array":
			return fmt.Sprintf("'%s' must contain at most %s items", field, param)
		default:
			return fmt.Sprintf("'%s' must be at most %s", field, param)
		}

	case "gte":
		return fmt.Sprintf("'%s' must be greater than or equal to %s", field, param)

	case "gt":
		return fmt.Sprintf("'%s' must be greater than %s", field, param)

	case "lte":
		return fmt.Sprintf("'%s' must be less than or equal to %s", field, param)

	case "lt":
		return fmt.Sprintf("'%s' must be less than %s", field, param)

	case "oneof":
		validValues := strings.ReplaceAll(param, " ", ", ")
		return fmt.Sprintf("'%s' must be one of: %s", field, validValues)

	case "url":
		return fmt.Sprintf("'%s' must be a valid URL", field)

	case "uuid":
		return fmt.Sprintf("'%s' must be a valid UUID", field)

	case "uuid4":
		return fmt.Sprintf("'%s' must be a valid UUID v4", field)

	case "len":
		return fmt.Sprintf("'%s' must be exactly %s characters long", field, param)

	case "eqfield":
		return fmt.Sprintf("'%s' must be equal to '%s'", field, formatFieldName(param))

	case "nefield":
		return fmt.Sprintf("'%s' must not be equal to '%s'", field, formatFieldName(param))

	case "numeric":
		return fmt.Sprintf("'%s' must be a numeric value", field)

	case "alpha":
		return fmt.Sprintf("'%s' must contain only alphabetic characters", field)

	case "alphanum":
		return fmt.Sprintf("'%s' must contain only alphanumeric characters", field)

	case "startswith":
		return fmt.Sprintf("'%s' must start with '%s'", field, param)

	case "endswith":
		return fmt.Sprintf("'%s' must end with '%s'", field, param)

	case "contains":
		return fmt.Sprintf("'%s' must contain '%s'", field, param)

	case "containsany":
		return fmt.Sprintf("'%s' must contain at least one of these characters: %s", field, param)

	case "excludes":
		return fmt.Sprintf("'%s' must not contain '%s'", field, param)

	case "dive":
		// Dive is for nested validation, the actual error will be in nested fields
		return fmt.Sprintf("'%s' contains invalid nested data", field)

	default:
		// Generic message for unknown validation tags
		if param != "" {
			return fmt.Sprintf("'%s' failed validation '%s' with parameter '%s'", field, tag, param)
		}
		return fmt.Sprintf("'%s' failed validation '%s'", field, tag)
	}
}
