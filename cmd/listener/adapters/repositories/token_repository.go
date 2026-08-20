package repositories

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
)

// TokenRepositoryAdapter implements the TokenRepository interface for the core layer
type TokenRepositoryAdapter struct {
	db *gorm.DB
}

// NewTokenRepository creates a new TokenRepositoryAdapter instance
func NewTokenRepository(db *gorm.DB) *TokenRepositoryAdapter {
	return &TokenRepositoryAdapter{
		db: db,
	}
}

// GetTokenByResourceID retrieves a token by resource ID
func (t *TokenRepositoryAdapter) GetTokenByResourceID(ctx context.Context, resourceID string) (*domain.Token, error) {
	if resourceID == "" {
		return nil, fmt.Errorf("resourceID cannot be empty")
	}

	var domainToken domain.Token
	result := t.db.WithContext(ctx).
		Where("resource_id = ?", resourceID).
		First(&domainToken)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil //nolint:nilnil // nil signals record not found without wrapping a sentinel error
		}
		return nil, fmt.Errorf("failed to get token by resourceID: %w", result.Error)
	}

	return &domainToken, nil
}

func (t *TokenRepositoryAdapter) GetByIssuerAndName(
	ctx context.Context,
	issuerChainId string,
	name string,
) (*domain.Token, error) {
	var token domain.Token
	if err := t.db.WithContext(ctx).First(&token, "issuer_id = ? AND name = ?", issuerChainId, name).Error; err != nil {
		return nil, fmt.Errorf("failed to get token by issuer and name: %w", err)
	}

	return &token, nil
}

func (t *TokenRepositoryAdapter) Upsert(ctx context.Context, token *domain.Token) error {
	if token == nil {
		return fmt.Errorf("token cannot be nil")
	}
	if token.ResourceId == "" {
		return fmt.Errorf("token resource_id cannot be empty")
	}

	// Update all columns except ID and created_at
	err := t.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "resource_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"name",
			"symbol",
			"metadata_url",
			"erc_standard",
			"decimals",
			"issuer_id",
			"status",
			"updated_at",
		}),
	}).Create(token).Error
	if err != nil {
		return fmt.Errorf("failed to upsert token: %w", err)
	}

	return nil
}
