package repositories

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/api/core"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
)

var _ core.HeaderProofRepository = (*headerProofRepository)(nil)

// headerProofRepository implements core.HeaderProofRepository using GORM
type headerProofRepository struct {
	db *gorm.DB
}

// NewHeaderProofRepository creates a new GORM-based header proof repository
func NewHeaderProofRepository(db *gorm.DB) core.HeaderProofRepository {
	return &headerProofRepository{db: db}
}

// FindByBlockRange retrieves header proofs within a block range with pagination
func (r *headerProofRepository) FindByBlockRange(
	ctx context.Context,
	chainId string,
	startBlock, endBlock int64,
	page, pageSize int,
) ([]domain.HeaderProofEvent, int64, error) {
	var headers []domain.HeaderProofEvent
	var totalCount int64

	// Build base query
	query := r.db.WithContext(ctx).Model(&domain.HeaderProofEvent{}).
		Where("chain_id = ?", chainId).
		Where("block_number >= ?", startBlock).
		Where("block_number <= ?", endBlock)

	// Get total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count header proofs: %w", err)
	}

	// If no records found, return empty result
	if totalCount == 0 {
		return []domain.HeaderProofEvent{}, 0, nil
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Fetch paginated results ordered by block_number ascending
	if err := query.
		Order("block_number ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&headers).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch header proofs: %w", err)
	}

	return headers, totalCount, nil
}
