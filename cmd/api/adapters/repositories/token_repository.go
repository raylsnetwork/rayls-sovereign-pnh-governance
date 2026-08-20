package repositories

import (
	"context"
	"strings"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/api/core"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/domain"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/dto"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/types"
)

var _ core.TokenRepository = (*tokenRepository)(nil)

// TokenAggregatedRow represents the result of the CTE query
type TokenAggregatedRow struct {
	domain.Token
	TotalSupply       decimal.Decimal
	CirculatingSupply domain.CirculatingSupplyJSON `gorm:"type:json"`
	FrozenChainIds    domain.FrozenChainIdsJSON    `gorm:"type:json"`
}

// tokenCTESQL is the base CTE query that pre-aggregates balances and freeze states
//
//nolint:gosec // G101 false positive
const tokenCTESQL = `
WITH token_balances AS (
    SELECT
        resource_id,
        SUM(amount) AS total_supply,
        json_agg(jsonb_build_object(
            'participantId', chain_id,
            'balance',       amount,
            'tokenId',       erc_id::text
        	)
			ORDER BY chain_id
		) AS circulating_supply
    FROM balances
    GROUP BY resource_id
),
token_frozen_states AS (
    SELECT resource_id, json_agg(chain_id) AS frozen_chain_ids
    FROM token_freeze_states
    WHERE is_frozen = true
    GROUP BY resource_id
)
SELECT
    t.id, t.created_at, t.updated_at, t.name, t.symbol,
    t.resource_id, t.metadata_url, t.erc_standard,
    t.decimals, t.issuer_id, t.status,
    COALESCE(tb.total_supply, 0) AS total_supply,
    COALESCE(tb.circulating_supply, '[]') AS circulating_supply,
    COALESCE(tr.frozen_chain_ids, '[]') AS frozen_chain_ids
FROM tokens t
LEFT JOIN token_balances     tb ON t.resource_id = tb.resource_id
LEFT JOIN token_frozen_states tr ON t.resource_id = tr.resource_id
`

// tokenRepository implements core.TokenRepository using GORM
type tokenRepository struct {
	db *gorm.DB
}

// NewTokenRepository creates a new GORM-based token repository
func NewTokenRepository(db *gorm.DB) core.TokenRepository {
	return &tokenRepository{db: db}
}

// FindByResourceIdWithBalances retrieves a token with its balances and freeze restrictions
func (r *tokenRepository) FindByResourceIdWithBalances(
	ctx context.Context,
	resourceId string,
) (*domain.TokenWithBalancesAndFreezeState, error) {
	var row TokenAggregatedRow

	err := r.db.WithContext(ctx).
		Raw(tokenCTESQL+"WHERE t.resource_id = ? LIMIT 1", resourceId).
		Scan(&row).Error
	if err != nil {
		return nil, err
	}

	if row.ResourceId == "" {
		return &domain.TokenWithBalancesAndFreezeState{}, nil
	}

	result := mapAggregatedRows([]TokenAggregatedRow{row})
	return &result[0], nil
}

// applyTokenFilters adds WHERE clauses to the query based on the provided filters
func (r *tokenRepository) applyTokenFilters(query *gorm.DB, filters dto.TokenListFilters) *gorm.DB {
	if filters.Name != "" {
		query = query.Where("t.name ILIKE ?", "%"+filters.Name+"%")
	}
	if filters.Symbol != "" {
		query = query.Where("t.symbol ILIKE ?", "%"+filters.Symbol+"%")
	}
	if filters.IssuerId != "" {
		query = query.Where("t.issuer_id = ?", filters.IssuerId)
	}
	if filters.Status != "" {
		if statusInt, exists := domain.StringToTokenStatus[strings.ToLower(filters.Status)]; exists {
			query = query.Where("t.status = ?", statusInt)
		}
	}
	if filters.ErcStandard != "" {
		if ercInt, exists := types.StringToAssetType[strings.ToLower(filters.ErcStandard)]; exists {
			query = query.Where("t.erc_standard = ?", ercInt)
		}
	}
	if filters.Decimals != nil {
		query = query.Where("t.decimals = ?", *filters.Decimals)
	}
	if filters.CreatedAfter != "" {
		query = query.Where("t.created_at >= ?", filters.CreatedAfter)
	}
	if filters.CreatedBefore != "" {
		query = query.Where("t.created_at < ?", filters.CreatedBefore)
	}
	return query
}

// FindByFilters finds tokens matching the provided filters with pagination
func (r *tokenRepository) FindByFilters(
	ctx context.Context,
	filters dto.TokenListFilters,
) ([]domain.TokenWithBalancesAndFreezeState, int64, error) {
	baseQuery := r.applyTokenFilters(r.db.WithContext(ctx).Table("tokens t"), filters)

	// Count distinct tokens (not JOIN rows)
	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []domain.TokenWithBalancesAndFreezeState{}, 0, nil
	}

	// Get paginated token IDs first
	offset := (filters.Page - 1) * filters.Limit
	var tokenIds []string
	if err := baseQuery.Select("t.resource_id").
		Order("t.created_at DESC").
		Limit(filters.Limit).
		Offset(offset).
		Pluck("t.resource_id", &tokenIds).Error; err != nil {
		return nil, 0, err
	}

	if len(tokenIds) == 0 {
		return []domain.TokenWithBalancesAndFreezeState{}, total, nil
	}

	// Fetch aggregated data via CTE for the paginated tokens
	var rows []TokenAggregatedRow
	err := r.db.WithContext(ctx).
		Raw(tokenCTESQL+"WHERE t.resource_id IN ? ORDER BY t.created_at DESC, t.resource_id", tokenIds).
		Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}

	return mapAggregatedRows(rows), total, nil
}

// mapAggregatedRows converts CTE query results into domain objects
func mapAggregatedRows(rows []TokenAggregatedRow) []domain.TokenWithBalancesAndFreezeState {
	result := make([]domain.TokenWithBalancesAndFreezeState, 0, len(rows))
	for _, row := range rows {
		circulatingSupply := []domain.CirculatingSupplyEntry(row.CirculatingSupply)
		if circulatingSupply == nil {
			circulatingSupply = []domain.CirculatingSupplyEntry{}
		}
		frozenChainIds := []string(row.FrozenChainIds)
		if frozenChainIds == nil {
			frozenChainIds = []string{}
		}
		result = append(result, domain.TokenWithBalancesAndFreezeState{
			Token:             row.Token,
			TotalSupply:       row.TotalSupply,
			CirculatingSupply: circulatingSupply,
			FrozenChainIds:    frozenChainIds,
		})
	}
	return result
}
