package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/core"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/events"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/mocks"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/testutil"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/contracts/EnygmaTeleport"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/types"
)

func TestEnygmaTeleportEventHandler_Name(t *testing.T) {
	handler := NewEnygmaTeleportEventHandler(nil, nil, nil, nil, nil)

	name := handler.Name()

	assert.Equal(t, "EnygmaTeleportHandler", name)
}

// Helper to create a test block
func createTestBlock(timestamp uint64) *ethTypes.Block {
	header := &ethTypes.Header{
		Time: timestamp,
	}
	return ethTypes.NewBlock(header, &ethTypes.Body{}, nil, nil)
}

// Helper to create EnygmaTransferBatch with specified number of transactions
func createEnygmaTransferBatch(numTxs int) types.EnygmaTransferBatch {
	txs := make([]*types.EnygmaTransferBatchTx, numTxs)
	for i := 0; i < numTxs; i++ {
		txs[i] = &types.EnygmaTransferBatchTx{
			MessageId:   "",
			ReferenceId: [32]byte{},
			FromAddress: common.Address{},
			ToAmount:    big.NewInt(1000),
			ToAddress:   common.Address{},
		}
	}
	return types.EnygmaTransferBatch{
		ResourceId:     "abcd1234",
		HubBlockNumber: big.NewInt(1000),
		FromChainID:    big.NewInt(1),
		ToChainID:      big.NewInt(2),
		ToRValueToAdd:  big.NewInt(999),
		Transactions:   txs,
		BatchId:        "",
	}
}

