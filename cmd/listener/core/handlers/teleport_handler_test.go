package handlers

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/core"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/events"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/mocks"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/testutil"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/TeleportV1"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/domain"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/types"
)

// TestNewTeleportEventHandler tests the constructor
func TestNewTeleportEventHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTxRepo := mocks.NewMockTransactionRepository(ctrl)
	mockRevertRepo := mocks.NewMockRevertDataTransactionRepository(ctrl)
	mockDecryptor := mocks.NewMockDecryptor(ctrl)
	logger := &testutil.StubLogger{}
	pnData := &core.PNodeDataAndSecrets{}

	handler := NewTeleportEventHandler(mockTxRepo, mockRevertRepo, mockDecryptor, logger, pnData)

	assert.NotNil(t, handler)
	assert.Equal(t, mockTxRepo, handler.txRepo)
	assert.Equal(t, mockRevertRepo, handler.revertRepo)
	assert.Equal(t, mockDecryptor, handler.decryptor)
	assert.Equal(t, logger, handler.log)
	assert.Equal(t, pnData, handler.pnData)
}

// TestTeleportEventHandler_Handle tests the event routing
func TestTeleportEventHandler_Handle(t *testing.T) {
	tests := []struct {
		name          string
		log           core.ContractLog
		setupMocks    func(*mocks.MockTransactionRepository, *mocks.MockRevertDataTransactionRepository, *mocks.MockDecryptor)
		expectedError bool
		errorContains string
	}{
		{
			name: "routes AtomicMessageAdditionalDataBatch to handler",
			log: core.ContractLog{
				ContractName: events.ContractTeleport,
				EventName:    events.AtomicMessageAdditionalDataBatch,
				RawEventData: testutil.MustMarshal(t, &TeleportV1.TeleportV1AtomicMessageAdditionalDataBatch{
					EncryptedData: "",
				}),
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, revertRepo *mocks.MockRevertDataTransactionRepository, decryptor *mocks.MockDecryptor) {
			},
			expectedError: false,
		},
		{
			name: "routes AtomicMessageStatusChangedBatch to handler",
			log: core.ContractLog{
				ContractName: events.ContractTeleport,
				EventName:    events.AtomicMessageStatusChangedBatch,
				RawEventData: testutil.MustMarshal(t, &TeleportV1.TeleportV1AtomicMessageStatusChangedBatch{
					MsgIds: []string{},
				}),
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, revertRepo *mocks.MockRevertDataTransactionRepository, decryptor *mocks.MockDecryptor) {
			},
			expectedError: false,
		},
		{
			name: "routes EncryptedDataBatchStored to handler",
			log: core.ContractLog{
				ContractName: events.ContractTeleport,
				EventName:    events.EncryptedDataBatchStored,
				RawEventData: testutil.MustMarshal(t, &TeleportV1.TeleportV1EncryptedDataBatchStored{
					Data: []byte{},
				}),
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, revertRepo *mocks.MockRevertDataTransactionRepository, decryptor *mocks.MockDecryptor) {
				decryptor.EXPECT().
					DecryptPayloadBytes(gomock.Any(), gomock.Any(), gomock.Any(), types.ParticipantSecret).
					Return([]byte("[]"), nil)
			},
			expectedError: false,
		},
		{
			name: "don't throw error if it is an unknown event",
			log: core.ContractLog{
				ContractName: events.ContractTeleport,
				EventName:    "UnknownEvent",
				RawEventData: nil,
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, revertRepo *mocks.MockRevertDataTransactionRepository, decryptor *mocks.MockDecryptor) {
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTxRepo := mocks.NewMockTransactionRepository(ctrl)
			mockRevertRepo := mocks.NewMockRevertDataTransactionRepository(ctrl)
			mockDecryptor := mocks.NewMockDecryptor(ctrl)
			logger := &testutil.StubLogger{}
			pnData := &core.PNodeDataAndSecrets{}

			handler := NewTeleportEventHandler(mockTxRepo, mockRevertRepo, mockDecryptor, logger, pnData)

			if tt.setupMocks != nil {
				tt.setupMocks(mockTxRepo, mockRevertRepo, mockDecryptor)
			}

			err := handler.Handle(context.Background(), tt.log)

			assert.NoError(t, err)
		})
	}
}

