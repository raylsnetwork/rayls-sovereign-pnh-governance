package utils

import "fmt"

// ErrorResponse creates a standardized error response with error and hint fields
type ErrorResponse struct {
	Error string `json:"error"`
	Hint  string `json:"hint,omitempty"`
}

// NewUnknownParamsError creates an error response for unknown query parameters
func NewUnknownParamsError(unknownParams string, allowedParams string) ErrorResponse {
	return ErrorResponse{
		Error: fmt.Sprintf("Unknown query parameter(s): %s", unknownParams),
		Hint:  fmt.Sprintf("Allowed parameters: %s", allowedParams),
	}
}

// NewEmptyParamsError creates an error response for empty query parameter values
func NewEmptyParamsError(emptyParams string) ErrorResponse {
	return ErrorResponse{
		Error: fmt.Sprintf("Empty value for query parameter(s): %s", emptyParams),
		Hint:  "Query parameters must have a non-empty value or be omitted entirely",
	}
}

// NewInvalidEnumValueError creates an error response for invalid enum parameter values
func NewInvalidEnumValueError(param string, value string, allowedValues string) ErrorResponse {
	return ErrorResponse{
		Error: fmt.Sprintf("Parameter '%s' has an invalid value: %s", param, value),
		Hint:  fmt.Sprintf("Allowed values: %s", allowedValues),
	}
}

// NewMalformedQueryError creates an error response for malformed URL query encoding (e.g. %)
func NewMalformedQueryError(details string) ErrorResponse {
	return ErrorResponse{
		Error: "Malformed query parameter encoding",
		Hint:  details,
	}
}
