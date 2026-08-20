package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/core"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
)

// Ensure ParticipantRepositoryAdapter implements core.ParticipantRepository at compile time
var _ core.ParticipantRepository = (*ParticipantRepositoryAdapter)(nil)

// ParticipantRepositoryAdapter implements the ParticipantRepository interface for the core layer
type ParticipantRepositoryAdapter struct {
	db *gorm.DB
}

// NewParticipantRepository creates a new ParticipantRepositoryAdapter instance
func NewParticipantRepository(db *gorm.DB) core.ParticipantRepository {
	if db == nil {
		panic("db is nil")
	}

	return &ParticipantRepositoryAdapter{
		db: db,
	}
}

// GetByChainId implements core.ParticipantRepository
func (r *ParticipantRepositoryAdapter) GetByChainId(ctx context.Context, chainId string) (domain.Participant, error) {
	var participant domain.Participant
	if err := r.db.WithContext(ctx).Where("chain_id = ?", chainId).First(&participant).Error; err != nil {
		return domain.Participant{}, fmt.Errorf("failed to get participant by chain id %s: %w", chainId, err)
	}

	participant.StatusStr = participant.GetStatusDescription()
	participant.RoleStr = participant.GetRoleDescription()

	return participant, nil
}

// Upsert implements core.ParticipantRepository
func (r *ParticipantRepositoryAdapter) Upsert(ctx context.Context, participant domain.Participant) error {
	// Convert StatusStr and RoleStr to numeric values (like original implementation)
	if participant.StatusStr != "" {
		if status, exists := domain.StringToMemberStatus[participant.StatusStr]; exists {
			participant.Status = uint8(status) // #nosec G115 -- MemberStatus values are 0-3
		}
	}
	if participant.RoleStr != "" {
		if role, exists := domain.StringToMemberRole[participant.RoleStr]; exists {
			participant.Role = uint8(role) // #nosec G115 -- MemberRole values are 0-2
		}
	}

	// Use chain_id as the unique identifier for upserts
	var existing domain.Participant
	err := r.db.WithContext(ctx).Where("chain_id = ?", participant.ChainId).First(&existing).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check existing participant: %w", err)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new participant
		if err := r.db.WithContext(ctx).Create(&participant).Error; err != nil {
			return fmt.Errorf("failed to create participant: %w", err)
		}
	} else {
		// Update existing participant
		if participant.ID == uuid.Nil {
			participant.ID = existing.ID
		}
		if participant.CreatedAt.IsZero() {
			participant.CreatedAt = existing.CreatedAt
		}

		if err := r.db.WithContext(ctx).Save(&participant).Error; err != nil {
			return fmt.Errorf("failed to update participant: %w", err)
		}
	}

	return nil
}
