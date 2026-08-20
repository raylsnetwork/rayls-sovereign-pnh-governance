package types

// Paginated wraps API responses with pagination metadata
type Paginated[T any] struct {
	Data  []T   `json:"data"`
	Total int64 `json:"total"`
	Limit int   `json:"limit"`
	Page  int   `json:"page"`
}