// TestConvertToTransaction tests the conversion from blockchain event to domain transaction
func TestConvertToTransaction(t *testing.T) {
	// Helper to create a message
	buildMessage := func() *types.DispatchedMessageToPrivateHub {
		return &types.DispatchedMessageToPrivateHub{
			MessageId: [32]byte{
				1,
				2,
				3,
				4,
				5,
				6,
				7,
				8,
				9,
				10,
				11,
				12,
				13,
				14,
				15,
				16,
				17,
				18,
				19,
				20,
				21,
				22,
				23,
				24,
				25,
				26,
				27,
				28,
				29,
				30,
				31,
				32,
			},
			From:        common.HexToAddress("0x1234567890123456789012345678901234567890"),
			To:          common.HexToAddress("0x0987654321098765432109876543210987654321"),
			ToChainId:   big.NewInt(1),
			FromChainId: big.NewInt(2),
			SharedId:    "shared-123",
			BatchId:     "batch-456",
			IsAtomic:    false,
			LogIdx:      5,
			Data: types.RaylzMessage{
				MessageMetadata: types.RaylzMessageMetadata{
					ResourceId: [32]byte{
						10,
						20,
						30,
						40,
						50,
						60,
						70,
						80,
						90,
						100,
						110,
						120,
						130,
						140,
						150,
						160,
						170,
						180,
						190,
						200,
						210,
						220,
						230,
						240,
						250,
						255,
						254,
						253,
						252,
						251,
						250,
						249,
					},
					TransferMetadata: types.BridgedTransferMetadata{
						Amount:    big.NewInt(1000000000000000000),
						AssetType: types.AssetTypeERC20,
						From:      "0x1234567890123456789012345678901234567890",
						To:        "0x0987654321098765432109876543210987654321",
						Id:        big.NewInt(0),
					},
				},
				Payload: []byte("test payload"),
			},
			TxHashDestination: common.HexToHash(
				"0x1230000000000000000000000000000000000000000000000000000000000456",
			),
		}
	}

	// Helper to validate invariant tx fields
	validateInvariants := func(t *testing.T, tx *domain.Transaction) {
		// Chain ids
		assert.Equal(t, "1", tx.ToChainId.Unwrap().String())
		assert.Equal(t, "2", tx.FromChainId.Unwrap().String())

		// Hex-encoded fields
		assert.Equal(t, "0a141e28323c46505a646e78828c96a0aab4bec8d2dce6f0fafffefdfcfbfaf9", tx.ResourceId)
		assert.Equal(t, "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20", tx.MessageId)

		// Hash conversions
		assert.Equal(t, "0x00000000000000000000000000000000000000000000000000000000000cc123", tx.HubTxHash)
		assert.Equal(t, "0x1230000000000000000000000000000000000000000000000000000000000456", tx.TxHashDestination)

		// Addresses
		assert.Equal(t, "0x1234567890123456789012345678901234567890", tx.From)
		assert.Equal(t, "0x0987654321098765432109876543210987654321", tx.To)

		// Strings
		assert.Equal(t, "shared-123", tx.SharedId)
		assert.Equal(t, "batch-456", tx.BatchId)

		// Numeric/Booleans
		assert.Equal(t, uint64(5), tx.LogIndex)
		assert.Equal(t, types.CrossChain, tx.TxType)

		// TeleportStatus validation depends on protocol
		if tx.Protocol == types.Vanilla {
			assert.Nil(t, tx.TeleportStatus, "Vanilla protocol should have nil TeleportStatus")
		} else {
			assert.NotNil(t, tx.TeleportStatus, "Non-Vanilla protocol should have non-nil TeleportStatus")
			assert.Equal(t, uint8(types.AtomicTeleportPending), *tx.TeleportStatus)
		}

		assert.NotEmpty(t, tx.Payload)
	}

	tests := []struct {
		name          string
		msg           *types.DispatchedMessageToPrivateHub
		blockNumber   string
		logHash       string
		expectedError bool
		errorContains string
		validate      func(t *testing.T, tx *domain.Transaction)
	}{
		{
			name:          "nil message returns error",
			msg:           nil,
			blockNumber:   "100",
			expectedError: true,
			errorContains: "message is nil",
		},
		{
			name: "nil amount defaults to zero",
			msg: func() *types.DispatchedMessageToPrivateHub {
				msg := buildMessage()
				msg.Data.MessageMetadata.TransferMetadata.Amount = nil
				return msg
			}(),
			blockNumber: "100",
			logHash:     common.HexToHash("0xcc123").String(),
			validate: func(t *testing.T, tx *domain.Transaction) {
				assert.Equal(t, "100", tx.BlockNumber.String())
				assert.Nil(t, tx.ErcId.Unwrap())
				assert.True(t, tx.Amount.IsZero())
				validateInvariants(t, tx)
			},
		},
		{
			name:          "invalid block number returns error",
			msg:           buildMessage(),
			blockNumber:   "abc",
			expectedError: true,
			errorContains: "can't convert",
		},
		{
			name:        "valid ERC20 transaction",
			msg:         buildMessage(),
			blockNumber: "100",
			logHash:     common.HexToHash("0xcc123").String(),
			validate: func(t *testing.T, tx *domain.Transaction) {
				// Validate variant fields
				assert.Equal(t, "100", tx.BlockNumber.String())
				assert.Equal(t, "1000000000000000000", tx.Amount.String())
				assert.Equal(t, uint8(types.AssetTypeERC20), tx.MsgType)
				assert.Nil(t, tx.ErcId.Unwrap())
				// Validate invariants
				validateInvariants(t, tx)
			},
		},
		{
			name: "valid ERC721 transaction with ercId",
			msg: func() *types.DispatchedMessageToPrivateHub {
				msg := buildMessage()
				msg.Data.MessageMetadata.TransferMetadata.AssetType = types.AssetTypeERC721
				msg.Data.MessageMetadata.TransferMetadata.Id = big.NewInt(42)
				return msg
			}(),
			blockNumber: "100",
			logHash:     common.HexToHash("0xcc123").String(),
			validate: func(t *testing.T, tx *domain.Transaction) {
				// Validate variant fields
				assert.Equal(t, "100", tx.BlockNumber.String())
				assert.Equal(t, "1000000000000000000", tx.Amount.String())
				assert.Equal(t, uint8(types.AssetTypeERC721), tx.MsgType)
				require.NotNil(t, tx.ErcId.Unwrap())
				assert.Equal(t, "42", tx.ErcId.Unwrap().String())
				// Validate invariants
				validateInvariants(t, tx)
			},
		},
		{
			name: "valid ERC1155 transaction with ercId",
			msg: func() *types.DispatchedMessageToPrivateHub {
				msg := buildMessage()
				msg.Data.MessageMetadata.TransferMetadata.AssetType = types.AssetTypeERC1155
				msg.Data.MessageMetadata.TransferMetadata.Id = big.NewInt(999)
				return msg
			}(),
			blockNumber: "100",
			logHash:     common.HexToHash("0xcc123").String(),
			validate: func(t *testing.T, tx *domain.Transaction) {
				// Validate variant fields
				assert.Equal(t, "100", tx.BlockNumber.String())
				assert.Equal(t, "1000000000000000000", tx.Amount.String())
				assert.Equal(t, uint8(types.AssetTypeERC1155), tx.MsgType)
				require.NotNil(t, tx.ErcId.Unwrap())
				assert.Equal(t, "999", tx.ErcId.Unwrap().String())
				// Validate invariants
				validateInvariants(t, tx)
			},
		},
		{
			name: "valid Enygma transaction",
			msg: func() *types.DispatchedMessageToPrivateHub {
				msg := buildMessage()
				msg.Data.MessageMetadata.TransferMetadata.AssetType = types.AssetTypeEnygma
				return msg
			}(),
			blockNumber: "100",
			logHash:     common.HexToHash("0xcc123").String(),
			validate: func(t *testing.T, tx *domain.Transaction) {
				// Validate variant fields
				assert.Equal(t, "100", tx.BlockNumber.String())
				assert.Equal(t, "1000000000000000000", tx.Amount.String())
				assert.Equal(t, uint8(types.AssetTypeEnygma), tx.MsgType)
				assert.Nil(t, tx.ErcId.Unwrap())
				// Validate invariants
				validateInvariants(t, tx)
			},
		},
		{
			name: "atomic transaction",
			msg: func() *types.DispatchedMessageToPrivateHub {
				msg := buildMessage()
				msg.IsAtomic = true
				return msg
			}(),
			blockNumber: "100",
			logHash:     common.HexToHash("0xcc123").String(),
			validate: func(t *testing.T, tx *domain.Transaction) {
				// Validate variant fields
				assert.Equal(t, "100", tx.BlockNumber.String())
				assert.Equal(t, "1000000000000000000", tx.Amount.String())
				assert.Nil(t, tx.ErcId.Unwrap())
				assert.Equal(t, types.Atomic, tx.Protocol)
				// Validate invariants
				validateInvariants(t, tx)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &TeleportEventHandler{}

			tx, err := handler.convertToTransaction(tt.msg, tt.blockNumber, tt.logHash, time.Time{})

			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				assert.Nil(t, tx)
			} else {
				require.NoError(t, err)
				require.NotNil(t, tx)
				if tt.validate != nil {
					tt.validate(t, tx)
				}
			}
		})
	}
}

