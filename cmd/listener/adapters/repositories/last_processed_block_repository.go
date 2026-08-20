package repositories

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/core"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
)

// Ensure LastProcessedBlockRepository implements core.BlockRepository at compile time
var _ core.BlockRepository = (*LastProcessedBlockRepository)(nil)

// LastProcessedBlockRepository implements core.BlockRepository
type LastProcessedBlockRepository struct {
	db *gorm.DB
}

// NewLastProcessedBlockRepository creates a new last processed block repository adapter
func NewLastProcessedBlockRepository(dbClient *gorm.DB) core.BlockRepository {
	if dbClient == nil {
		panic("dbClient is nil")
	}

	return &LastProcessedBlockRepository{
		db: dbClient,
	}
}

// GetLatestProcessedBlock implements core.BlockRepository
func (r *LastProcessedBlockRepository) GetLatestProcessedBlock(ctx context.Context) (*big.Int, error) {
	var block domain.LastProcessedBlock
	if err := r.db.WithContext(ctx).First(&block).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no last processed block found in database: %w", err)
		}
		return nil, fmt.Errorf("failed to query last processed block: %w", err)
	}

	return block.Number.Unwrap(), nil
}

// UpdateLatestProcessedBlock implements core.BlockRepository
func (r *LastProcessedBlockRepository) UpdateLatestProcessedBlock(ctx context.Context, blockNumber *big.Int) error {
	block := domain.LastProcessedBlock{
		ID:     1,
		Number: domain.NewBigInt(blockNumber),
	}

	// Use Create with OnConflict to handle upsert properly
	result := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"number", "updated_at"}),
		}).
		Create(&block)

	if result.Error != nil {
		return fmt.Errorf("failed to save last processed block %s: %w", blockNumber.String(), result.Error)
	}

	return nil
}
