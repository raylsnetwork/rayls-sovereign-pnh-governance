package core

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/dto"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/logger"
)

const (
	// MaxPageSize is the maximum allowed page size for paginated queries
	MaxPageSize = 1000
)

var _ HeaderProofService = (*headerProofService)(nil)

// headerProofService implements the HeaderProofService interface
type headerProofService struct {
	repo HeaderProofRepository
	log  logger.Logger
}

// NewHeaderProofService creates a new header proof service
func NewHeaderProofService(repo HeaderProofRepository, log logger.Logger) HeaderProofService {
	return &headerProofService{
		repo: repo,
		log:  log,
	}
}

// GetHeaderProofsList retrieves header proofs for a block range with pagination
func (s *headerProofService) GetHeaderProofsList(
	ctx context.Context,
	filters dto.HeaderProofFilters,
) (*dto.HeaderProofListResponse, error) {
	// Set default pagination values
	if filters.Page == 0 {
		filters.Page = 1
	}
	if filters.PageSize == 0 {
		filters.PageSize = 50
	}

	// Validate pageSize limit
	if filters.PageSize > MaxPageSize {
		return nil, NewValidationError("pageSize", fmt.Sprintf("pageSize cannot exceed %d", MaxPageSize))
	}

	// Validate and convert startBlock
	startBlock, err := strconv.ParseInt(filters.StartBlock, 10, 64)
	if err != nil || startBlock < 0 {
		return nil, NewValidationError("startBlock", "startBlock must be a non-negative number")
	}

	// Validate and convert endBlock
	endBlock, err := strconv.ParseInt(filters.EndBlock, 10, 64)
	if err != nil || endBlock < 0 {
		return nil, NewValidationError("endBlock", "endBlock must be a non-negative number")
	}

	// Validate block range
	if endBlock < startBlock {
		return nil, NewValidationError("endBlock", "endBlock must be greater than or equal to startBlock")
	}

	// Fetch from repository
	headers, totalCount, err := s.repo.FindByBlockRange(
		ctx,
		filters.ChainID,
		startBlock,
		endBlock,
		filters.Page,
		filters.PageSize,
	)
	if err != nil {
		return nil, NewInternalError("FindByBlockRange", err)
	}

	// Convert domain models to response DTOs
	responseData := make([]dto.HeaderProofResponse, len(headers))
	for i, header := range headers {
		responseData[i] = dto.HeaderProofResponse{
			ID:          header.ID,
			ChainID:     header.ChainID.String(),
			BlockNumber: header.BlockNumber.String(),
			BlockHash:   header.BlockHash,
			CreatedAt:   header.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalCount) / float64(filters.PageSize)))

	// Build response with pagination metadata
	response := &dto.HeaderProofListResponse{
		Data: responseData,
		Pagination: dto.PaginationMetadata{
			CurrentPage:  filters.Page,
			PageSize:     filters.PageSize,
			TotalRecords: totalCount,
			TotalPages:   totalPages,
		},
	}

	return response, nil
}
