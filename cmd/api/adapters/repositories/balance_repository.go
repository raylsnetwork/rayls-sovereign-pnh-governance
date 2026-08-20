package repositories

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/api/core"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
)

var _ core.BalanceRepository = (*balanceRepository)(nil)

// balanceRepository implements core.BalanceRepository
type balanceRepository struct {
	db *gorm.DB
}

// NewBalanceRepository creates a new balance repository
func NewBalanceRepository(db *gorm.DB) core.BalanceRepository {
	return &balanceRepository{db: db}
}

// FindAllInChain retrieves all balances for a specific chain
func (r *balanceRepository) FindAllInChain(ctx context.Context, chainId string) ([]domain.Balance, error) {
	var balances []domain.Balance

	err := r.db.WithContext(ctx).
		Where("chain_id = ?", chainId).
		Find(&balances).Error
	if err != nil {
		return nil, err
	}

	return balances, nil
}

// FindByChainAndResource retrieves balance for specific resource in chain
func (r *balanceRepository) FindByChainAndResource(
	ctx context.Context,
	chainId, resourceId string,
) (*domain.Balance, error) {
	var balance domain.Balance

	err := r.db.WithContext(ctx).
		Where("chain_id = ? AND resource_id = ?", chainId, resourceId).
		First(&balance).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrRecordNotFound
		}
		return nil, err
	}

	return &balance, nil
}

// FindAcrossAllChains retrieves balances for a resource across all chains
func (r *balanceRepository) FindAcrossAllChains(ctx context.Context, resourceId string) ([]domain.Balance, error) {
	var balances []domain.Balance

	err := r.db.WithContext(ctx).
		Where("resource_id = ?", resourceId).
		Find(&balances).Error
	if err != nil {
		return nil, err
	}

	return balances, nil
}

// FindAcrossSpecificChains retrieves balances for a resource across specific chains
func (r *balanceRepository) FindAcrossSpecificChains(
	ctx context.Context,
	resourceId string,
	chains []string,
) ([]domain.Balance, error) {
	var balances []domain.Balance

	err := r.db.WithContext(ctx).
		Where("resource_id = ? AND chain_id IN (?)", resourceId, chains).
		Find(&balances).Error
	if err != nil {
		return nil, err
	}

	return balances, nil
}