// TestApplyAtomicAdditionalData tests applying atomic teleport additional data to transactions
func TestBuildRevertDataIfPresent(t *testing.T) {
	tests := []struct {
		name           string
		data           types.AtomicTeleportAdditionalData
		expectNil      bool
		validateRevert func(t *testing.T, revert *domain.RevertDataTransaction)
	}{
		{
			name:      "returns nil when no revert data",
			data:      types.AtomicTeleportAdditionalData{},
			expectNil: true,
		},
		{
			name: "destination revert hash only",
			data: types.AtomicTeleportAdditionalData{
				TxHashDestinationRevert: common.HexToHash("0x123"),
			},
			expectNil: false,
			validateRevert: func(t *testing.T, revert *domain.RevertDataTransaction) {
				assert.Equal(
					t,
					"0x0000000000000000000000000000000000000000000000000000000000000123",
					revert.TxHashDestinationRevert,
				)
				assert.Equal(t, uint8(0), revert.TxHashDestinationRevertStatus)
				assert.True(t, revert.TxHashSourceRevert == "")
				assert.Equal(t, uint8(0), revert.TxHashSourceRevertStatus)
			},
		},
		{
			name: "destination revert with status",
			data: types.AtomicTeleportAdditionalData{
				TxHashDestinationRevert:       common.HexToHash("0x789"),
				TxHashDestinationRevertStatus: 1,
			},
			expectNil: false,
			validateRevert: func(t *testing.T, revert *domain.RevertDataTransaction) {
				assert.Equal(
					t,
					"0x0000000000000000000000000000000000000000000000000000000000000789",
					revert.TxHashDestinationRevert,
				)
				assert.Equal(t, uint8(1), revert.TxHashDestinationRevertStatus)
			},
		},
		{
			name: "destination revert all fields",
			data: types.AtomicTeleportAdditionalData{
				TxHashDestinationRevert:       common.HexToHash("0xaaa"),
				TxHashDestinationRevertStatus: 1,
			},
			expectNil: false,
			validateRevert: func(t *testing.T, revert *domain.RevertDataTransaction) {
				assert.Equal(
					t,
					"0x0000000000000000000000000000000000000000000000000000000000000aaa",
					revert.TxHashDestinationRevert,
				)
				assert.Equal(t, uint8(1), revert.TxHashDestinationRevertStatus)
			},
		},
		{
			name: "source revert hash only",
			data: types.AtomicTeleportAdditionalData{
				TxHashSourceRevert: common.HexToHash("0xabc"),
			},
			expectNil: false,
			validateRevert: func(t *testing.T, revert *domain.RevertDataTransaction) {
				assert.Equal(
					t,
					"0x0000000000000000000000000000000000000000000000000000000000000abc",
					revert.TxHashSourceRevert,
				)
				assert.Equal(t, uint8(0), revert.TxHashSourceRevertStatus)
			},
		},
		{
			name: "source revert with status",
			data: types.AtomicTeleportAdditionalData{
				TxHashSourceRevert:       common.HexToHash("0x111"),
				TxHashSourceRevertStatus: 2,
			},
			expectNil: false,
			validateRevert: func(t *testing.T, revert *domain.RevertDataTransaction) {
				assert.Equal(
					t,
					"0x0000000000000000000000000000000000000000000000000000000000000111",
					revert.TxHashSourceRevert,
				)
				assert.Equal(t, uint8(2), revert.TxHashSourceRevertStatus)
			},
		},
		{
			name: "source revert all fields",
			data: types.AtomicTeleportAdditionalData{
				TxHashSourceRevert:       common.HexToHash("0xbbb"),
				TxHashSourceRevertStatus: 2,
			},
			expectNil: false,
			validateRevert: func(t *testing.T, revert *domain.RevertDataTransaction) {
				assert.Equal(
					t,
					"0x0000000000000000000000000000000000000000000000000000000000000bbb",
					revert.TxHashSourceRevert,
				)
				assert.Equal(t, uint8(2), revert.TxHashSourceRevertStatus)
			},
		},
		{
			name: "both destination and source reverts",
			data: types.AtomicTeleportAdditionalData{
				TxHashDestinationRevert:       common.HexToHash("0x222"),
				TxHashDestinationRevertStatus: 1,
				TxHashSourceRevert:            common.HexToHash("0x333"),
				TxHashSourceRevertStatus:      2,
			},
			expectNil: false,
			validateRevert: func(t *testing.T, revert *domain.RevertDataTransaction) {
				assert.Equal(
					t,
					"0x0000000000000000000000000000000000000000000000000000000000000222",
					revert.TxHashDestinationRevert,
				)
				assert.Equal(t, uint8(1), revert.TxHashDestinationRevertStatus)
				assert.Equal(
					t,
					"0x0000000000000000000000000000000000000000000000000000000000000333",
					revert.TxHashSourceRevert,
				)
				assert.Equal(t, uint8(2), revert.TxHashSourceRevertStatus)
			},
		},
		{
			name: "ignores destination unlock data",
			data: types.AtomicTeleportAdditionalData{
				TxHashDestination:          common.HexToHash("0x444"),
				TxHashDestinationTimestamp: 1672531200,
				BatchHubTxHash:             common.HexToHash("0x555"),
				// Only revert field:
				TxHashDestinationRevert: common.HexToHash("0x777"),
			},
			expectNil: false,
			validateRevert: func(t *testing.T, revert *domain.RevertDataTransaction) {
				// Only revert field should be set
				assert.Equal(
					t,
					"0x0000000000000000000000000000000000000000000000000000000000000777",
					revert.TxHashDestinationRevert,
				)
				assert.Equal(t, uint8(0), revert.TxHashDestinationRevertStatus)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			revert := buildRevertDataIfPresent(tt.data)

			if tt.expectNil {
				assert.Nil(t, revert, "expected nil but got revert data")
			} else {
				require.NotNil(t, revert, "expected revert data but got nil")
				if tt.validateRevert != nil {
					tt.validateRevert(t, revert)
				}
			}
		})
	}
}

