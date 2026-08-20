package config

import (
	"context"
	"math/big"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/core"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/config"
)

// Ensure ConfigProvider implements core.ConfigProvider at compile time
var _ core.ConfigProvider = (*ConfigProvider)(nil)

// ConfigProvider implements the ConfigProvider interface using the actual config system
type ConfigProvider struct {
	config *config.Config
}

// NewConfigProvider creates a new ConfigProvider
func NewConfigProvider(config *config.Config) core.ConfigProvider {
	return &ConfigProvider{
		config: config,
	}
}

// GetStartingBlockNumber implements the ConfigProvider interface
func (c *ConfigProvider) GetStartingBlockNumber() *big.Int {
	blockNumber, ok := new(big.Int).SetString(c.config.PrivateHub.StartingBlock, 10)
	if !ok {
		// If conversion fails, default to block 0
		return big.NewInt(0)
	}

	return blockNumber
}

// GetBatchSize implements the ConfigProvider interface
func (c *ConfigProvider) GetBatchSize() int64 {
	if c.config.PrivateHub.BatchSize <= 0 {
		// Default batch size
		return 10
	}
	return c.config.PrivateHub.BatchSize
}

// GetChainURL returns the blockchain RPC URL
func (c *ConfigProvider) GetChainURL() string {
	return c.config.PrivateHub.URL
}

// GetConfig implements core.ConfigProvider interface
func (c *ConfigProvider) GetConfig(ctx context.Context) (*config.Config, error) {
	return c.config, nil
}
