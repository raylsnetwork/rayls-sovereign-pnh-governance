package core

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/api/utils"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/domain"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/dto"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/logger"
)

var _ ParticipantService = (*participantService)(nil)

// participantService implements the ParticipantService interface
type participantService struct {
	repo ParticipantRepository
	log  logger.Logger
}

// NewParticipantService creates a new participant service
func NewParticipantService(repo ParticipantRepository, log logger.Logger) ParticipantService {
	return &participantService{
		repo: repo,
		log:  log,
	}
}

// GetParticipantByChainId retrieves a single participant by chain ID
func (s *participantService) GetParticipantByChainId(ctx context.Context, chainId string) (*domain.Participant, error) {
	// Validate chainId
	if chainId == "" {
		return nil, NewValidationError("chainId", "chainId cannot be empty")
	}
	if _, err := strconv.ParseInt(chainId, 10, 64); err != nil {
		return nil, NewValidationError("chainId", "chainId must be an integer")
	}

	// Fetch from repository
	participant, err := s.repo.FindByChainId(ctx, chainId)
	if err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			return nil, NewNotFoundError("participant", chainId)
		}
		return nil, NewInternalError("FindByChainId", err)
	}

	return participant, nil
}

// GetParticipantsList retrieves a list of participants with filters
func (s *participantService) GetParticipantsList(
	ctx context.Context,
	filters dto.ParticipantListFilters,
) ([]domain.Participant, error) {
	// Validate and normalize status enum if provided
	if filters.Status != "" {
		filters.Status = strings.ToLower(filters.Status)
		if _, exists := domain.StringToMemberStatus[filters.Status]; !exists {
			allowedValues := domain.GetAllowedMemberStatuses()
			return nil, NewValidationError("status", "invalid status, allowed values: "+allowedValues)
		}
	}

	// Validate and normalize role enum if provided
	if filters.Role != "" {
		filters.Role = strings.ToLower(filters.Role)
		if _, exists := domain.StringToMemberRole[filters.Role]; !exists {
			allowedValues := domain.GetAllowedMemberRoles()
			return nil, NewValidationError("role", "invalid role, allowed values: "+allowedValues)
		}
	}

	// Validate timestamp range
	if _, _, err := utils.ValidateTimestampRange(filters.CreatedAfter, filters.CreatedBefore); err != nil {
		return nil, NewValidationError("createdAfter/createdBefore", err.Error())
	}

	// Fetch from repository
	participants, err := s.repo.FindByFilters(ctx, filters)
	if err != nil {
		return nil, NewInternalError("FindByFilters", err)
	}

	return participants, nil
}
