package config

import (
	"context"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/config"
)

func TestNewConfigProvider(t *testing.T) {
	tests := []struct {
		name   string
		config *config.Config
	}{
		{
			name:   "creates provider with empty config",
			config: &config.Config{},
		},
		{
			name: "creates provider with full config",
			config: &config.Config{
				PrivateHub: config.PrivateHub{
					URL:           "http://localhost:8545",
					StartingBlock: "100",
					BatchSize:     50,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			provider := NewConfigProvider(tt.config)

			// Assert
			assert.NotNil(t, provider)

			// Verify provider works by calling a method
			blockNum := provider.GetStartingBlockNumber()
			assert.NotNil(t, blockNum)
		})
	}
}

func TestConfigProvider_GetStartingBlockNumber(t *testing.T) {
	tests := []struct {
		name          string
		startingBlock string
		expectedBlock *big.Int
	}{
		{
			name:          "returns valid block number from config",
			startingBlock: "12345",
			expectedBlock: big.NewInt(12345),
		},
		{
			name:          "returns block 0 for empty string",
			startingBlock: "",
			expectedBlock: big.NewInt(0),
		},
		{
			name:          "returns block 0 for invalid string",
			startingBlock: "invalid",
			expectedBlock: big.NewInt(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			cfg := &config.Config{
				PrivateHub: config.PrivateHub{
					StartingBlock: tt.startingBlock,
				},
			}
			provider := NewConfigProvider(cfg)

			// Execute
			blockNumber := provider.GetStartingBlockNumber()

			// Assert
			require.NotNil(t, blockNumber)
			assert.Equal(t, 0, tt.expectedBlock.Cmp(blockNumber),
				"Expected %s but got %s", tt.expectedBlock.String(), blockNumber.String())
		})
	}
}

func TestConfigProvider_GetBatchSize(t *testing.T) {
	tests := []struct {
		name              string
		batchSize         int64
		expectedBatchSize int64
	}{
		{
			name:              "returns configured batch size",
			batchSize:         50,
			expectedBatchSize: 50,
		},
		{
			name:              "returns default batch size for zero",
			batchSize:         0,
			expectedBatchSize: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			cfg := &config.Config{
				PrivateHub: config.PrivateHub{
					BatchSize: tt.batchSize,
				},
			}
			provider := NewConfigProvider(cfg)

			// Execute
			batchSize := provider.GetBatchSize()

			// Assert
			assert.Equal(t, tt.expectedBatchSize, batchSize)
		})
	}
}

func TestConfigProvider_GetChainURL(t *testing.T) {
	tests := []struct {
		name        string
		chainURL    string
		expectedURL string
	}{
		{
			name:        "returns configured chain URL",
			chainURL:    "http://localhost:8545",
			expectedURL: "http://localhost:8545",
		},
		{
			name:        "returns empty string when not configured",
			chainURL:    "",
			expectedURL: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			cfg := &config.Config{
				PrivateHub: config.PrivateHub{
					URL: tt.chainURL,
				},
			}
			provider := NewConfigProvider(cfg)

			// Execute
			url := provider.(*ConfigProvider).GetChainURL()

			// Assert
			assert.Equal(t, tt.expectedURL, url)
		})
	}
}

func TestConfigProvider_GetConfig(t *testing.T) {
	tests := []struct {
		name   string
		config *config.Config
		ctx    context.Context
	}{
		{
			name: "returns config successfully",
			config: &config.Config{
				PrivateHub: config.PrivateHub{
					URL:           "http://localhost:8545",
					StartingBlock: "100",
					BatchSize:     50,
				},
			},
			ctx: context.Background(),
		},
		{
			name:   "returns empty config",
			config: &config.Config{},
			ctx:    context.Background(),
		},
		{
			name: "works with cancelled context",
			config: &config.Config{
				PrivateHub: config.PrivateHub{
					URL: "http://localhost:8545",
				},
			},
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			provider := NewConfigProvider(tt.config)

			// Execute
			cfg, err := provider.GetConfig(tt.ctx)

			// Assert
			require.NoError(t, err)
			assert.NotNil(t, cfg)
			assert.Equal(t, tt.config, cfg)

			// Verify the returned config is the same instance
			assert.Same(t, tt.config, cfg, "should return the same config instance")
		})
	}
}