func TestApplyTransactionUpdates(t *testing.T) {
	buildTransaction := func() *domain.Transaction {
		return &domain.Transaction{
			MessageId:            "msg-123",
			SharedId:             "shared-456",
			TxHashDestination:    "0xoldDest",
			HubTxHash:            "0xoldCC",
			DestinationTimestamp: time.Unix(1234567890, 0),
		}
	}

	tests := []struct {
		name       string
		data       types.AtomicTeleportAdditionalData
		validateTx func(t *testing.T, tx *domain.Transaction)
	}{
		{
			name: "no updates when data is empty",
			data: types.AtomicTeleportAdditionalData{},
			validateTx: func(t *testing.T, tx *domain.Transaction) {
				assert.Equal(t, "0xoldDest", tx.TxHashDestination)
				assert.Equal(t, "0xoldCC", tx.HubTxHash)
				assert.Equal(t, time.Unix(1234567890, 0), tx.DestinationTimestamp)
			},
		},
		{
			name: "updates TxHashDestination only",
			data: types.AtomicTeleportAdditionalData{
				TxHashDestination: common.HexToHash("0x444"),
			},
			validateTx: func(t *testing.T, tx *domain.Transaction) {
				assert.Equal(
					t,
					"0x0000000000000000000000000000000000000000000000000000000000000444",
					tx.TxHashDestination,
				)
				assert.Equal(t, time.Unix(1234567890, 0), tx.DestinationTimestamp)
				assert.Equal(t, "0xoldCC", tx.HubTxHash)
			},
		},
		{
			name: "updates TxHashDestinationTimestamp only",
			data: types.AtomicTeleportAdditionalData{
				TxHashDestinationTimestamp: 1672531200,
			},
			validateTx: func(t *testing.T, tx *domain.Transaction) {
				assert.Equal(t, "0xoldDest", tx.TxHashDestination)
				assert.Equal(t, time.Unix(1672531200, 0), tx.DestinationTimestamp)
				assert.Equal(t, "0xoldCC", tx.HubTxHash)
			},
		},
		{
			name: "updates all destination fields together",
			data: types.AtomicTeleportAdditionalData{
				TxHashDestination:          common.HexToHash("0x444"),
				TxHashDestinationTimestamp: 1672531200,
			},
			validateTx: func(t *testing.T, tx *domain.Transaction) {
				assert.Equal(
					t,
					"0x0000000000000000000000000000000000000000000000000000000000000444",
					tx.TxHashDestination,
				)
				assert.Equal(t, time.Unix(1672531200, 0), tx.DestinationTimestamp)
				assert.Equal(t, "0xoldCC", tx.HubTxHash)
			},
		},
		{
			name: "updates BatchHubTxHash only",
			data: types.AtomicTeleportAdditionalData{
				BatchHubTxHash: common.HexToHash("0x555"),
			},
			validateTx: func(t *testing.T, tx *domain.Transaction) {
				assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000555", tx.HubTxHash)
				assert.Equal(t, "0xoldDest", tx.TxHashDestination)
				assert.Equal(t, time.Unix(1234567890, 0), tx.DestinationTimestamp)
			},
		},
		{
			name: "updates all fields",
			data: types.AtomicTeleportAdditionalData{
				TxHashDestination:          common.HexToHash("0x666"),
				TxHashDestinationTimestamp: 1672531200,
				BatchHubTxHash:             common.HexToHash("0x888"),
			},
			validateTx: func(t *testing.T, tx *domain.Transaction) {
				assert.Equal(
					t,
					"0x0000000000000000000000000000000000000000000000000000000000000666",
					tx.TxHashDestination,
				)
				assert.Equal(t, time.Unix(1672531200, 0), tx.DestinationTimestamp)
				assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000888", tx.HubTxHash)
			},
		},
		{
			name: "ignores revert data",
			data: types.AtomicTeleportAdditionalData{
				TxHashDestination: common.HexToHash("0x444"),
				// Revert fields that should be ignored:
				TxHashDestinationRevert:       common.HexToHash("0x777"),
				TxHashDestinationRevertStatus: 2,
			},
			validateTx: func(t *testing.T, tx *domain.Transaction) {
				// Only unlock data should be applied
				assert.Equal(
					t,
					"0x0000000000000000000000000000000000000000000000000000000000000444",
					tx.TxHashDestination,
				)
				assert.Equal(t, time.Unix(1234567890, 0), tx.DestinationTimestamp)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := buildTransaction()
			applyTransactionUpdates(tx, tt.data)

			assert.Equal(t, "msg-123", tx.MessageId)
			assert.Equal(t, "shared-456", tx.SharedId)

			if tt.validateTx != nil {
				tt.validateTx(t, tx)
			}
		})
	}
}

