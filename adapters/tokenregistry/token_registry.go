package tokenregistry

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/config"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/domain"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/types"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/withstack"
)

// TokenRegistryAdapter provides access to the TokenRegistry smart contract
type TokenRegistryAdapter struct {
	tokenRegistry *contracts.TokenRegistryClient
	client        *ethclient.Client
	config        *config.Config
}

// NewTokenRegistryAdapter creates a new TokenRegistryAdapter instance
func NewTokenRegistryAdapter(
	ctx context.Context,
	client *ethclient.Client,
	conf *config.Config,
) (*TokenRegistryAdapter, error) {
	tokenRegistry, err := contracts.CreateTokenRegistry(ctx, conf.PrivateHub.TokenRegistry, client)
	if err != nil {
		return nil, withstack.Wrap(err)
	}

	return &TokenRegistryAdapter{
		tokenRegistry: tokenRegistry,
		client:        client,
		config:        conf,
	}, nil
}

// GetTokenByResourceId retrieves a specific token by resource ID from the smart contract registry
func (t *TokenRegistryAdapter) GetTokenByResourceId(ctx context.Context, resourceId string) (*domain.Token, error) {
	// Convert hex string to [32]byte
	resourceIdBytes, err := hex.DecodeString(resourceId)
	if err != nil {
		return nil, fmt.Errorf("failed to decode resource ID: %w", err)
	}

	var resourceIdArray [32]byte
	copy(resourceIdArray[:], resourceIdBytes)

	token, err := t.tokenRegistry.GetTokenByResourceId(&bind.CallOpts{Context: ctx}, resourceIdArray)
	if err != nil {
		return nil, fmt.Errorf("failed to get token by resource ID from contract: %w", err)
	}

	result := &domain.Token{
		Name:        token.Name,
		Symbol:      token.Symbol,
		ResourceId:  hex.EncodeToString(token.ResourceId[:]),
		MetadataUrl: token.Metadata.Url,
		// Fold FACTORY-mode test ordinals (ERC20TEST..) back to their production base so the
		// example token indexes identically to production (ercStandard string, ercId, balances).
		ErcStandard: types.NormalizeAssetType(token.ErcStandard),
		Decimals:    token.Metadata.Decimals,
		IssuerId:    token.IssuerChainId.String(),
		Status:      token.Status,
	}

	if token.CreatedAt != nil && token.CreatedAt.Sign() > 0 {
		result.CreatedAt = time.Unix(token.CreatedAt.Int64(), 0).UTC()
	}
	if token.UpdatedAt != nil && token.UpdatedAt.Sign() > 0 {
		result.UpdatedAt = time.Unix(token.UpdatedAt.Int64(), 0).UTC()
	}

	return result, nil
}
