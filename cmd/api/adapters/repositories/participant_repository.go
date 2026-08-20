package repositories

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/api/core"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/api/utils"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/domain"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/dto"
)

var _ core.ParticipantRepository = (*participantRepository)(nil)

// participantRepository implements core.ParticipantRepository using GORM
type participantRepository struct {
	db *gorm.DB
}

// NewParticipantRepository creates a new GORM-based participant repository
func NewParticipantRepository(db *gorm.DB) core.ParticipantRepository {
	return &participantRepository{db: db}
}

// FindByChainId retrieves a single participant by chain ID
func (r *participantRepository) FindByChainId(ctx context.Context, chainId string) (*domain.Participant, error) {
	var participant domain.Participant

	err := r.db.WithContext(ctx).
		Where("chain_id = ?", chainId).
		First(&participant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrRecordNotFound
		}
		return nil, err
	}

	participant.StatusStr = participant.GetStatusDescription()
	participant.RoleStr = participant.GetRoleDescription()

	return &participant, nil
}

// FindByFilters finds participants matching the provided filters
func (r *participantRepository) FindByFilters(
	ctx context.Context,
	filters dto.ParticipantListFilters,
) ([]domain.Participant, error) {
	var participants []domain.Participant

	query := r.db.WithContext(ctx).Model(&domain.Participant{})

	// Apply filters
	if filters.Name != "" {
		query = query.Where("name = ?", filters.Name)
	}

	if filters.ChainId != nil {
		query = query.Where("chain_id = ?", *filters.ChainId)
	}

	if filters.Status != "" {
		if statusInt, exists := domain.StringToMemberStatus[filters.Status]; exists {
			query = query.Where("status = ?", statusInt)
		}
	}

	if filters.Role != "" {
		if roleInt, exists := domain.StringToMemberRole[filters.Role]; exists {
			query = query.Where("role = ?", roleInt)
		}
	}

	if filters.CreatedAfter != "" {
		afterTime, _ := utils.ParseTime(filters.CreatedAfter)
		query = query.Where("created_at >= ?", afterTime)
	}

	if filters.CreatedBefore != "" {
		beforeTime, _ := utils.ParseTime(filters.CreatedBefore)
		query = query.Where("created_at <= ?", beforeTime)
	}

	err := query.Find(&participants).Error
	if err != nil {
		return nil, err
	}

	for i := range participants {
		participants[i].StatusStr = participants[i].GetStatusDescription()
		participants[i].RoleStr = participants[i].GetRoleDescription()
	}

	return participants, nil
}
