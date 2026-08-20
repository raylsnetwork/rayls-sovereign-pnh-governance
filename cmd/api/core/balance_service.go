package core

import (
	"context"
	"errors"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/domain"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/logger"
)

var _ BalanceService = (*balanceService)(nil)

// balanceService implements the BalanceService interface
type balanceService struct {
	balanceRepo BalanceRepository
	log         logger.Logger
}

// NewBalanceService creates a new balance service
func NewBalanceService(balanceRepo BalanceRepository, log logger.Logger) BalanceService {
	return &balanceService{
		balanceRepo: balanceRepo,
		log:         log,
	}
}

// GetBalancesInChain retrieves balance(s) for a chain
// If resourceId is "/" or empty, returns all balances in the chain (it can be "/" due to the wildcard *)
// Otherwise returns specific resource balance
func (s *balanceService) GetBalancesInChain(ctx context.Context, chainId, resourceId string) (any, error) {
	s.log.Info("Fetching balances", "chainId", chainId, "resourceId", resourceId)

	// Validate chainId
	if chainId == "" {
		return nil, NewValidationError("chainId", "Chain ID cannot be empty")
	}

	// Check if requesting all balances or specific resource
	if resourceId == "/" || resourceId == "" {
		// Get all balances in the chain
		balances, err := s.balanceRepo.FindAllInChain(ctx, chainId)
		if err != nil {
			s.log.Error("Failed to fetch all balances in chain", "error", err, "chainId", chainId)
			return nil, NewInternalError("balance retrieval", err)
		}

		if len(balances) == 0 {
			s.log.Warn("No tokens found for chain", "chainId", chainId)
			return nil, NewNotFoundError("balances", "No tokens for this chain id found")
		}

		s.log.Info("Successfully fetched all balances", "chainId", chainId, "count", len(balances))
		return balances, nil
	}

	// Get specific resource balance
	// Strip leading "/" if present
	if len(resourceId) > 0 && resourceId[0] == '/' {
		resourceId = resourceId[1:]
	}

	balance, err := s.balanceRepo.FindByChainAndResource(ctx, chainId, resourceId)
	if err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			s.log.Warn("Resource not found", "chainId", chainId, "resourceId", resourceId)
			return nil, NewNotFoundError("resource", "This resourceid does not exist")
		}
		s.log.Error("Failed to fetch balance for resource", "error", err, "chainId", chainId, "resourceId", resourceId)
		return nil, NewInternalError("balance retrieval", err)
	}

	s.log.Info("Successfully fetched balance", "chainId", chainId, "resourceId", resourceId)
	return balance, nil
}

// GetBalanceAcrossAllChains retrieves balance for a resource across all chains
func (s *balanceService) GetBalanceAcrossAllChains(ctx context.Context, resourceId string) ([]domain.Balance, error) {
	s.log.Info("Fetching balance across all chains", "resourceId", resourceId)

	// Validate resourceId
	if resourceId == "" {
		return nil, NewValidationError("resourceId", "Resource ID cannot be empty")
	}

	// Get balances across all chains
	balances, err := s.balanceRepo.FindAcrossAllChains(ctx, resourceId)
	if err != nil {
		s.log.Error("Failed to fetch balance across all chains", "error", err, "resourceId", resourceId)
		return nil, NewInternalError("balance retrieval", err)
	}

	if len(balances) == 0 {
		s.log.Warn("Resource not found in any chain", "resourceId", resourceId)
		return nil, NewNotFoundError("resource", "Resource id not present in any pl")
	}

	s.log.Info("Successfully fetched balance across all chains", "resourceId", resourceId, "count", len(balances))
	return balances, nil
}

// GetBalanceAcrossSpecificChains retrieves balance for a resource across specific chains
// If chains is empty, delegates to GetBalanceAcrossAllChains
func (s *balanceService) GetBalanceAcrossSpecificChains(
	ctx context.Context,
	resourceId string,
	chains []string,
) ([]domain.Balance, error) {
	s.log.Info("Fetching balance across specific chains", "resourceId", resourceId, "chains", chains)

	// Validate resourceId
	if resourceId == "" {
		return nil, NewValidationError("resourceId", "Resource ID cannot be empty")
	}

	// If chains is empty, delegate to GetBalanceAcrossAllChains (reuse existing logic)
	if len(chains) == 0 {
		s.log.Info("Chains list empty, delegating to GetBalanceAcrossAllChains")
		return s.GetBalanceAcrossAllChains(ctx, resourceId)
	}

	// Get balances for specific chains
	balances, err := s.balanceRepo.FindAcrossSpecificChains(ctx, resourceId, chains)
	if err != nil {
		s.log.Error(
			"Failed to fetch balance across specific chains",
			"error",
			err,
			"resourceId",
			resourceId,
			"chains",
			chains,
		)
		return nil, NewInternalError("balance retrieval", err)
	}

	if len(balances) == 0 {
		s.log.Warn("Resource not found in any of the specified chains", "resourceId", resourceId, "chains", chains)
		return nil, NewNotFoundError("resource", "Resource id not present in any of pls list")
	}

	s.log.Info(
		"Successfully fetched balance across specific chains",
		"resourceId",
		resourceId,
		"count",
		len(balances),
	)
	return balances, nil
}
