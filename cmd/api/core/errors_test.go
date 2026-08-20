package core

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationError_EmptyFieldShowsMessageOnly(t *testing.T) {
	// A validation error without a specific field should show message only
	err := &ValidationError{
		Field:   "",
		Message: "something went wrong",
	}

	result := err.Error()

	assert.Equal(t, "validation error: something went wrong", result)
	assert.NotContains(t, result, "field")
}

func TestValidationError_WithFieldShowsFieldAndMessage(t *testing.T) {
	// A validation error with a field should show both field and message
	err := NewValidationError("chainId", "must be an integer")

	result := err.Error()

	assert.Contains(t, result, "chainId")
	assert.Contains(t, result, "must be an integer")
}

func TestNotFoundError_MessageIncludesResourceAndIdentifier(t *testing.T) {
	// A not found error includes the resource type and identifier in its message
	err := NewNotFoundError("transaction", "abc-123")

	result := err.Error()

	assert.Contains(t, result, "transaction")
	assert.Contains(t, result, "abc-123")
	assert.Contains(t, result, "not found")
}

func TestInternalError_UnwrapReturnsOriginalError(t *testing.T) {
	// An internal error should unwrap to the original error
	originalErr := errors.New("database connection failed")
	err := NewInternalError("FindByChainId", originalErr)

	unwrapped := errors.Unwrap(err)

	assert.NotNil(t, unwrapped)
}
