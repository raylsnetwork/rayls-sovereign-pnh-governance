package dto

// HeaderProofFilters defines query parameters for filtering header proofs
type HeaderProofFilters struct {
	ChainID    string `form:"chainId"    binding:"required"                 example:"11155111"`
	StartBlock string `form:"startBlock" binding:"required,numeric"         example:"1"`
	EndBlock   string `form:"endBlock"   binding:"required,numeric"         example:"100"`
	Page       int    `form:"page"       binding:"omitempty,min=1"          example:"1"`
	PageSize   int    `form:"pageSize"   binding:"omitempty,min=1,max=1000" example:"50"`
}

// HeaderProofResponse represents a single header proof in the API response
type HeaderProofResponse struct {
	ID          int    `json:"id"          example:"1"`
	ChainID     string `json:"chainId"     example:"1337"`
	BlockNumber string `json:"blockNumber" example:"100"`
	BlockHash   string `json:"blockHash"   example:"0x4de2a17c6f8b9e3c7a5d2f1e6b9c4a8d3e7f1c5a9b2d8e6f0c4b7a1d9e3f8c5b"`

	CreatedAt string `json:"createdAt" example:"2024-10-31T10:30:00Z"`
}

// HeaderProofListResponse represents a paginated list of header proofs
type HeaderProofListResponse struct {
	Data       []HeaderProofResponse `json:"data"`
	Pagination PaginationMetadata    `json:"pagination"`
}

// PaginationMetadata contains pagination information
type PaginationMetadata struct {
	CurrentPage  int   `json:"currentPage"`
	PageSize     int   `json:"pageSize"`
	TotalRecords int64 `json:"totalRecords"`
	TotalPages   int   `json:"totalPages"`
}