// TestTeleportEventHandler_ProcessAtomicMessageStatusChangedBatch tests the processAtomicMessageStatusChangedBatch method
func TestTeleportEventHandler_ProcessAtomicMessageStatusChangedBatch(t *testing.T) {
	// Helper to create transactions with proper IDs
	createTx := func(sharedId string, status types.AtomicTeleportStatus) *domain.Transaction {
		tx := &domain.Transaction{
			SharedId:       sharedId,
			TeleportStatus: func() *uint8 { v := uint8(status); return &v }(),
		}
		tx.ID = uuid.New()
		return tx
	}

	tests := []struct {
		name          string
		log           core.ContractLog
		setupMocks    func(*mocks.MockTransactionRepository, *mocks.MockRevertDataTransactionRepository)
		expectedError string
	}{
		{
			name: "success - status changed to executed",
			log: core.ContractLog{
				ContractName: events.ContractTeleport,
				EventName:    events.AtomicMessageStatusChangedBatch,
				RawEventData: testutil.MustMarshal(t, &TeleportV1.TeleportV1AtomicMessageStatusChangedBatch{
					MsgIds: []string{"msg-123", "msg-456"},
					Status: uint8(types.AtomicTeleportExecuted),
				}),
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, revertRepo *mocks.MockRevertDataTransactionRepository) {
				tx1 := createTx("msg-123", types.AtomicTeleportPending)
				tx2 := createTx("msg-456", types.AtomicTeleportPending)

				// Mock GetTransactionsBySharedIDs (called by BuildTransactionMap)
				txRepo.EXPECT().
					GetTransactionsBySharedIDs(gomock.Any(), []string{"msg-123", "msg-456"}, false).
					Return([]domain.Transaction{*tx1, *tx2}, nil)

				// Mock UpdateTransactionsBulk
				txRepo.EXPECT().
					UpdateTransactionsBulk(gomock.Any(), gomock.Any(), "teleport_status", "updated_at").
					Return(nil)
			},
		},
		{
			name: "success - status changed to rejected",
			log: core.ContractLog{
				ContractName: events.ContractTeleport,
				EventName:    events.AtomicMessageStatusChangedBatch,
				RawEventData: testutil.MustMarshal(t, &TeleportV1.TeleportV1AtomicMessageStatusChangedBatch{
					MsgIds: []string{"msg-123", "msg-456"},
					Status: uint8(types.AtomicTeleportRejected),
				}),
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, revertRepo *mocks.MockRevertDataTransactionRepository) {
				tx1 := createTx("msg-123", types.AtomicTeleportPending)
				tx2 := createTx("msg-456", types.AtomicTeleportPending)

				// Mock GetTransactionsBySharedIDs (called by BuildTransactionMap)
				txRepo.EXPECT().
					GetTransactionsBySharedIDs(gomock.Any(), []string{"msg-123", "msg-456"}, false).
					Return([]domain.Transaction{*tx1, *tx2}, nil)

				// Mock CreateCreditMemo for EACH rejected transaction
				txRepo.EXPECT().
					CreateCreditMemo(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(2)

				// Mock UpdateTransactionsBulk
				txRepo.EXPECT().
					UpdateTransactionsBulk(gomock.Any(), gomock.Any(), "teleport_status", "updated_at").
					Return(nil)
			},
		},
		{
			name: "error - type assertion failure",
			log: core.ContractLog{
				ContractName: events.ContractTeleport,
				EventName:    events.AtomicMessageStatusChangedBatch,
				RawEventData: nil,
			},
			setupMocks:    func(txRepo *mocks.MockTransactionRepository, revertRepo *mocks.MockRevertDataTransactionRepository) {},
			expectedError: "failed to unmarshal event data",
		},
		{
			name: "error - BuildTransactionMap fails",
			log: core.ContractLog{
				ContractName: events.ContractTeleport,
				EventName:    events.AtomicMessageStatusChangedBatch,
				RawEventData: testutil.MustMarshal(t, &TeleportV1.TeleportV1AtomicMessageStatusChangedBatch{
					MsgIds: []string{"msg-123"},
					Status: uint8(types.AtomicTeleportExecuted),
				}),
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, revertRepo *mocks.MockRevertDataTransactionRepository) {
				txRepo.EXPECT().
					GetTransactionsBySharedIDs(gomock.Any(), []string{"msg-123"}, false).
					Return(nil, fmt.Errorf("db error"))
			},
			expectedError: "failed to retrieve transactions from database",
		},
		{
			name: "error - UpdateTransactionsBulk fails",
			log: core.ContractLog{
				ContractName: events.ContractTeleport,
				EventName:    events.AtomicMessageStatusChangedBatch,
				RawEventData: testutil.MustMarshal(t, &TeleportV1.TeleportV1AtomicMessageStatusChangedBatch{
					MsgIds: []string{"msg-123"},
					Status: uint8(types.AtomicTeleportExecuted),
				}),
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, revertRepo *mocks.MockRevertDataTransactionRepository) {
				tx1 := createTx("msg-123", types.AtomicTeleportPending)
				txRepo.EXPECT().
					GetTransactionsBySharedIDs(gomock.Any(), []string{"msg-123"}, false).
					Return([]domain.Transaction{*tx1}, nil)
				txRepo.EXPECT().
					UpdateTransactionsBulk(gomock.Any(), gomock.Any(), "teleport_status", "updated_at").
					Return(fmt.Errorf("update failed"))
			},
			expectedError: "update failed",
		},
		{
			name: "error - CreateCreditMemo fails",
			log: core.ContractLog{
				ContractName: events.ContractTeleport,
				EventName:    events.AtomicMessageStatusChangedBatch,
				RawEventData: testutil.MustMarshal(t, &TeleportV1.TeleportV1AtomicMessageStatusChangedBatch{
					MsgIds: []string{"msg-123"},
					Status: uint8(types.AtomicTeleportRejected),
				}),
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, revertRepo *mocks.MockRevertDataTransactionRepository) {
				tx1 := createTx("msg-123", types.AtomicTeleportPending)
				txRepo.EXPECT().
					GetTransactionsBySharedIDs(gomock.Any(), []string{"msg-123"}, false).
					Return([]domain.Transaction{*tx1}, nil)
				txRepo.EXPECT().
					CreateCreditMemo(gomock.Any(), gomock.Any()).
					Return(fmt.Errorf("credit memo creation failed"))
			},
			expectedError: "failed to create credit memo",
		},
		{
			name: "skips transactions not in pending status",
			log: core.ContractLog{
				ContractName: events.ContractTeleport,
				EventName:    events.AtomicMessageStatusChangedBatch,
				RawEventData: testutil.MustMarshal(t, &TeleportV1.TeleportV1AtomicMessageStatusChangedBatch{
					MsgIds: []string{"msg-123", "msg-456", "msg-789"},
					Status: uint8(types.AtomicTeleportExecuted),
				}),
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, revertRepo *mocks.MockRevertDataTransactionRepository) {
				tx1 := createTx("msg-123", types.AtomicTeleportPending)  // Will be updated
				tx2 := createTx("msg-456", types.AtomicTeleportExecuted) // Already executed - skipped
				tx3 := createTx("msg-789", types.AtomicTeleportRejected) // Already rejected - skipped

				txRepo.EXPECT().
					GetTransactionsBySharedIDs(gomock.Any(), []string{"msg-123", "msg-456", "msg-789"}, false).
					Return([]domain.Transaction{*tx1, *tx2, *tx3}, nil)

				// Only msg-123 should be updated (tx1 is pending)
				txRepo.EXPECT().
					UpdateTransactionsBulk(gomock.Any(), gomock.Any(), "teleport_status", "updated_at").
					Return(nil)
			},
		},
		{
			name: "skips transactions not found in database",
			log: core.ContractLog{
				ContractName: events.ContractTeleport,
				EventName:    events.AtomicMessageStatusChangedBatch,
				RawEventData: testutil.MustMarshal(t, &TeleportV1.TeleportV1AtomicMessageStatusChangedBatch{
					MsgIds: []string{"msg-123", "msg-nonexistent"},
					Status: uint8(types.AtomicTeleportExecuted),
				}),
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, revertRepo *mocks.MockRevertDataTransactionRepository) {
				tx1 := createTx("msg-123", types.AtomicTeleportPending)

				txRepo.EXPECT().
					GetTransactionsBySharedIDs(gomock.Any(), []string{"msg-123", "msg-nonexistent"}, false).
					Return([]domain.Transaction{*tx1}, nil)

				txRepo.EXPECT().
					UpdateTransactionsBulk(gomock.Any(), gomock.Any(), "teleport_status", "updated_at").
					Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTxRepo := mocks.NewMockTransactionRepository(ctrl)
			mockRevertRepo := mocks.NewMockRevertDataTransactionRepository(ctrl)
			mockDecryptor := mocks.NewMockDecryptor(ctrl)
			stubLogger := &testutil.StubLogger{}
			pnData := &core.PNodeDataAndSecrets{}

			if tt.setupMocks != nil {
				tt.setupMocks(mockTxRepo, mockRevertRepo)
			}

			handler := NewTeleportEventHandler(mockTxRepo, mockRevertRepo, mockDecryptor, stubLogger, pnData)
			err := handler.Handle(context.Background(), tt.log)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestTeleportEventHandler_ProcessAtomicMessageAdditionalDataBatch tests the processAtomicMessageAdditionalDataBatch method
func TestTeleportEventHandler_ProcessAtomicMessageAdditionalDataBatch(t *testing.T) {
	// Helper to create test data and serialize it to JSON
	createTestData := func() []byte {
		data := []types.AtomicTeleportAdditionalData{
			{
				SharedId:                   "shared-123",
				TxHashDestination:          common.HexToHash("0x456"),
				TxHashDestinationTimestamp: 1672531200,
				BatchHubTxHash:             common.HexToHash("0x789"),
			},
			{
				SharedId:                      "shared-456",
				TxHashDestinationRevert:       common.HexToHash("0xabc"),
				TxHashDestinationRevertStatus: 2,
			},
		}
		jsonBytes, _ := json.Marshal(data)
		return jsonBytes
	}

	// Helper to create a transaction
	createTransaction := func(sharedId string) *domain.Transaction {
		tx := &domain.Transaction{
			SharedId: sharedId,
		}
		tx.ID = uuid.New()
		return tx
	}

	tests := []struct {
		name          string
		log           core.ContractLog
		setupMocks    func(*mocks.MockTransactionRepository, *mocks.MockRevertDataTransactionRepository, *mocks.MockDecryptor)
		expectedError bool
		errorContains string
	}{
		{
			name: "type assertion failure returns error",
			log: core.ContractLog{
				ContractName: events.ContractTeleport,
				EventName:    events.AtomicMessageAdditionalDataBatch,
				RawEventData: nil,
				BlockNumber:  100,
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, revertRepo *mocks.MockRevertDataTransactionRepository, decryptor *mocks.MockDecryptor) {
			},
			expectedError: true,
			errorContains: "failed to unmarshal event data",
		},
		{
			name: "empty encrypted data returns nil",
			log: core.ContractLog{
				ContractName: events.ContractTeleport,
				EventName:    events.AtomicMessageAdditionalDataBatch,
				RawEventData: testutil.MustMarshal(t, &TeleportV1.TeleportV1AtomicMessageAdditionalDataBatch{
					EncryptedData: "",
				}),
				BlockNumber: 100,
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, revertRepo *mocks.MockRevertDataTransactionRepository, decryptor *mocks.MockDecryptor) {
			},
			expectedError: false,
		},
		{
			name: "hex decode failure returns error",
			log: core.ContractLog{
				ContractName: events.ContractTeleport,
				EventName:    events.AtomicMessageAdditionalDataBatch,
				RawEventData: testutil.MustMarshal(t, &TeleportV1.TeleportV1AtomicMessageAdditionalDataBatch{
					EncryptedData: "not-valid-hex",
				}),
				BlockNumber: 100,
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, revertRepo *mocks.MockRevertDataTransactionRepository, decryptor *mocks.MockDecryptor) {
			},
			expectedError: true,
			errorContains: "failed to decode hex",
		},
		{
			name: "decryption failure returns error",
			log: core.ContractLog{
				ContractName: events.ContractTeleport,
				EventName:    events.AtomicMessageAdditionalDataBatch,
				RawEventData: testutil.MustMarshal(t, &TeleportV1.TeleportV1AtomicMessageAdditionalDataBatch{
					EncryptedData: hex.EncodeToString([]byte("encrypted")),
				}),
				BlockNumber: 100,
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, revertRepo *mocks.MockRevertDataTransactionRepository, decryptor *mocks.MockDecryptor) {
				decryptor.EXPECT().
					DecryptPayloadBytes(gomock.Any(), uint64(100), gomock.Any(), types.AtomicSecret).
					Return(nil, errors.New("decryption failed"))
			},
			expectedError: true,
			errorContains: "decryption failed",
		},
		{
			name: "BuildTransactionMap failure returns error",
			log: core.ContractLog{
				ContractName: events.ContractTeleport,
				EventName:    events.AtomicMessageAdditionalDataBatch,
				RawEventData: testutil.MustMarshal(t, &TeleportV1.TeleportV1AtomicMessageAdditionalDataBatch{
					EncryptedData: hex.EncodeToString([]byte("encrypted")),
					MsgIds:        []string{"msg-1", "msg-2"},
				}),
				BlockNumber: 100,
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, revertRepo *mocks.MockRevertDataTransactionRepository, decryptor *mocks.MockDecryptor) {
				decryptor.EXPECT().
					DecryptPayloadBytes(gomock.Any(), uint64(100), gomock.Any(), types.AtomicSecret).
					Return(createTestData(), nil)

				txRepo.EXPECT().
					GetTransactionsBySharedIDs(gomock.Any(), gomock.Any(), false).
					Return(nil, errors.New("some error"))
			},
			expectedError: true,
			errorContains: "failed to fetch transactions",
		},
		{
			name: "transaction not found in map - warn and skip (non-blocking)",
			log: core.ContractLog{
				ContractName: events.ContractTeleport,
				EventName:    events.AtomicMessageAdditionalDataBatch,
				RawEventData: testutil.MustMarshal(t, &TeleportV1.TeleportV1AtomicMessageAdditionalDataBatch{
					EncryptedData: hex.EncodeToString([]byte("encrypted")),
					MsgIds:        []string{"msg-1"},
				}),
				BlockNumber: 100,
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, revertRepo *mocks.MockRevertDataTransactionRepository, decryptor *mocks.MockDecryptor) {
				data := []types.AtomicTeleportAdditionalData{
					{
						SharedId:          "shared-123",
						TxHashDestination: common.HexToHash("0x456"),
					},
				}
				jsonBytes, _ := json.Marshal(data)

				decryptor.EXPECT().
					DecryptPayloadBytes(gomock.Any(), uint64(100), gomock.Any(), types.AtomicSecret).
					Return(jsonBytes, nil)

				// Return empty slice - no transactions found
				txRepo.EXPECT().
					GetTransactionsBySharedIDs(gomock.Any(), []string{"shared-123"}, false).
					Return([]domain.Transaction{}, nil)

				txRepo.EXPECT().
					UpdateTransactionsBulk(gomock.Any(), gomock.Any(), "tx_hash_destination", "destination_timestamp", "hub_tx_hash", "updated_at").
					DoAndReturn(func(ctx context.Context, txs []*domain.Transaction, _ ...string) error {
						require.Len(t, txs, 0, "expected empty transactions slice")
						return nil
					})
			},
			expectedError: false,
		},
		{
			name: "UpdateTransactionsBulk failure returns error",
			log: core.ContractLog{
				ContractName: events.ContractTeleport,
				EventName:    events.AtomicMessageAdditionalDataBatch,
				RawEventData: testutil.MustMarshal(t, &TeleportV1.TeleportV1AtomicMessageAdditionalDataBatch{
					EncryptedData: hex.EncodeToString([]byte("encrypted")),
					MsgIds:        []string{"msg-1"},
				}),
				BlockNumber: 100,
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, revertRepo *mocks.MockRevertDataTransactionRepository, decryptor *mocks.MockDecryptor) {
				data := []types.AtomicTeleportAdditionalData{
					{
						SharedId:          "shared-123",
						TxHashDestination: common.HexToHash("0x456"),
					},
				}
				jsonBytes, _ := json.Marshal(data)

				decryptor.EXPECT().
					DecryptPayloadBytes(gomock.Any(), uint64(100), gomock.Any(), types.AtomicSecret).
					Return(jsonBytes, nil)

				tx := createTransaction("shared-123")
				txRepo.EXPECT().
					GetTransactionsBySharedIDs(gomock.Any(), []string{"shared-123"}, false).
					Return([]domain.Transaction{*tx}, nil)

				txRepo.EXPECT().
					UpdateTransactionsBulk(gomock.Any(), gomock.Any(), "tx_hash_destination", "destination_timestamp", "hub_tx_hash", "updated_at").
					Return(errors.New("update failed"))
			},
			expectedError: true,
			errorContains: "update failed",
		},
		{
			name: "CreateRevertTransactions failure returns error",
			log: core.ContractLog{
				ContractName: events.ContractTeleport,
				EventName:    events.AtomicMessageAdditionalDataBatch,
				RawEventData: testutil.MustMarshal(t, &TeleportV1.TeleportV1AtomicMessageAdditionalDataBatch{
					EncryptedData: hex.EncodeToString([]byte("encrypted")),
					MsgIds:        []string{"msg-1"},
				}),
				BlockNumber: 100,
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, revertRepo *mocks.MockRevertDataTransactionRepository, decryptor *mocks.MockDecryptor) {
				// Data with revert information
				data := []types.AtomicTeleportAdditionalData{
					{
						SharedId:                "shared-123",
						TxHashDestinationRevert: common.HexToHash("0xabc"),
					},
				}
				jsonBytes, _ := json.Marshal(data)

				decryptor.EXPECT().
					DecryptPayloadBytes(gomock.Any(), uint64(100), gomock.Any(), types.AtomicSecret).
					Return(jsonBytes, nil)

				tx := createTransaction("shared-123")
				txRepo.EXPECT().
					GetTransactionsBySharedIDs(gomock.Any(), []string{"shared-123"}, false).
					Return([]domain.Transaction{*tx}, nil)

				txRepo.EXPECT().
					UpdateTransactionsBulk(gomock.Any(), gomock.Any(), "tx_hash_destination", "destination_timestamp", "hub_tx_hash", "updated_at").
					Return(nil)

				revertRepo.EXPECT().
					CreateRevertTransactions(gomock.Any(), gomock.Any()).
					Return(errors.New("create revert failed"))
			},
			expectedError: true,
			errorContains: "create revert failed",
		},
		{
			name: "success - no revert data, only transactions updated",
			log: core.ContractLog{
				ContractName: events.ContractTeleport,
				EventName:    events.AtomicMessageAdditionalDataBatch,
				RawEventData: testutil.MustMarshal(t, &TeleportV1.TeleportV1AtomicMessageAdditionalDataBatch{
					EncryptedData: hex.EncodeToString([]byte("encrypted")),
					MsgIds:        []string{"msg-1"},
				}),
				BlockNumber: 100,
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, revertRepo *mocks.MockRevertDataTransactionRepository, decryptor *mocks.MockDecryptor) {
				// Data without revert information
				data := []types.AtomicTeleportAdditionalData{
					{
						SharedId:                   "shared-123",
						TxHashDestination:          common.HexToHash("0x456"),
						TxHashDestinationTimestamp: 1672531200,
					},
				}
				jsonBytes, _ := json.Marshal(data)

				decryptor.EXPECT().
					DecryptPayloadBytes(gomock.Any(), uint64(100), gomock.Any(), types.AtomicSecret).
					Return(jsonBytes, nil)

				tx := createTransaction("shared-123")
				txRepo.EXPECT().
					GetTransactionsBySharedIDs(gomock.Any(), []string{"shared-123"}, false).
					Return([]domain.Transaction{*tx}, nil)

				txRepo.EXPECT().
					UpdateTransactionsBulk(gomock.Any(), gomock.Any(), "tx_hash_destination", "destination_timestamp", "hub_tx_hash", "updated_at").
					DoAndReturn(func(ctx context.Context, txs []*domain.Transaction, _ ...string) error {
						require.Len(t, txs, 1)
						assert.Equal(t, "shared-123", txs[0].SharedId)
						assert.Equal(
							t,
							"0x0000000000000000000000000000000000000000000000000000000000000456",
							txs[0].TxHashDestination,
						)
						return nil
					})
			},
			expectedError: false,
		},
		{
			name: "success - with revert data, transactions and reverts created",
			log: core.ContractLog{
				ContractName: events.ContractTeleport,
				EventName:    events.AtomicMessageAdditionalDataBatch,
				RawEventData: testutil.MustMarshal(t, &TeleportV1.TeleportV1AtomicMessageAdditionalDataBatch{
					EncryptedData: hex.EncodeToString([]byte("encrypted")),
					MsgIds:        []string{"msg-1", "msg-2"},
				}),
				BlockNumber: 100,
			},
			setupMocks: func(txRepo *mocks.MockTransactionRepository, revertRepo *mocks.MockRevertDataTransactionRepository, decryptor *mocks.MockDecryptor) {
				decryptor.EXPECT().
					DecryptPayloadBytes(gomock.Any(), uint64(100), gomock.Any(), types.AtomicSecret).
					Return(createTestData(), nil)

				tx1 := createTransaction("shared-123")
				tx2 := createTransaction("shared-456")
				txRepo.EXPECT().
					GetTransactionsBySharedIDs(gomock.Any(), []string{"shared-123", "shared-456"}, false).
					Return([]domain.Transaction{*tx1, *tx2}, nil)

				txRepo.EXPECT().
					UpdateTransactionsBulk(gomock.Any(), gomock.Any(), "tx_hash_destination", "destination_timestamp", "hub_tx_hash", "updated_at").
					DoAndReturn(func(ctx context.Context, txs []*domain.Transaction, _ ...string) error {
						require.Len(t, txs, 2)
						// First transaction has destination data
						assert.Equal(t, "shared-123", txs[0].SharedId)
						assert.Equal(
							t,
							"0x0000000000000000000000000000000000000000000000000000000000000456",
							txs[0].TxHashDestination,
						)
						// Second transaction has revert data
						assert.Equal(t, "shared-456", txs[1].SharedId)
						return nil
					})

				revertRepo.EXPECT().
					CreateRevertTransactions(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, reverts []*domain.RevertDataTransaction) error {
						require.Len(t, reverts, 1)
						assert.Equal(t, tx2.ID, reverts[0].TransactionId)
						assert.Equal(
							t,
							"0x0000000000000000000000000000000000000000000000000000000000000abc",
							reverts[0].TxHashDestinationRevert,
						)
						return nil
					})
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTxRepo := mocks.NewMockTransactionRepository(ctrl)
			mockRevertRepo := mocks.NewMockRevertDataTransactionRepository(ctrl)
			mockDecryptor := mocks.NewMockDecryptor(ctrl)
			logger := &testutil.StubLogger{}
			pnData := &core.PNodeDataAndSecrets{}

			handler := NewTeleportEventHandler(mockTxRepo, mockRevertRepo, mockDecryptor, logger, pnData)

			tt.setupMocks(mockTxRepo, mockRevertRepo, mockDecryptor)

			err := handler.processAtomicMessageAdditionalDataBatch(context.Background(), tt.log)

			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestTeleportEventHandler_PersistTransferMessages tests the persistTransferMessages method
func TestTeleportEventHandler_PersistTransferMessages(t *testing.T) {
	// buildMessage is a test helper to create a DispatchedMessageToPrivateHub
	buildMessage := func(sharedID string, amount int64, txType types.BridgeTransactionType) types.DispatchedMessageToPrivateHub {
		return types.DispatchedMessageToPrivateHub{
			MessageId:       [32]byte{1, 2, 3},
			From:            common.HexToAddress("0x1234567890123456789012345678901234567890"),
			To:              common.HexToAddress("0x0987654321098765432109876543210987654321"),
			ToChainId:       big.NewInt(1),
			FromChainId:     big.NewInt(2),
			SharedId:        sharedID,
			TransactionType: txType,
			Data: types.RaylzMessage{
				MessageMetadata: types.RaylzMessageMetadata{
					ResourceId: [32]byte{10, 20, 30},
					TransferMetadata: types.BridgedTransferMetadata{
						Amount:    big.NewInt(amount),
						AssetType: types.AssetTypeERC20,
						From:      "0x1234567890123456789012345678901234567890",
						To:        "0x0987654321098765432109876543210987654321",
					},
				},
			},
		}
	}

	tests := []struct {
		name          string
		batch         []types.DispatchedMessageToPrivateHub
		blockNumber   uint64
		setupMocks    func(*mocks.MockTransactionRepository)
		expectedError string
	}{
		{
			name: "success - converts all transfer messages and persists them",
			batch: []types.DispatchedMessageToPrivateHub{
				buildMessage("shared-123", 1000, types.Transfer),
				buildMessage("shared-456", 2000, types.Transfer),
			},
			blockNumber: 100,
			setupMocks: func(txRepo *mocks.MockTransactionRepository) {
				txRepo.EXPECT().
					CreateTransactionsWithPromotion(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedError: "",
		},
		{
			name: "success - filters out non-transfer messages",
			batch: []types.DispatchedMessageToPrivateHub{
				buildMessage("shared-123", 1000, types.Transfer),
				buildMessage("shared-456", 2000, types.Proof), // Should be filtered out
				buildMessage("shared-789", 5000, types.Transfer),
			},
			blockNumber: 100,
			setupMocks: func(txRepo *mocks.MockTransactionRepository) {
				txRepo.EXPECT().
					CreateTransactionsWithPromotion(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedError: "",
		},
		{
			name: "error - CreateTransactionsWithPromotion fails",
			batch: []types.DispatchedMessageToPrivateHub{
				buildMessage("shared-123", 1000, types.Transfer),
			},
			blockNumber: 100,
			setupMocks: func(txRepo *mocks.MockTransactionRepository) {
				txRepo.EXPECT().
					CreateTransactionsWithPromotion(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("database connection lost"))
			},
			expectedError: "database connection lost",
		},
		{
			name:        "success - empty batch",
			batch:       []types.DispatchedMessageToPrivateHub{},
			blockNumber: 100,
			setupMocks: func(txRepo *mocks.MockTransactionRepository) {
				// CreateTransactionsWithPromotion should NOT be called for empty batch
			},
			expectedError: "",
		},
		{
			name: "success - all messages are non-transfer (filtered out)",
			batch: []types.DispatchedMessageToPrivateHub{
				buildMessage("shared-123", 1000, types.Proof),
				buildMessage("shared-456", 2000, types.Proof),
			},
			blockNumber: 100,
			setupMocks: func(txRepo *mocks.MockTransactionRepository) {
				// CreateTransactionsWithPromotion should NOT be called when all filtered out
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTxRepo := mocks.NewMockTransactionRepository(ctrl)
			mockRevertRepo := mocks.NewMockRevertDataTransactionRepository(ctrl)
			mockDecryptor := mocks.NewMockDecryptor(ctrl)
			stubLogger := &testutil.StubLogger{}
			pnData := &core.PNodeDataAndSecrets{}

			tt.setupMocks(mockTxRepo)

			handler := NewTeleportEventHandler(mockTxRepo, mockRevertRepo, mockDecryptor, stubLogger, pnData)
			err := handler.persistTransferMessages(
				context.Background(),
				tt.batch,
				tt.blockNumber,
				common.HexToHash("0xcc123").String(),
				time.Time{},
			)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTeleportEventHandler_ProcessEncryptedDataBatchStored(t *testing.T) {
	tests := []struct {
		name          string
		log           core.ContractLog
		setupMocks    func(*mocks.MockDecryptor, *mocks.MockTransactionRepository)
		expectedError string
	}{
		{
			name: "error - type assertion failure",
			log: core.ContractLog{
				ContractName: events.ContractTeleport,
				EventName:    events.EncryptedDataBatchStored,
				RawEventData: nil,
				BlockNumber:  100,
			},
			setupMocks:    func(mockDecryptor *mocks.MockDecryptor, mockTxRepo *mocks.MockTransactionRepository) {},
			expectedError: "failed to unmarshal event data",
		},
		{
			name: "error - decryption failure",
			log: core.ContractLog{
				ContractName: events.ContractTeleport,
				EventName:    events.EncryptedDataBatchStored,
				RawEventData: testutil.MustMarshal(t, &TeleportV1.TeleportV1EncryptedDataBatchStored{
					Data: []byte("encrypted-data"),
				}),
				BlockNumber: 100,
			},
			setupMocks: func(mockDecryptor *mocks.MockDecryptor, mockTxRepo *mocks.MockTransactionRepository) {
				// Mock decryption to fail
				mockDecryptor.EXPECT().
					DecryptPayloadBytes([]byte("encrypted-data"), uint64(100), gomock.Any(), types.ParticipantSecret).
					Return(nil, errors.New("decryption failed"))
			},
			expectedError: "decryption failed",
		},
		{
			name: "error - persistTransferMessages fails",
			log: core.ContractLog{
				ContractName: events.ContractTeleport,
				EventName:    events.EncryptedDataBatchStored,
				RawEventData: testutil.MustMarshal(t, &TeleportV1.TeleportV1EncryptedDataBatchStored{
					Data: []byte("encrypted-data"),
				}),
				BlockNumber: 100,
			},
			setupMocks: func(mockDecryptor *mocks.MockDecryptor, mockTxRepo *mocks.MockTransactionRepository) {
				// Mock decryption to return a transfer message
				batch := []types.DispatchedMessageToPrivateHub{
					{
						SharedId:        "transfer-1",
						TransactionType: types.Transfer,
						MessageId:       [32]byte{1, 2, 3},
						From:            common.HexToAddress("0x1234567890123456789012345678901234567890"),
						To:              common.HexToAddress("0x0987654321098765432109876543210987654321"),
						ToChainId:       big.NewInt(1),
						FromChainId:     big.NewInt(2),
						Data: types.RaylzMessage{
							MessageMetadata: types.RaylzMessageMetadata{
								ResourceId: [32]byte{10, 20, 30},
								TransferMetadata: types.BridgedTransferMetadata{
									Amount:    big.NewInt(100),
									AssetType: types.AssetTypeERC20,
									From:      "0x1234567890123456789012345678901234567890",
									To:        "0x0987654321098765432109876543210987654321",
								},
							},
						},
					},
				}
				batchJSON, _ := json.Marshal(batch)
				mockDecryptor.EXPECT().
					DecryptPayloadBytes([]byte("encrypted-data"), uint64(100), gomock.Any(), types.ParticipantSecret).
					Return(batchJSON, nil)

				mockTxRepo.EXPECT().
					CreateTransactionsWithPromotion(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("database error"))
			},
			expectedError: "database error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTxRepo := mocks.NewMockTransactionRepository(ctrl)
			mockRevertRepo := mocks.NewMockRevertDataTransactionRepository(ctrl)
			mockDecryptor := mocks.NewMockDecryptor(ctrl)
			logger := &testutil.StubLogger{}
			pnData := &core.PNodeDataAndSecrets{}

			tc.setupMocks(mockDecryptor, mockTxRepo)

			handler := NewTeleportEventHandler(mockTxRepo, mockRevertRepo, mockDecryptor, logger, pnData)

			err := handler.processEncryptedDataBatchStored(context.Background(), tc.log)

			if tc.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTeleportEventHandler_Name(t *testing.T) {
	handler := NewTeleportEventHandler(nil, nil, nil, nil, nil)

	name := handler.Name()

	assert.Equal(t, "TeleportHandler", name)
}

func TestTeleportEventHandler_classifyProtocol(t *testing.T) {
	tests := []struct {
		name     string
		msg      types.DispatchedMessageToPrivateHub
		expected types.ProtocolType
	}{
		{
			name: "atomic protocol",
			msg: types.DispatchedMessageToPrivateHub{
				IsAtomic: true,
			},
			expected: types.Atomic,
		},
		{
			name: "vanilla protocol",
			msg: types.DispatchedMessageToPrivateHub{
				IsAtomic:        false,
				TransactionType: types.Proof,
				Data: types.RaylzMessage{
					MessageMetadata: types.RaylzMessageMetadata{
						TransferMetadata: types.BridgedTransferMetadata{
							AssetType: types.AssetTypeCustom,
						},
					},
				},
			},
			expected: types.Vanilla,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewTeleportEventHandler(nil, nil, nil, &testutil.StubLogger{}, nil)

			protocol, _ := handler.classifyProtocol(&tt.msg)

			assert.Equal(t, tt.expected, protocol, "expected protocol %v but got %v", tt.expected, protocol)
		})
	}
}
