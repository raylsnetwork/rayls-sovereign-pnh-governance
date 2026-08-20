package indexer

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/mocks"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/config"
)

func TestNewProvider(t *testing.T) {
	tests := []struct {
		name   string
		config *config.Config
	}{
		{
			name:   "creates provider with valid config",
			config: &config.Config{},
		},
		{
			name: "creates provider with populated config",
			config: &config.Config{
				PrivateHub: config.PrivateHub{
					URL: "http://localhost:8545",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mocks.NewMockEthereumClient(ctrl)
			provider := NewProvider(mockClient, tt.config)

			assert.NotNil(t, provider)
		})
	}
}

func TestProvider_GetLatestBlock(t *testing.T) {
	expectedBlock := types.NewBlockWithHeader(&types.Header{
		Number: big.NewInt(12345),
	})

	tests := []struct {
		name        string
		mockSetup   func(m *mocks.MockEthereumClient)
		ctx         context.Context
		expectError bool
		expectBlock *types.Block
		errorCheck  func(t *testing.T, err error)
	}{
		{
			name: "successfully returns latest block",
			mockSetup: func(m *mocks.MockEthereumClient) {
				m.EXPECT().BlockByNumber(gomock.Any(), nil).Return(expectedBlock, nil)
			},
			ctx:         context.Background(),
			expectError: false,
			expectBlock: expectedBlock,
		},
		{
			name: "returns error when client fails",
			mockSetup: func(m *mocks.MockEthereumClient) {
				m.EXPECT().BlockByNumber(gomock.Any(), nil).Return(nil, errors.New("connection failed"))
			},
			ctx:         context.Background(),
			expectError: true,
			expectBlock: nil,
			errorCheck: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "connection failed")
			},
		},
		{
			name: "handles context cancellation",
			mockSetup: func(m *mocks.MockEthereumClient) {
				m.EXPECT().BlockByNumber(gomock.Any(), nil).Return(nil, context.Canceled)
			},
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			expectError: true,
			expectBlock: nil,
			errorCheck: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, context.Canceled), "error should be context.Canceled")
			},
		},
		{
			name: "handles wrapped context cancellation error",
			mockSetup: func(m *mocks.MockEthereumClient) {
				wrappedErr := errors.Join(errors.New("rpc error"), context.Canceled)
				m.EXPECT().BlockByNumber(gomock.Any(), nil).Return(nil, wrappedErr)
			},
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			expectError: true,
			expectBlock: nil,
			errorCheck: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, context.Canceled), "wrapped error should contain context.Canceled")
			},
		},
		{
			name: "handles context deadline exceeded",
			mockSetup: func(m *mocks.MockEthereumClient) {
				m.EXPECT().BlockByNumber(gomock.Any(), nil).Return(nil, context.DeadlineExceeded)
			},
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				return ctx
			}(),
			expectError: true,
			expectBlock: nil,
			errorCheck: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, context.DeadlineExceeded), "error should be context.DeadlineExceeded")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mocks.NewMockEthereumClient(ctrl)
			tt.mockSetup(mockClient)

			provider := NewProvider(mockClient, &config.Config{})

			// Execute
			block, err := provider.GetLatestBlock(tt.ctx)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, block)
				if tt.errorCheck != nil {
					tt.errorCheck(t, err)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectBlock, block)
			}
		})
	}
}

func TestProvider_GetBlockByNumber(t *testing.T) {
	blockNumber := uint64(12345)
	expectedBlock := types.NewBlockWithHeader(&types.Header{
		Number: big.NewInt(int64(blockNumber)), // #nosec G115 -- test constant
	})

	tests := []struct {
		name        string
		blockNumber uint64
		mockSetup   func(m *mocks.MockEthereumClient)
		ctx         context.Context
		expectError bool
		expectBlock *types.Block
		errorCheck  func(t *testing.T, err error)
	}{
		{
			name:        "successfully returns block by number",
			blockNumber: blockNumber,
			mockSetup: func(m *mocks.MockEthereumClient) {
				m.EXPECT().
					BlockByNumber(gomock.Any(), big.NewInt(int64(blockNumber))).
					Return(expectedBlock, nil)
				// #nosec G115 -- test constant
			},
			ctx:         context.Background(),
			expectError: false,
			expectBlock: expectedBlock,
		},
		{
			name:        "returns error when client fails",
			blockNumber: blockNumber,
			mockSetup: func(m *mocks.MockEthereumClient) {
				m.EXPECT().
					BlockByNumber(gomock.Any(), big.NewInt(int64(blockNumber))). // #nosec G115 -- test constant
					Return(nil, errors.New("block not found"))
			},
			ctx:         context.Background(),
			expectError: true,
			expectBlock: nil,
			errorCheck: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "block not found")
			},
		},
		{
			name:        "handles context cancellation",
			blockNumber: blockNumber,
			mockSetup: func(m *mocks.MockEthereumClient) {
				m.EXPECT().
					BlockByNumber(gomock.Any(), big.NewInt(int64(blockNumber))).
					Return(nil, context.Canceled)
				// #nosec G115 -- test constant
			},
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			expectError: true,
			expectBlock: nil,
			errorCheck: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, context.Canceled))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mocks.NewMockEthereumClient(ctrl)
			tt.mockSetup(mockClient)

			provider := NewProvider(mockClient, &config.Config{})

			// Execute
			block, err := provider.GetBlockByNumber(tt.ctx, tt.blockNumber)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, block)
				if tt.errorCheck != nil {
					tt.errorCheck(t, err)
				}
			} else {
				require.NoError(t, err)
				assert.NotNil(t, block)
				if tt.expectBlock != nil {
					assert.Equal(t, tt.expectBlock.Number(), block.Number())
				}
			}
		})
	}
}
