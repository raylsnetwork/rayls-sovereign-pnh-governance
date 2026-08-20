package core

import (
	"context"
	"encoding/hex"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/api/utils"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/dto"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/logger"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/types"
)

var _ TokenService = (*tokenService)(nil)

// tokenService implements the TokenService interface
type tokenService struct {
	repo TokenRepository
	log  logger.Logger
}

// NewTokenService creates a new token service
func NewTokenService(repo TokenRepository, log logger.Logger) TokenService {
	return &tokenService{
		repo: repo,
		log:  log,
	}
}

// GetTokenByResourceId retrieves a single token by resource ID with balances
func (s *tokenService) GetTokenByResourceId(
	ctx context.Context,
	resourceId string,
) (*domain.TokenWithBalancesAndFreezeState, error) {
	// Validate resourceId is not empty
	if resourceId == "" {
		return nil, NewValidationError("resourceId", "resourceId cannot be empty")
	}

	// Validate resourceId is valid hex
	buff := make([]byte, hex.DecodedLen(len(resourceId)))
	if _, err := hex.Decode(buff, []byte(resourceId)); err != nil {
		return nil, NewValidationError("resourceId", "resourceId must be a valid hex string")
	}

	// Fetch from repository (with balances/circulating supply)
	token, err := s.repo.FindByResourceIdWithBalances(ctx, resourceId)
	if err != nil {
		return nil, NewInternalError("FindByResourceIdWithBalances", err)
	}

	// Check if token exists
	if token.ResourceId == "" {
		return nil, NewNotFoundError("token", resourceId)
	}

	return token, nil
}

// GetTokenRegistryStatus retrieves token registry status (subset of fields) by resource ID
func (s *tokenService) GetTokenRegistryStatus(
	ctx context.Context,
	resourceId string,
) (*dto.TokenRegistryStatusDto, error) {
	token, err := s.GetTokenByResourceId(ctx, resourceId)
	if err != nil {
		return nil, err
	}

	return &dto.TokenRegistryStatusDto{
		CreatedAt:   token.CreatedAt,
		UpdatedAt:   token.UpdatedAt,
		Name:        token.Name,
		Symbol:      token.Symbol,
		ResourceId:  token.ResourceId,
		MetadataUrl: token.MetadataUrl,
		ErcStandard: token.ErcStandard,
		Decimals:    token.Decimals,
		IssuerId:    token.IssuerId,
		Status:      token.Status,
	}, nil
}

// GetTokensList retrieves a paginated list of tokens with filters
func (s *tokenService) GetTokensList(
	ctx context.Context,
	filters dto.TokenListFilters,
) (*types.Paginated[domain.TokenWithBalancesAndFreezeState], error) {
	// Validate and normalize filters
	if err := s.validateFilters(&filters); err != nil {
		return nil, err
	}

	// Fetch from repository
	tokens, total, err := s.repo.FindByFilters(ctx, filters)
	if err != nil {
		return nil, NewInternalError("FindByFilters", err)
	}

	return &types.Paginated[domain.TokenWithBalancesAndFreezeState]{
		Data:  tokens,
		Total: total,
		Limit: filters.Limit,
		Page:  filters.Page,
	}, nil
}

// validateFilters validates and sanitizes token list query parameters
func (s *tokenService) validateFilters(filters *dto.TokenListFilters) error {
	// Normalize pagination
	if filters.Limit < 1 || filters.Limit > 100 {
		filters.Limit = 10
	}
	if filters.Page < 1 {
		filters.Page = 1
	}

	// Validate timestamp range
	if _, _, err := utils.ValidateTimestampRange(filters.CreatedAfter, filters.CreatedBefore); err != nil {
		return NewValidationError("createdAfter/createdBefore", err.Error())
	}

	return nil
}
