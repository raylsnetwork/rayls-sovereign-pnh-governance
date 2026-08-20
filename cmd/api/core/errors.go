package core

import (
	"fmt"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/withstack"
)

// ============================================================================
// Domain Errors
// These errors express business-level failures
// HTTP adapters map these to appropriate status codes (404, 400, etc.)
// ============================================================================

// NotFoundError indicates a requested resource does not exist
type NotFoundError struct {
	Resource string // "transaction", "batch"
	ID       string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found: %s", e.Resource, e.ID)
}

// NewNotFoundError creates a NotFoundError
func NewNotFoundError(resource, id string) *NotFoundError {
	return &NotFoundError{
		Resource: resource,
		ID:       id,
	}
}

// ValidationError indicates invalid input parameters
type ValidationError struct {
	Field   string // The field that failed validation
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

// NewValidationError creates a ValidationError
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// InternalError wraps unexpected errors that should be logged and returned as 500
type InternalError struct {
	Operation string
	Err       error
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("internal error during %s: %v", e.Operation, e.Err)
}

func (e *InternalError) Unwrap() error {
	return e.Err
}

// NewInternalError creates an InternalError with stack trace
func NewInternalError(operation string, err error) *InternalError {
	return &InternalError{
		Operation: operation,
		Err:       withstack.Wrap(err),
	}
}

// ConflictError indicates a resource already exists (e.g. duplicate unique key)
type ConflictError struct {
	Resource string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("conflict: %s already exists", e.Resource)
}

// NewConflictError creates a ConflictError
func NewConflictError(resource string) *ConflictError {
	return &ConflictError{Resource: resource}
}

// ErrRecordConflict is a sentinel returned by repositories on unique constraint violations
var ErrRecordConflict = fmt.Errorf("record already exists")
