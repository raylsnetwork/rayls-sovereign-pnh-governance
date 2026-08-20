package indexer

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/core"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/config"
)

// Ensure Provider implements core.Provider at compile time
var _ core.Provider = (*Provider)(nil)

// EthereumClient exists solely to enable mocking the Ethereum client in unit tests
type EthereumClient interface {
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
}

// Provider is an adapter that provides block data from the blockchain
type Provider struct {
	client EthereumClient
	config *config.Config
}

// NewProvider creates a new provider with an injected Ethereum client and config
func NewProvider(client EthereumClient, config *config.Config) core.Provider {
	return &Provider{
		client: client,
		config: config,
	}
}

// GetLatestBlock retrieves the latest block from the Ethereum client
func (p *Provider) GetLatestBlock(ctx context.Context) (*types.Block, error) {
	block, err := p.client.BlockByNumber(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block: %w", err)
	}
	return block, nil
}

// GetBlockByNumber retrieves a specific block by its number from the Ethereum client.
func (p *Provider) GetBlockByNumber(ctx context.Context, number uint64) (*types.Block, error) {
	block, err := p.client.BlockByNumber(ctx, new(big.Int).SetUint64(number))
	if err != nil {
		return nil, fmt.Errorf("failed to get block by number %d from ethereum client: %w", number, err)
	}
	return block, nil
}