func TestEnygmaTeleportEventHandler_Handle(t *testing.T) {
	tests := []struct {
		name          string
		log           core.ContractLog
		setupMocks    func(*mocks.MockTransactionRepository, *mocks.MockDecryptor, *mocks.MockProvider)
		expectedError bool
		errorContains string
	}{
		{
			name: "routes EnygmaTransfer to handler",
			log: core.ContractLog{
				ContractName:    events.ContractEnygmaTeleport,
				EventName:       events.EnygmaTransfer,
				BlockNumber:     100,
				TransactionHash: "0xabc123",
				RawEventData: testutil.MustMarshal(t, &EnygmaTeleport.EnygmaTeleportEnygmaTransfer{
					ResourceId:       [32]byte{1, 2, 3},
					EncryptedMessage: []byte{},
				}),
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, decryptor *mocks.MockDecryptor, provider *mocks.MockProvider) {
			},
			expectedError: false,
		},
		{
			name: "routes EnygmaTransferCompleted to handler",
			log: core.ContractLog{
				ContractName:    events.ContractEnygmaTeleport,
				EventName:       events.EnygmaTransferCompleted,
				BlockNumber:     100,
				TransactionHash: "0xabc123",
				RawEventData: testutil.MustMarshal(t, &EnygmaTeleport.EnygmaTeleportEnygmaTransferCompleted{
					EncryptedMessage: []byte{},
				}),
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, decryptor *mocks.MockDecryptor, provider *mocks.MockProvider) {
			},
			expectedError: false,
		},
		{
			name: "routes EnygmaSupplyUpdated to handler",
			log: core.ContractLog{
				ContractName:    events.ContractEnygmaTeleport,
				EventName:       events.EnygmaSupplyUpdated,
				BlockNumber:     100,
				TransactionHash: "0xabc123",
				RawEventData: testutil.MustMarshal(t, &EnygmaTeleport.EnygmaTeleportEnygmaSupplyUpdated{
					ResourceId:  [32]byte{1, 2, 3},
					BlockNumber: big.NewInt(100),
					Update: struct {
						Amount *big.Int
						TxType uint8
					}{
						Amount: big.NewInt(1000),
						TxType: uint8(types.EnygmaTxTypeMint),
					},
					ChainId: big.NewInt(1),
				}),
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, decryptor *mocks.MockDecryptor, provider *mocks.MockProvider) {
				// Mock transaction creation
				txRepo.EXPECT().
					CreateTransaction(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedError: false,
		},
		{
			name: "don't throw error if it is an unknown event",
			log: core.ContractLog{
				ContractName: events.ContractEnygmaTeleport,
				EventName:    "UnknownEvent",
				RawEventData: nil,
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, decryptor *mocks.MockDecryptor, provider *mocks.MockProvider) {
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTxRepo := mocks.NewMockTransactionRepository(ctrl)
			mockDecryptor := mocks.NewMockDecryptor(ctrl)
			mockProvider := mocks.NewMockProvider(ctrl)
			logger := &testutil.StubLogger{}
			pnData := make(core.PNodeDataAndSecrets)

			handler := NewEnygmaTeleportEventHandler(
				mockTxRepo,
				mockDecryptor,
				mockProvider,
				logger,
				&pnData,
			)

			if tt.setupMocks != nil {
				tt.setupMocks(mockTxRepo, mockDecryptor, mockProvider)
			}

			err := handler.Handle(context.Background(), tt.log)

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

func TestEnygmaTeleportEventHandler_ProcessEnygmaTransferEvent(t *testing.T) {
	tests := []struct {
		name          string
		log           core.ContractLog
		setupMocks    func(*mocks.MockTransactionRepository, *mocks.MockDecryptor, *mocks.MockProvider)
		expectedError string
	}{
		{
			name: "error - type assertion failure",
			log: core.ContractLog{
				ContractName: events.ContractEnygmaTeleport,
				EventName:    events.EnygmaTransfer,
				RawEventData: nil,
				BlockNumber:  100,
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, decryptor *mocks.MockDecryptor, provider *mocks.MockProvider) {
			},
			expectedError: "failed to unmarshal event data",
		},
		{
			name: "success - empty encrypted message returns nil",
			log: core.ContractLog{
				ContractName: events.ContractEnygmaTeleport,
				EventName:    events.EnygmaTransfer,
				RawEventData: testutil.MustMarshal(t, &EnygmaTeleport.EnygmaTeleportEnygmaTransfer{
					ResourceId:       [32]byte{1, 2, 3},
					EncryptedMessage: []byte{},
				}),
				BlockNumber: 100,
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, decryptor *mocks.MockDecryptor, provider *mocks.MockProvider) {
			},
			expectedError: "",
		},
		{
			name: "error - CreateTransactionsWithEnygmaData failure",
			log: core.ContractLog{
				ContractName: events.ContractEnygmaTeleport,
				EventName:    events.EnygmaTransfer,
				RawEventData: testutil.MustMarshal(t, &EnygmaTeleport.EnygmaTeleportEnygmaTransfer{
					ResourceId:       [32]byte{1, 2, 3},
					EncryptedMessage: []byte("encrypted-data"),
				}),
				BlockNumber:     100,
				TransactionHash: "0xabc123",
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, decryptor *mocks.MockDecryptor, provider *mocks.MockProvider) {
				batch := createEnygmaTransferBatch(1)
				batchJSON, _ := json.Marshal(batch)

				decryptor.EXPECT().
					DecryptPayloadBytes(gomock.Any(), uint64(100), gomock.Any(), types.ParticipantSecret).
					Return(batchJSON, nil)

				// CreateTransactionsWithEnygmaData fails
				txRepo.EXPECT().
					CreateTransactionsWithEnygmaData(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("db error"))
			},
			expectedError: "db error",
		},
		{
			name: "returns nil on empty transactions after decryption",
			log: core.ContractLog{
				ContractName: events.ContractEnygmaTeleport,
				EventName:    events.EnygmaTransfer,
				RawEventData: testutil.MustMarshal(t, &EnygmaTeleport.EnygmaTeleportEnygmaTransfer{
					ResourceId:       [32]byte{1, 2, 3},
					EncryptedMessage: []byte("encrypted-data"),
				}),
				BlockNumber:     100,
				TransactionHash: "0xabc123",
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, decryptor *mocks.MockDecryptor, provider *mocks.MockProvider) {
				// Return batch without transactions
				batch := createEnygmaTransferBatch(0)
				batchJSON, _ := json.Marshal(batch)

				decryptor.EXPECT().
					DecryptPayloadBytes(gomock.Any(), uint64(100), gomock.Any(), types.ParticipantSecret).
					Return(batchJSON, nil)
			},
			expectedError: "",
		},
		{
			name: "verifies equal length arrays are passed to CreateTransactionsWithEnygmaData",
			log: core.ContractLog{
				ContractName: events.ContractEnygmaTeleport,
				EventName:    events.EnygmaTransfer,
				RawEventData: testutil.MustMarshal(t, &EnygmaTeleport.EnygmaTeleportEnygmaTransfer{
					ResourceId:       [32]byte{1, 2, 3},
					EncryptedMessage: []byte("encrypted-data"),
				}),
				BlockNumber:     100,
				TransactionHash: "0xabc123",
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, decryptor *mocks.MockDecryptor, provider *mocks.MockProvider) {
				batch := createEnygmaTransferBatch(3)
				batchJSON, _ := json.Marshal(batch)

				decryptor.EXPECT().
					DecryptPayloadBytes(gomock.Any(), uint64(100), gomock.Any(), types.ParticipantSecret).
					Return(batchJSON, nil)

				// Verify that CreateTransactionsWithEnygmaData is called with equal length arrays
				txRepo.EXPECT().
					CreateTransactionsWithEnygmaData(
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
					).
					Do(func(ctx context.Context, transactions []*domain.Transaction, enygmaTransactions []*domain.EnygmaTransaction) {
						assert.Equal(t, len(transactions), len(enygmaTransactions),
							"transactions and enygmaTransactions arrays must have equal length")
						assert.Equal(t, 3, len(transactions), "expected 3 transactions")
						assert.Equal(t, 3, len(enygmaTransactions), "expected 3 enygma transactions")
						pendingStatus := uint8(types.EnygmaTransferPending)
						for _, tx := range transactions {
							assert.Equal(t, &pendingStatus, tx.TeleportStatus, "expected TeleportStatus to be Pending")
						}
					}).
					Return(nil)
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTxRepo := mocks.NewMockTransactionRepository(ctrl)
			mockDecryptor := mocks.NewMockDecryptor(ctrl)
			mockProvider := mocks.NewMockProvider(ctrl)
			logger := &testutil.StubLogger{}
			pnData := make(core.PNodeDataAndSecrets)

			handler := NewEnygmaTeleportEventHandler(
				mockTxRepo,
				mockDecryptor,
				mockProvider,
				logger,
				&pnData,
			)

			if tt.setupMocks != nil {
				tt.setupMocks(mockTxRepo, mockDecryptor, mockProvider)
			}

			err := handler.processEnygmaTransferEvent(context.Background(), tt.log)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEnygmaTeleportEventHandler_ProcessEnygmaTransferCompletedEvent(t *testing.T) {
	tests := []struct {
		name          string
		log           core.ContractLog
		setupMocks    func(*mocks.MockTransactionRepository, *mocks.MockDecryptor, *mocks.MockProvider)
		expectedError string
	}{
		{
			name: "error - type assertion failure",
			log: core.ContractLog{
				ContractName: events.ContractEnygmaTeleport,
				EventName:    events.EnygmaTransferCompleted,
				RawEventData: nil,
				BlockNumber:  100,
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, decryptor *mocks.MockDecryptor, provider *mocks.MockProvider) {
			},
			expectedError: "failed to unmarshal event data",
		},
		{
			name: "empty encrypted message returns nil",
			log: core.ContractLog{
				ContractName: events.ContractEnygmaTeleport,
				EventName:    events.EnygmaTransferCompleted,
				RawEventData: testutil.MustMarshal(t, &EnygmaTeleport.EnygmaTeleportEnygmaTransferCompleted{
					EncryptedMessage: []byte{},
				}),
				BlockNumber: 100,
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, decryptor *mocks.MockDecryptor, provider *mocks.MockProvider) {
			},
			expectedError: "",
		},
		{
			name: "empty decrypted data returns nil",
			log: core.ContractLog{
				ContractName: events.ContractEnygmaTeleport,
				EventName:    events.EnygmaTransferCompleted,
				RawEventData: testutil.MustMarshal(t, &EnygmaTeleport.EnygmaTeleportEnygmaTransferCompleted{
					EncryptedMessage: []byte("encrypted-data"),
				}),
				BlockNumber: 100,
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, decryptor *mocks.MockDecryptor, provider *mocks.MockProvider) {
				emptyData := []types.EnygmaTransferCompleted{}
				emptyJSON, _ := json.Marshal(emptyData)

				decryptor.EXPECT().
					DecryptPayloadBytes(gomock.Any(), uint64(100), gomock.Any(), types.ParticipantSecret).
					Return(emptyJSON, nil)
			},
			expectedError: "",
		},
		{
			name: "transaction not found in map",
			log: core.ContractLog{
				ContractName: events.ContractEnygmaTeleport,
				EventName:    events.EnygmaTransferCompleted,
				RawEventData: testutil.MustMarshal(t, &EnygmaTeleport.EnygmaTeleportEnygmaTransferCompleted{
					EncryptedMessage: []byte("encrypted-data"),
				}),
				BlockNumber: 100,
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, decryptor *mocks.MockDecryptor, provider *mocks.MockProvider) {
				data := []types.EnygmaTransferCompleted{
					{
						MessageId:       "msg-123",
						TransactionHash: "0xabc",
					},
				}
				jsonData, _ := json.Marshal(data)

				decryptor.EXPECT().
					DecryptPayloadBytes(gomock.Any(), uint64(100), gomock.Any(), types.ParticipantSecret).
					Return(jsonData, nil)

				// Return empty slice - transaction not found
				txRepo.EXPECT().
					GetTransactionsByMessageIDs(gomock.Any(), []string{"msg-123"}, false).
					Return([]domain.Transaction{}, nil)

				// UpdateTransactionsBulk should be called with empty slice
				txRepo.EXPECT().
					UpdateTransactionsBulk(gomock.Any(), gomock.Any(), "tx_hash_destination", "teleport_status", "updated_at").
					Return(nil)
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTxRepo := mocks.NewMockTransactionRepository(ctrl)
			mockDecryptor := mocks.NewMockDecryptor(ctrl)
			mockProvider := mocks.NewMockProvider(ctrl)
			logger := &testutil.StubLogger{}
			pnData := make(core.PNodeDataAndSecrets)

			handler := NewEnygmaTeleportEventHandler(
				mockTxRepo,
				mockDecryptor,
				mockProvider,
				logger,
				&pnData,
			)

			if tt.setupMocks != nil {
				tt.setupMocks(mockTxRepo, mockDecryptor, mockProvider)
			}

			err := handler.processEnygmaTransferCompletedEvent(context.Background(), tt.log)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEnygmaTeleportEventHandler_ProcessEnygmaSupplyUpdatedEvent(t *testing.T) {
	tests := []struct {
		name          string
		log           core.ContractLog
		setupMocks    func(*mocks.MockTransactionRepository, *mocks.MockDecryptor, *mocks.MockProvider)
		expectedError string
	}{
		{
			name: "error - type assertion failure",
			log: core.ContractLog{
				ContractName: events.ContractEnygmaTeleport,
				EventName:    events.EnygmaSupplyUpdated,
				RawEventData: nil,
				BlockNumber:  100,
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, decryptor *mocks.MockDecryptor, provider *mocks.MockProvider) {
			},
			expectedError: "failed to unmarshal event data",
		},
		{
			name: "success - mint transaction",
			log: core.ContractLog{
				ContractName:    events.ContractEnygmaTeleport,
				EventName:       events.EnygmaSupplyUpdated,
				BlockNumber:     100,
				TransactionHash: "0xabc123",
				RawEventData: testutil.MustMarshal(t, &EnygmaTeleport.EnygmaTeleportEnygmaSupplyUpdated{
					ResourceId:  [32]byte{1, 2, 3, 4},
					BlockNumber: big.NewInt(100),
					Update: struct {
						Amount *big.Int
						TxType uint8
					}{
						Amount: big.NewInt(1000),
						TxType: uint8(types.EnygmaTxTypeMint),
					},
					ChainId: big.NewInt(1),
				}),
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, decryptor *mocks.MockDecryptor, provider *mocks.MockProvider) {
				txRepo.EXPECT().
					CreateTransaction(gomock.Any(), gomock.Any()).
					Do(func(ctx context.Context, tx *domain.Transaction) {
						// Verify business rule: Mint transactions have positive amounts
						assert.True(t, tx.Amount.IsPositive(), "Mint transaction amount should be positive")
						assert.Equal(t, types.Mint, tx.TxType, "Transaction type should be Mint")
					}).
					Return(nil)
			},
			expectedError: "",
		},
		{
			name: "success - burn transaction (negative amount)",
			log: core.ContractLog{
				ContractName:    events.ContractEnygmaTeleport,
				EventName:       events.EnygmaSupplyUpdated,
				BlockNumber:     100,
				TransactionHash: "0xdef456",
				RawEventData: testutil.MustMarshal(t, &EnygmaTeleport.EnygmaTeleportEnygmaSupplyUpdated{
					ResourceId:  [32]byte{5, 6, 7, 8},
					BlockNumber: big.NewInt(200),
					Update: struct {
						Amount *big.Int
						TxType uint8
					}{
						Amount: big.NewInt(2000),
						TxType: uint8(types.EnygmaTxTypeBurn),
					},
					ChainId: big.NewInt(2),
				}),
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, decryptor *mocks.MockDecryptor, provider *mocks.MockProvider) {
				txRepo.EXPECT().
					CreateTransaction(gomock.Any(), gomock.Any()).
					Do(func(ctx context.Context, tx *domain.Transaction) {
						// Verify business rule: Burn transactions have negative amounts
						assert.True(t, tx.Amount.IsNegative(), "Burn transaction amount should be negative")
						assert.Equal(t, types.Burn, tx.TxType, "Transaction type should be Burn")
					}).
					Return(nil)
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTxRepo := mocks.NewMockTransactionRepository(ctrl)
			mockDecryptor := mocks.NewMockDecryptor(ctrl)
			mockProvider := mocks.NewMockProvider(ctrl)
			logger := &testutil.StubLogger{}
			pnData := make(core.PNodeDataAndSecrets)

			handler := NewEnygmaTeleportEventHandler(
				mockTxRepo,
				mockDecryptor,
				mockProvider,
				logger,
				&pnData,
			)

			if tt.setupMocks != nil {
				tt.setupMocks(mockTxRepo, mockDecryptor, mockProvider)
			}

			err := handler.processEnygmaSupplyUpdatedEvent(context.Background(), tt.log)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEnygmaTeleportEventHandler_ProcessEnygmaDvpBalanceUpdated(t *testing.T) {
	tests := []struct {
		name          string
		log           core.ContractLog
		setupMocks    func(*mocks.MockTransactionRepository, *mocks.MockDecryptor, *mocks.MockProvider)
		expectedError bool
		errorContains string
	}{
		{
			name: "success - processes EnygmaDvpBalanceUpdated event",
			log: core.ContractLog{
				RawEventData: testutil.MustMarshal(t, &EnygmaTeleport.EnygmaTeleportEnygmaDvpBalanceUpdated{
					EncryptedMessage: []byte("encrypted-data"),
				}),
				BlockNumber:     1000,
				TransactionHash: "0xabc123",
				LogIndex:        1,
			},
			setupMocks: func(mockTxRepo *mocks.MockTransactionRepository, mockDecryptor *mocks.MockDecryptor, mockProvider *mocks.MockProvider) {
				// Decryption succeeds
				mockDecryptor.EXPECT().
					DecryptPayloadBytes(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func([]byte, uint64, interface{}, types.SecretType) ([]byte, error) {
						data := types.DvpBalanceUpdated{
							ErcId:              big.NewInt(123),
							TokenType:          uint8(types.AssetTypeERC721),
							ResourceId:         "0x1234",
							Amount:             big.NewInt(1000),
							UpdateType:         types.Mint,
							SourceChainId:      big.NewInt(1),
							DestinationChainId: big.NewInt(2),
						}
						return json.Marshal(data)
					})

				// Transaction creation succeeds
				mockTxRepo.EXPECT().
					CreateTransaction(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedError: false,
		},
		{
			name: "error - type assertion fails",
			log: core.ContractLog{
				RawEventData: nil,
			},
			setupMocks: func(mockTxRepo *mocks.MockTransactionRepository, mockDecryptor *mocks.MockDecryptor, mockProvider *mocks.MockProvider) {
			},
			expectedError: true,
			errorContains: "failed to unmarshal event data",
		},
		{
			name: "error - decryption fails",
			log: core.ContractLog{
				RawEventData: testutil.MustMarshal(t, &EnygmaTeleport.EnygmaTeleportEnygmaDvpBalanceUpdated{
					EncryptedMessage: []byte("encrypted-data"),
				}),
				BlockNumber:     1000,
				TransactionHash: "0xabc123",
				LogIndex:        1,
			},
			setupMocks: func(mockTxRepo *mocks.MockTransactionRepository, mockDecryptor *mocks.MockDecryptor, mockProvider *mocks.MockProvider) {
				// Decryption fails
				mockDecryptor.EXPECT().
					DecryptPayloadBytes(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("decryption failed"))
			},
			expectedError: true,
			errorContains: "decryption failed",
		},
		{
			name: "error - CreateTransaction fails",
			log: core.ContractLog{
				RawEventData: testutil.MustMarshal(t, &EnygmaTeleport.EnygmaTeleportEnygmaDvpBalanceUpdated{
					EncryptedMessage: []byte("encrypted-data"),
				}),
				BlockNumber:     1000,
				TransactionHash: "0xabc123",
				LogIndex:        1,
			},
			setupMocks: func(mockTxRepo *mocks.MockTransactionRepository, mockDecryptor *mocks.MockDecryptor, mockProvider *mocks.MockProvider) {
				// Decryption succeeds
				mockDecryptor.EXPECT().
					DecryptPayloadBytes(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func([]byte, uint64, interface{}, types.SecretType) ([]byte, error) {
						data := types.DvpBalanceUpdated{
							ErcId:              big.NewInt(123),
							TokenType:          uint8(types.AssetTypeERC721),
							ResourceId:         "0x1234",
							Amount:             big.NewInt(1000),
							UpdateType:         types.Mint,
							SourceChainId:      big.NewInt(1),
							DestinationChainId: big.NewInt(2),
						}
						return json.Marshal(data)
					})

				// Transaction creation fails
				mockTxRepo.EXPECT().
					CreateTransaction(gomock.Any(), gomock.Any()).
					Return(errors.New("failed to create transaction"))
			},
			expectedError: true,
			errorContains: "failed to create transaction",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			mockTxRepo := mocks.NewMockTransactionRepository(ctrl)
			mockDecryptor := mocks.NewMockDecryptor(ctrl)
			mockProvider := mocks.NewMockProvider(ctrl)

			handler := &EnygmaTeleportEventHandler{
				txRepo:    mockTxRepo,
				decryptor: mockDecryptor,
				provider:  mockProvider,
				log:       &testutil.StubLogger{},
				pnData:    &core.PNodeDataAndSecrets{},
			}

			if tt.setupMocks != nil {
				tt.setupMocks(mockTxRepo, mockDecryptor, mockProvider)
			}

			// Execute
			err := handler.processEnygmaDvpBalanceUpdated(ctx, tt.log)

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
