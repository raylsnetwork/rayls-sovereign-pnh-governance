package handlers

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
	goUUID "github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/core"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/mocks"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/testutil"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/TokenCoreV1"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/domain"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/types"
)

// TestEnsureTokenExists removed - functionality replaced by Upsert method using OnConflict

func TestProcessTokenRegistered(t *testing.T) {
	tests := []struct {
		name          string
		resourceID    [32]byte
		eventIssuerId *big.Int
		setupMocks    func(*mocks.MockTokenService, *mocks.MockTokenRepository)
		expectedToken *domain.Token
		expectedError bool
		errorContains string
	}{
		{
			name:          "tokenService.GetTokenByResourceId fails - returns error",
			resourceID:    [32]byte{1, 2, 3},
			eventIssuerId: big.NewInt(123),
			setupMocks: func(mockTokenService *mocks.MockTokenService, mockTokenRepo *mocks.MockTokenRepository) {
				// tokenService fails to get token from token registry
				mockTokenService.EXPECT().
					GetTokenByResourceId(gomock.Any(), "0102030000000000000000000000000000000000000000000000000000000000").
					Return(nil, errors.New("connection error"))
			},
			expectedToken: nil,
			expectedError: true,
			errorContains: "connection error",
		},
		{
			name:          "blockchain returns nil - creates and returns stub token",
			resourceID:    [32]byte{4, 5, 6},
			eventIssuerId: big.NewInt(456),
			setupMocks: func(mockTokenService *mocks.MockTokenService, mockTokenRepo *mocks.MockTokenRepository) {
				// tokenService returns nil (token not yet in blockchain registry, e.g. timing issue)
				mockTokenService.EXPECT().
					GetTokenByResourceId(gomock.Any(), "0405060000000000000000000000000000000000000000000000000000000000").
					Return(nil, nil)
				// A stub token must be upserted to keep the tokens table consistent with transaction records
				mockTokenRepo.EXPECT().
					Upsert(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedToken: &domain.Token{
				ResourceId: "0405060000000000000000000000000000000000000000000000000000000000",
				IssuerId:   "456",
			},
			expectedError: false,
		},
		{
			name:          "issuer mismatch - returns ErrIssuerMismatch",
			resourceID:    [32]byte{7, 8, 9},
			eventIssuerId: big.NewInt(999),
			setupMocks: func(mockTokenService *mocks.MockTokenService, mockTokenRepo *mocks.MockTokenRepository) {
				// tokenService returns a token with different issuer
				tokenWithWrongIssuer := &domain.Token{
					IssuerId: "111",
				}
				mockTokenService.EXPECT().
					GetTokenByResourceId(gomock.Any(), "0708090000000000000000000000000000000000000000000000000000000000").
					Return(tokenWithWrongIssuer, nil)
			},
			expectedToken: nil,
			expectedError: true,
			errorContains: "token issuer mismatch",
		},
		{
			name: "Upsert fails - returns error",
			// Note: This test verifies error propagation from Upsert.
			resourceID:    [32]byte{10, 11, 12},
			eventIssuerId: big.NewInt(789),
			setupMocks: func(mockTokenService *mocks.MockTokenService, mockTokenRepo *mocks.MockTokenRepository) {
				// tokenService returns a valid token with matching issuer
				validToken := &domain.Token{
					ResourceId: "0a0b0c0000000000000000000000000000000000000000000000000000000000",
					IssuerId:   "789",
				}

				mockTokenService.EXPECT().
					GetTokenByResourceId(gomock.Any(), "0a0b0c0000000000000000000000000000000000000000000000000000000000").
					Return(validToken, nil)

				// Upsert fails
				mockTokenRepo.EXPECT().
					Upsert(gomock.Any(), gomock.Any()).
					Return(errors.New("database error"))
			},
			expectedToken: nil,
			expectedError: true,
			errorContains: "database error",
		},
		{
			name: "success - token registered and upserted",
			// Note: This test verifies the complete happy path
			// This integration test ensures:
			// 1. Valid token with matching issuer passes validation
			// 2. Upsert is called and succeeds
			// 3. The token from registry is returned to caller
			resourceID:    [32]byte{13, 14, 15},
			eventIssuerId: big.NewInt(888),
			setupMocks: func(mockTokenService *mocks.MockTokenService, mockTokenRepo *mocks.MockTokenRepository) {
				// tokenService returns a valid token with matching issuer
				validToken := &domain.Token{
					ResourceId: "0d0e0f0000000000000000000000000000000000000000000000000000000000",
					IssuerId:   "888",
				}

				mockTokenService.EXPECT().
					GetTokenByResourceId(gomock.Any(), "0d0e0f0000000000000000000000000000000000000000000000000000000000").
					Return(validToken, nil)

				// Upsert succeeds
				mockTokenRepo.EXPECT().
					Upsert(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedToken: &domain.Token{
				ResourceId: "0d0e0f0000000000000000000000000000000000000000000000000000000000",
				IssuerId:   "888",
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			mockTokenService := mocks.NewMockTokenService(ctrl)
			mockTokenRepo := mocks.NewMockTokenRepository(ctrl)

			handler := &TokenCoreEventHandler{
				tokenService: mockTokenService,
				tokenRepo:    mockTokenRepo,
				log:          &testutil.StubLogger{},
			}

			if tt.setupMocks != nil {
				tt.setupMocks(mockTokenService, mockTokenRepo)
			}

			// Execute
			token, err := handler.persistTokenRegistry(ctx, tt.resourceID, tt.eventIssuerId)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, token)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				if tt.expectedToken != nil {
					assert.NotNil(t, token)
					assert.Equal(t, tt.expectedToken.ResourceId, token.ResourceId)
					assert.Equal(t, tt.expectedToken.IssuerId, token.IssuerId)
				} else {
					assert.Nil(t, token)
				}
			}
		})
	}
}

func TestProcessTokenBalanceUpdated(t *testing.T) {
	tests := []struct {
		name          string
		log           core.ContractLog
		setupMocks    func(*mocks.MockTokenService, *mocks.MockTokenRepository, *mocks.MockTransactionRepository)
		expectedError bool
		errorContains string
	}{
		{
			name: "type assertion fails - returns error",
			log: core.ContractLog{
				RawEventData: nil,
			},
			expectedError: true,
			errorContains: "failed to unmarshal event data",
		},
		{
			name: "processTokenRegistered returns error - returns error",
			log: core.ContractLog{
				RawEventData: testutil.MustMarshal(t, &TokenCoreV1.TokenCoreV1TokenBalanceUpdated{
					ResourceId:    [32]byte{1, 2, 3},
					IssuerChainId: big.NewInt(123),
				}),
				TransactionHash: "0xabc123",
				BlockNumber:     1000,
			},
			setupMocks: func(mockTokenService *mocks.MockTokenService, mockTokenRepo *mocks.MockTokenRepository, mockTxRepo *mocks.MockTransactionRepository) {
				// tokenService fails to get token (this is the first call in processTokenRegistered)
				mockTokenService.EXPECT().
					GetTokenByResourceId(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("connection failed"))
			},
			expectedError: true,
			errorContains: "connection failed",
		},
		{
			name: "blockchain returns nil - creates stub and persists update transaction",
			log: core.ContractLog{
				RawEventData: testutil.MustMarshal(t, &TokenCoreV1.TokenCoreV1TokenBalanceUpdated{
					ResourceId:    [32]byte{4, 5, 6},
					IssuerChainId: big.NewInt(456),
				}),
				TransactionHash: "0xdef456",
				BlockNumber:     2000,
			},
			setupMocks: func(mockTokenService *mocks.MockTokenService, mockTokenRepo *mocks.MockTokenRepository, mockTxRepo *mocks.MockTransactionRepository) {
				// tokenService returns nil (token not yet in blockchain registry, e.g. timing issue)
				mockTokenService.EXPECT().
					GetTokenByResourceId(gomock.Any(), gomock.Any()).
					Return(nil, nil)
				// Stub token must be created to keep tokens table consistent
				mockTokenRepo.EXPECT().
					Upsert(gomock.Any(), gomock.Any()).
					Return(nil)
				// Update transaction is persisted because the stub token is non-nil
				mockTxRepo.EXPECT().
					PersistUpdateTransaction(gomock.Any(), gomock.Any(), "456", "456").
					Return(goUUID.Nil, nil)
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			mockTokenService := mocks.NewMockTokenService(ctrl)
			mockTokenRepo := mocks.NewMockTokenRepository(ctrl)
			mockTxRepo := mocks.NewMockTransactionRepository(ctrl)

			handler := &TokenCoreEventHandler{
				tokenService: mockTokenService,
				tokenRepo:    mockTokenRepo,
				txRepo:       mockTxRepo,
				log:          &testutil.StubLogger{},
			}

			if tt.setupMocks != nil {
				tt.setupMocks(mockTokenService, mockTokenRepo, mockTxRepo)
			}

			// Execute
			err := handler.processTokenBalanceUpdated(ctx, tt.log)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProcessTokenStatusUpdated(t *testing.T) {
	tests := []struct {
		name          string
		log           core.ContractLog
		setupMocks    func(*mocks.MockTokenRepository)
		expectedError bool
		errorContains string
	}{
		{
			name: "type assertion fails - returns error",
			log: core.ContractLog{
				RawEventData: nil,
			},
			expectedError: true,
			errorContains: "failed to unmarshal event data",
		},
		{
			name: "GetByIssuerAndName fails - returns error",
			log: core.ContractLog{
				RawEventData: testutil.MustMarshal(t, &TokenCoreV1.TokenCoreV1TokenStatusUpdated{
					IssuerChainId: big.NewInt(123),
					Name:          "TestToken",
				}),
			},
			setupMocks: func(mockTokenRepo *mocks.MockTokenRepository) {
				// GetByIssuerAndName call fails
				mockTokenRepo.EXPECT().
					GetByIssuerAndName(gomock.Any(), "123", "TestToken").
					Return(nil, errors.New("db error"))
			},
			expectedError: true,
			errorContains: "db error",
		},
		{
			name: "success - token status updated",
			log: core.ContractLog{
				RawEventData: testutil.MustMarshal(t, &TokenCoreV1.TokenCoreV1TokenStatusUpdated{
					IssuerChainId: big.NewInt(789),
					Name:          "SuccessToken",
					Status:        3,
				}),
			},
			setupMocks: func(mockTokenRepo *mocks.MockTokenRepository) {
				// GetByIssuerAndName succeeds
				existingToken := &domain.Token{
					Status: 0,
				}
				mockTokenRepo.EXPECT().
					GetByIssuerAndName(gomock.Any(), "789", "SuccessToken").
					Return(existingToken, nil)

				// Upsert succeeds and verify status was updated to NEW status (3)
				mockTokenRepo.EXPECT().
					Upsert(gomock.Any(), gomock.Any()).
					Do(func(ctx context.Context, token *domain.Token) {
						assert.Equal(t, uint8(3), token.Status)
					}).
					Return(nil)
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			mockTokenRepo := mocks.NewMockTokenRepository(ctrl)

			handler := &TokenCoreEventHandler{
				tokenRepo: mockTokenRepo,
				log:       &testutil.StubLogger{},
			}

			if tt.setupMocks != nil {
				tt.setupMocks(mockTokenRepo)
			}

			// Execute
			err := handler.processTokenStatusUpdated(ctx, tt.log)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProcessErc20TokenRegistered(t *testing.T) {
	tests := []struct {
		name          string
		log           core.ContractLog
		setupMocks    func(*mocks.MockTokenService, *mocks.MockTokenRepository, *mocks.MockTransactionRepository)
		expectedError bool
		errorContains string
	}{
		{
			name: "type assertion fails - returns error",
			log: core.ContractLog{
				RawEventData: nil,
			},
			expectedError: true,
			errorContains: "failed to unmarshal event data",
		},
		{
			name: "blockchain returns nil - creates stub and persists initial mint transaction",
			log: core.ContractLog{
				RawEventData: testutil.MustMarshal(t, &TokenCoreV1.TokenCoreV1Erc20TokenRegistered{
					ResourceId:    [32]byte{4, 5, 6},
					IssuerChainId: big.NewInt(456),
					InitialSupply: big.NewInt(2000000),
					Raw: ethTypes.Log{
						TxHash:      common.HexToHash("0xdef456"),
						BlockNumber: 2000,
					},
				}),
			},
			setupMocks: func(mockTokenService *mocks.MockTokenService, mockTokenRepo *mocks.MockTokenRepository, mockTxRepo *mocks.MockTransactionRepository) {
				expectedResourceId := "0405060000000000000000000000000000000000000000000000000000000000"
				// tokenService returns nil (token not yet in blockchain registry)
				mockTokenService.EXPECT().
					GetTokenByResourceId(gomock.Any(), expectedResourceId).
					Return(nil, nil)
				// Stub token created to maintain consistency
				mockTokenRepo.EXPECT().
					Upsert(gomock.Any(), gomock.Any()).
					Return(nil)
				// Initial mint is persisted because stub token is non-nil
				mockTxRepo.EXPECT().
					PersistUpdateTransactions(gomock.Any(), gomock.Any(), "456", "456").
					Return(nil, nil)
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			mockTokenService := mocks.NewMockTokenService(ctrl)
			mockTokenRepo := mocks.NewMockTokenRepository(ctrl)
			mockTxRepo := mocks.NewMockTransactionRepository(ctrl)

			handler := &TokenCoreEventHandler{
				tokenService: mockTokenService,
				tokenRepo:    mockTokenRepo,
				txRepo:       mockTxRepo,
				log:          &testutil.StubLogger{},
			}

			if tt.setupMocks != nil {
				tt.setupMocks(mockTokenService, mockTokenRepo, mockTxRepo)
			}

			// Execute
			err := handler.processErc20TokenRegistered(ctx, tt.log)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProcessErc721TokenRegistered(t *testing.T) {
	tests := []struct {
		name          string
		log           core.ContractLog
		expectedError bool
		errorContains string
	}{
		{
			name: "type assertion fails - returns error",
			log: core.ContractLog{
				RawEventData: nil,
			},
			expectedError: true,
			errorContains: "failed to unmarshal event data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			ctx := context.Background()

			handler := &TokenCoreEventHandler{
				log: &testutil.StubLogger{},
			}

			// Execute
			err := handler.processErc721TokenRegistered(ctx, tt.log)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProcessErc1155TokenRegistered(t *testing.T) {
	tests := []struct {
		name          string
		log           core.ContractLog
		expectedError bool
		errorContains string
	}{
		{
			name: "type assertion fails - returns error",
			log: core.ContractLog{
				RawEventData: nil,
			},
			expectedError: true,
			errorContains: "failed to unmarshal event data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			ctx := context.Background()

			handler := &TokenCoreEventHandler{
				log: &testutil.StubLogger{},
			}

			// Execute
			err := handler.processErc1155TokenRegistered(ctx, tt.log)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTokenCoreEventHandler_Name(t *testing.T) {
	handler := NewTokenCoreEventHandler(nil, nil, nil, nil)

	name := handler.Name()

	assert.Equal(t, "TokenCoreHandler", name)
}

func TestTokenCoreEventHandler_PersistInitialMintTransactions(t *testing.T) {
	tests := []struct {
		name          string
		log           core.ContractLog
		resourceID    [32]byte
		msgType       uint8
		issuerChainID *big.Int
		supplies      []MintSupply
		setupMocks    func(*mocks.MockTransactionRepository)
		expectedError bool
		errorContains string
	}{
		{
			name: "success - single mint transaction",
			log: core.ContractLog{
				TransactionHash: "0xabc123",
				BlockNumber:     1000,
				LogIndex:        1,
			},
			resourceID:    [32]byte{1, 2, 3, 4},
			msgType:       uint8(types.AssetTypeERC721),
			issuerChainID: big.NewInt(1),
			supplies: []MintSupply{
				{
					ercID:  nil,
					amount: big.NewInt(1000),
				},
			},
			setupMocks: func(mockTxRepo *mocks.MockTransactionRepository) {
				mockTxRepo.EXPECT().
					PersistUpdateTransactions(gomock.Any(), gomock.Any(), "1", "1").
					Do(func(ctx context.Context, txs *[]types.UpdateTransaction, fromChain, toChain string) {
						// Verify the transaction was created correctly
						assert.Len(t, *txs, 1)
						assert.Equal(t, [32]byte{1, 2, 3, 4}, (*txs)[0].ResourceId)
						assert.Equal(t, big.NewInt(1000), (*txs)[0].Amount)
						assert.Equal(t, types.Mint, (*txs)[0].UpdateType)
						assert.Equal(t, "0xabc123", (*txs)[0].TxHash)
						assert.Equal(t, "1000", (*txs)[0].BlockNumber)
						assert.Equal(t, uint8(types.AssetTypeERC721), (*txs)[0].MsgType)
						assert.Equal(t, uint64(1), (*txs)[0].LogIndex)
					}).
					Return(nil, nil)
			},
			expectedError: false,
		},
		{
			name: "success - multiple mint transactions with LogIndex increment",
			log: core.ContractLog{
				TransactionHash: "0xdef456",
				BlockNumber:     2000,
				LogIndex:        10,
			},
			resourceID:    [32]byte{5, 6, 7, 8},
			msgType:       uint8(types.AssetTypeERC721),
			issuerChainID: big.NewInt(2),
			supplies: []MintSupply{
				{
					ercID:  big.NewInt(1),
					amount: big.NewInt(1),
				},
				{
					ercID:  big.NewInt(2),
					amount: big.NewInt(1),
				},
				{
					ercID:  big.NewInt(3),
					amount: big.NewInt(1),
				},
			},
			setupMocks: func(mockTxRepo *mocks.MockTransactionRepository) {
				mockTxRepo.EXPECT().
					PersistUpdateTransactions(gomock.Any(), gomock.Any(), "2", "2").
					Do(func(ctx context.Context, txs *[]types.UpdateTransaction, fromChain, toChain string) {
						// Verify multiple transactions with incremented LogIndex
						assert.Len(t, *txs, 3)
						assert.Equal(t, uint64(10), (*txs)[0].LogIndex, "First transaction should have base LogIndex")
						assert.Equal(t, uint64(11), (*txs)[1].LogIndex, "Second transaction should have LogIndex + 1")
						assert.Equal(t, uint64(12), (*txs)[2].LogIndex, "Third transaction should have LogIndex + 2")

						// Verify ERC IDs
						assert.Equal(t, big.NewInt(1), (*txs)[0].ErcId)
						assert.Equal(t, big.NewInt(2), (*txs)[1].ErcId)
						assert.Equal(t, big.NewInt(3), (*txs)[2].ErcId)
					}).
					Return(nil, nil)
			},
			expectedError: false,
		},
		{
			name: "success - empty supplies returns without error",
			log: core.ContractLog{
				TransactionHash: "0x123",
				BlockNumber:     100,
				LogIndex:        1,
			},
			resourceID:    [32]byte{},
			msgType:       uint8(types.AssetTypeERC20),
			issuerChainID: big.NewInt(1),
			supplies:      []MintSupply{},
			setupMocks: func(mockTxRepo *mocks.MockTransactionRepository) {
				// No expectations - should not call PersistUpdateTransactions
			},
			expectedError: false,
		},
		{
			name: "error - PersistUpdateTransactions fails",
			log: core.ContractLog{
				TransactionHash: "0x789",
				BlockNumber:     3000,
				LogIndex:        15,
			},
			resourceID:    [32]byte{9, 10, 11, 12},
			msgType:       uint8(types.AssetTypeERC1155),
			issuerChainID: big.NewInt(3),
			supplies: []MintSupply{
				{
					ercID:  big.NewInt(100),
					amount: big.NewInt(500),
				},
			},
			setupMocks: func(mockTxRepo *mocks.MockTransactionRepository) {
				mockTxRepo.EXPECT().
					PersistUpdateTransactions(gomock.Any(), gomock.Any(), "3", "3").
					Return(nil, errors.New("failed to create update transaction"))
			},
			expectedError: true,
			errorContains: "failed to create update transaction",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			mockTxRepo := mocks.NewMockTransactionRepository(ctrl)

			handler := &TokenCoreEventHandler{
				txRepo: mockTxRepo,
				log:    &testutil.StubLogger{},
			}

			if tt.setupMocks != nil {
				tt.setupMocks(mockTxRepo)
			}

			// Execute
			err := handler.persistInitialMintTransactions(
				ctx,
				tt.log,
				tt.resourceID,
				tt.msgType,
				tt.issuerChainID,
				tt.supplies,
			)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
