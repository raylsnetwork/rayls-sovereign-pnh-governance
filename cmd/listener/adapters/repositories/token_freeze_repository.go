package repositories

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/core"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/withstack"
)

// TokenFreezeRepositoryAdapter implements token freeze operations
type TokenFreezeRepositoryAdapter struct {
	db *gorm.DB
}

// NewTokenFreezeRepository creates a new TokenFreezeRepositoryAdapter instance
func NewTokenFreezeRepository(db *gorm.DB) core.TokenFreezeRepository {
	return &TokenFreezeRepositoryAdapter{
		db: db,
	}
}

// UpdateTokenFreezeStatus updates both state and audit tables for frozen tokens
func (r *TokenFreezeRepositoryAdapter) UpdateTokenFreezeStatus(
	ctx context.Context,
	resourceId string,
	chainIds []*big.Int,
	action uint8,
	blockNumber *big.Int,
	txHash string,
	blockTimestamp time.Time,
) error {
	if resourceId == "" {
		return fmt.Errorf("resourceId cannot be empty")
	}
	if len(chainIds) == 0 {
		return fmt.Errorf("chainIds cannot be empty")
	}
	if blockNumber == nil {
		return fmt.Errorf("blockNumber cannot be nil")
	}
	if txHash == "" {
		return fmt.Errorf("txHash cannot be empty")
	}

	// action: 0 = UNFREEZE, 1 = FREEZE
	isFrozen := action == 1

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Normalize chain identifiers: create one database row per chainId
		for _, chainId := range chainIds {
			chainIdStr := chainId.String()

			// 1. Insert audit record with idempotency protection
			auditRecord := domain.TokenFreezeAudit{
				ID:              uuid.New(),
				ResourceId:      resourceId,
				ChainId:         chainIdStr,
				Action:          action,
				BlockNumber:     domain.NewBigInt(blockNumber),
				TransactionHash: txHash,
				CreatedAt:       blockTimestamp.UTC(),
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{
					{Name: "resource_id"},
					{Name: "chain_id"},
					{Name: "transaction_hash"},
					{Name: "action"},
				},
				DoNothing: true,
			}).Create(&auditRecord).Error; err != nil {
				return withstack.Wrap(err)
			}

			// 2. Upsert state table with idempotency protection
			stateRecord := domain.TokenFreezeState{
				ResourceId: resourceId,
				ChainId:    chainIdStr,
				IsFrozen:   isFrozen,
				CreatedAt:  blockTimestamp.UTC(),
				UpdatedAt:  blockTimestamp.UTC(),
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "resource_id"}, {Name: "chain_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"is_frozen", "updated_at"}),
			}).Create(&stateRecord).Error; err != nil {
				return withstack.Wrap(err)
			}
		}

		return nil
	})
}
