package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// Default pagination values
const (
	DefaultPage          = 1
	DefaultBatchPageSize = 10
	MaxPageSize          = 100
)

// PaginationParams holds validated pagination parameters
type PaginationParams struct {
	Page  int
	Limit int
}

// ParsePaginationParams extracts and validates page/limit from query params.
// Returns defaults if params are missing or invalid.
func ParsePaginationParams(c *gin.Context, defaultPage, defaultLimit int) PaginationParams {
	page := defaultPage
	limit := defaultLimit

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= MaxPageSize {
			limit = parsed
		}
	}

	return PaginationParams{Page: page, Limit: limit}
}
