package handlers

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/core"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/mocks"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/testutil"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/contracts/DvpTeleport"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/types"
)

func TestNewDvpTeleportEventHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTxRepo := mocks.NewMockTransactionRepository(ctrl)
	mockDecryptor := mocks.NewMockDecryptor(ctrl)
	mockSwapSalts := mocks.NewMockSwapSaltsStore(ctrl)
	mockLogger := &testutil.StubLogger{}
	pnData := &core.PNodeDataAndSecrets{}

	handler := NewDvpTeleportEventHandler(mockTxRepo, mockDecryptor, nil, mockLogger, pnData, mockSwapSalts)

	assert.NotNil(t, handler)
	assert.Equal(t, mockTxRepo, handler.txRepo)
	assert.Equal(t, mockDecryptor, handler.decryptor)
	assert.Equal(t, mockLogger, handler.log)
	assert.Equal(t, pnData, handler.pnData)
	assert.Equal(t, mockSwapSalts, handler.swapSalts)
}

func TestDvpTeleportEventHandler_Name(t *testing.T) {
	handler := NewDvpTeleportEventHandler(nil, nil, nil, nil, nil, nil)

	assert.Equal(t, "DvpTeleportHandler", handler.Name())
}

func TestDvpTeleportEventHandler_Handle_UnknownEventReturnsNoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler := NewDvpTeleportEventHandler(
		mocks.NewMockTransactionRepository(ctrl),
		mocks.NewMockDecryptor(ctrl),
		nil,
		&testutil.StubLogger{},
		&core.PNodeDataAndSecrets{},
		mocks.NewMockSwapSaltsStore(ctrl),
	)

	err := handler.Handle(context.Background(), core.ContractLog{EventName: "UnknownEvent"})

	assert.NoError(t, err)
}

func TestDvpTeleportEventHandler_ProcessSwapInitiatedEvent_PersistsTokenInAndTokenOutTxs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockTxRepo := mocks.NewMockTransactionRepository(ctrl)
	mockDecryptor := mocks.NewMockDecryptor(ctrl)
	mockSwapSalts := mocks.NewMockSwapSaltsStore(ctrl)

	sharedIdBytes := [32]byte{}
	copy(sharedIdBytes[:], "shared-topic")
	sharedIdKey := hex.EncodeToString(sharedIdBytes[:])

	initiatorSelfSalt := big.NewInt(0xC0FFEE)
	initiatorCtxtSalt := []byte("ctxt-derived-salt")
	decryptedData := types.DvpSwapMessage{
		SharedId:           "shared-123",
		To:                 "0xRecipient",
		ChainId:            big.NewInt(1),
		PNTxHash:           "0xpnhash",
		PNTxTimestamp:      time.Now(),
		TokenInType:        types.DvpEnygma,
		TokenInID:          "1",
		TokenInResourceID:  "0xenygma-resource",
		TokenInAmount:      big.NewInt(500),
		TokenOutType:       types.DvpERC721,
		TokenOutID:         "123",
		TokenOutResourceID: "0xresource",
		TokenOutAmount:     big.NewInt(1000),
		InitiatorSelfSalt:  initiatorSelfSalt,
	}

	mockDecryptor.EXPECT().
		DecryptSwapPayload([]byte("ctxt"), []byte("encrypted"), uint64(1000), gomock.Any()).
		DoAndReturn(func(_, _ []byte, _ uint64, _ core.PNodeDataAndSecrets) ([]byte, []byte, error) {
			plaintext, err := json.Marshal(decryptedData)
			return plaintext, initiatorCtxtSalt, err
		})
	mockSwapSalts.EXPECT().
		Put(gomock.Any(), sharedIdKey, types.DvpSwapSalts{
			InitiatorSelfSalt: initiatorSelfSalt.Bytes(),
			InitiatorCtxtSalt: initiatorCtxtSalt,
		}).
		Return(nil)

	mockTxRepo.EXPECT().
		CreateTransactions(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, txs []*domain.Transaction) error {
			assert.Len(t, txs, 2)
			assert.Equal(t, types.CrossChain, txs[0].TxType)
			assert.Equal(t, types.CrossChain, txs[1].TxType)
			assert.Equal(t, "shared-123", txs[0].SharedId)
			assert.Equal(t, "shared-123", txs[1].SharedId)
			assert.Equal(t, types.DvpSwap, txs[0].Protocol)
			assert.NotNil(t, txs[0].TeleportStatus)
			assert.Equal(t, uint8(types.DvpSwapStatePending), *txs[0].TeleportStatus)
			return nil
		})

	handler := &DvpTeleportEventHandler{
		txRepo:    mockTxRepo,
		decryptor: mockDecryptor,
		log:       &testutil.StubLogger{},
		pnData:    &core.PNodeDataAndSecrets{},
		swapSalts: mockSwapSalts,
	}

	log := core.ContractLog{
		RawEventData: testutil.MustMarshal(t, &DvpTeleport.DvpTeleportSwapInitiated{
			SharedId:      sharedIdBytes,
			EncryptedData: []byte("encrypted"),
			Ctxt:          []byte("ctxt"),
		}),
		BlockNumber:     1000,
		TransactionHash: "0xabc",
		LogIndex:        1,
	}

	err := handler.processSwapInitiatedEvent(ctx, log)

	assert.NoError(t, err)
}

func TestDvpTeleportEventHandler_ProcessSwapInitiatedEvent_ReturnsErrorOnUnmarshal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler := &DvpTeleportEventHandler{
		txRepo:    mocks.NewMockTransactionRepository(ctrl),
		decryptor: mocks.NewMockDecryptor(ctrl),
		log:       &testutil.StubLogger{},
		pnData:    &core.PNodeDataAndSecrets{},
		swapSalts: mocks.NewMockSwapSaltsStore(ctrl),
	}

	err := handler.processSwapInitiatedEvent(context.Background(), core.ContractLog{RawEventData: nil})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal event data")
}

func TestDvpTeleportEventHandler_ProcessSwapInitiatedEvent_ReturnsErrorOnDecrypt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDecryptor := mocks.NewMockDecryptor(ctrl)
	mockDecryptor.EXPECT().
		DecryptSwapPayload(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil, assert.AnError)

	handler := &DvpTeleportEventHandler{
		txRepo:    mocks.NewMockTransactionRepository(ctrl),
		decryptor: mockDecryptor,
		log:       &testutil.StubLogger{},
		pnData:    &core.PNodeDataAndSecrets{},
		swapSalts: mocks.NewMockSwapSaltsStore(ctrl),
	}

	log := core.ContractLog{
		RawEventData: testutil.MustMarshal(t, &DvpTeleport.DvpTeleportSwapInitiated{
			EncryptedData: []byte("encrypted"),
			Ctxt:          []byte("ctxt"),
		}),
		BlockNumber: 1000,
	}

	err := handler.processSwapInitiatedEvent(context.Background(), log)

	assert.Error(t, err)
}

func TestDvpTeleportEventHandler_ProcessSwapInitiatedEvent_ReturnsErrorWhenInitiatorSelfSaltMissing(t *testing.T) {
	// Without InitiatorSelfSalt we can't decrypt the matching SwapCompleted,
	// so refuse the event rather than persist unusable salts.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDecryptor := mocks.NewMockDecryptor(ctrl)
	mockDecryptor.EXPECT().
		DecryptSwapPayload(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(_, _ []byte, _ uint64, _ core.PNodeDataAndSecrets) ([]byte, []byte, error) {
			plaintext, err := json.Marshal(types.DvpSwapMessage{SharedId: "shared-nosalt"})
			return plaintext, []byte("ctxt-salt"), err
		})

	handler := &DvpTeleportEventHandler{
		txRepo:    mocks.NewMockTransactionRepository(ctrl),
		decryptor: mockDecryptor,
		log:       &testutil.StubLogger{},
		pnData:    &core.PNodeDataAndSecrets{},
		swapSalts: mocks.NewMockSwapSaltsStore(ctrl),
	}

	log := core.ContractLog{
		RawEventData: testutil.MustMarshal(t, &DvpTeleport.DvpTeleportSwapInitiated{
			EncryptedData: []byte("encrypted"),
			Ctxt:          []byte("ctxt"),
		}),
		BlockNumber: 1000,
	}

	err := handler.processSwapInitiatedEvent(context.Background(), log)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InitiatorSelfSalt")
}

func TestDvpTeleportEventHandler_ProcessSwapCompletedEvent_FillsConfirmationAndMarksCompleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTxRepo := mocks.NewMockTransactionRepository(ctrl)
	mockDecryptor := mocks.NewMockDecryptor(ctrl)
	mockSwapSalts := mocks.NewMockSwapSaltsStore(ctrl)

	sharedIdBytes := [32]byte{}
	copy(sharedIdBytes[:], "shared-complete")
	sharedIdKey := hex.EncodeToString(sharedIdBytes[:])
	storedSalts := types.DvpSwapSalts{
		InitiatorSelfSalt: []byte("stored-self-salt"),
		InitiatorCtxtSalt: []byte("stored-ctxt-salt"),
	}

	decryptedData := types.DvpSwapMessage{
		SharedId:      "shared-abc",
		To:            "0xConfirmer",
		ChainId:       big.NewInt(137),
		PNTxHash:      "0xconfirmhash",
		PNTxTimestamp: time.Now(),
		TokenInType:   types.DvpERC721,
	}

	mockSwapSalts.EXPECT().
		Get(gomock.Any(), sharedIdKey).
		Return(storedSalts, nil)
	mockDecryptor.EXPECT().
		DecryptWithSalt([]byte("encrypted"), storedSalts.InitiatorSelfSalt).
		DoAndReturn(func(_, _ []byte) ([]byte, error) {
			return json.Marshal(decryptedData)
		})

	mockTxRepo.EXPECT().
		UpdateDvpTeleportConfirmation(
			gomock.Any(),
			"shared-abc",
			gomock.Any(),
			"0xConfirmer",
			"0xconfirmhash",
			gomock.Any(),
			"",
		).
		Return(nil)

	mockTxRepo.EXPECT().
		UpdateTeleportStatusBySharedID(
			gomock.Any(),
			"shared-abc",
			uint8(types.DvpSwapStateCompleted),
		).
		Return(nil)

	handler := &DvpTeleportEventHandler{
		txRepo:    mockTxRepo,
		decryptor: mockDecryptor,
		log:       &testutil.StubLogger{},
		pnData:    &core.PNodeDataAndSecrets{},
		swapSalts: mockSwapSalts,
	}

	log := core.ContractLog{
		RawEventData: testutil.MustMarshal(t, &DvpTeleport.DvpTeleportSwapCompleted{
			SharedId:      sharedIdBytes,
			EncryptedData: []byte("encrypted"),
		}),
		BlockNumber: 1000,
	}

	err := handler.processSwapCompletedEvent(context.Background(), log)

	assert.NoError(t, err)
}

func TestDvpTeleportEventHandler_ProcessSwapCompletedEvent_SkipsWhenSaltsMissing(t *testing.T) {
	// When the listener never saw the matching SwapInitiated (foreign swap or
	// init before starting block) there are no stored salts. The event must
	// be acked, not requeued — mirrors the relayer's `if swap == nil { return nil }`.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSwapSalts := mocks.NewMockSwapSaltsStore(ctrl)
	mockSwapSalts.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(types.DvpSwapSalts{}, core.ErrSwapSaltsNotFound)

	handler := &DvpTeleportEventHandler{
		txRepo:    mocks.NewMockTransactionRepository(ctrl),
		decryptor: mocks.NewMockDecryptor(ctrl),
		log:       &testutil.StubLogger{},
		pnData:    &core.PNodeDataAndSecrets{},
		swapSalts: mockSwapSalts,
	}

	log := core.ContractLog{
		RawEventData: testutil.MustMarshal(t, &DvpTeleport.DvpTeleportSwapCompleted{
			EncryptedData: []byte("encrypted"),
		}),
	}

	err := handler.processSwapCompletedEvent(context.Background(), log)

	assert.NoError(t, err)
}

func TestDvpTeleportEventHandler_ProcessSwapCompletedEvent_ReturnsErrorWhenConfirmationFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTxRepo := mocks.NewMockTransactionRepository(ctrl)
	mockDecryptor := mocks.NewMockDecryptor(ctrl)
	mockSwapSalts := mocks.NewMockSwapSaltsStore(ctrl)

	decryptedData := types.DvpSwapMessage{
		SharedId: "shared-err",
		ChainId:  big.NewInt(1),
	}

	mockSwapSalts.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(types.DvpSwapSalts{InitiatorSelfSalt: []byte("self"), InitiatorCtxtSalt: []byte("ctxt")}, nil)
	mockDecryptor.EXPECT().
		DecryptWithSalt(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_, _ []byte) ([]byte, error) {
			return json.Marshal(decryptedData)
		})
	mockTxRepo.EXPECT().
		UpdateDvpTeleportConfirmation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(assert.AnError)

	handler := &DvpTeleportEventHandler{
		txRepo:    mockTxRepo,
		decryptor: mockDecryptor,
		log:       &testutil.StubLogger{},
		pnData:    &core.PNodeDataAndSecrets{},
		swapSalts: mockSwapSalts,
	}

	log := core.ContractLog{
		RawEventData: testutil.MustMarshal(t, &DvpTeleport.DvpTeleportSwapCompleted{
			EncryptedData: []byte("encrypted"),
		}),
	}

	err := handler.processSwapCompletedEvent(context.Background(), log)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update DVP transactions")
}

func TestDvpTeleportEventHandler_ProcessSwapCancelledEvent_MarksCancelled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTxRepo := mocks.NewMockTransactionRepository(ctrl)

	sharedIdBytes := [32]byte{}
	copy(sharedIdBytes[:], "shared-cancel")
	expectedSharedId := hex.EncodeToString(sharedIdBytes[:])

	mockTxRepo.EXPECT().
		UpdateTeleportStatusBySharedID(gomock.Any(), expectedSharedId, uint8(types.DvpSwapStateCancelled)).
		Return(nil)

	handler := &DvpTeleportEventHandler{txRepo: mockTxRepo, log: &testutil.StubLogger{}}

	log := core.ContractLog{
		RawEventData: testutil.MustMarshal(t, &DvpTeleport.DvpTeleportSwapCancelled{SharedId: sharedIdBytes}),
	}

	err := handler.processSwapCancelledEvent(context.Background(), log)

	assert.NoError(t, err)
}

func TestDvpTeleportEventHandler_ProcessSwapCancelledEvent_PropagatesUpdateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTxRepo := mocks.NewMockTransactionRepository(ctrl)
	mockTxRepo.EXPECT().
		UpdateTeleportStatusBySharedID(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(assert.AnError)

	handler := &DvpTeleportEventHandler{txRepo: mockTxRepo, log: &testutil.StubLogger{}}

	log := core.ContractLog{
		RawEventData: testutil.MustMarshal(t, &DvpTeleport.DvpTeleportSwapCancelled{SharedId: [32]byte{}}),
	}

	err := handler.processSwapCancelledEvent(context.Background(), log)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update teleport_status")
}

func TestDvpTeleportEventHandler_ProcessSwapTimedOutEvent_MarksCancelled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTxRepo := mocks.NewMockTransactionRepository(ctrl)

	sharedIdBytes := [32]byte{}
	copy(sharedIdBytes[:], "shared-timeout")
	expectedSharedId := hex.EncodeToString(sharedIdBytes[:])

	// SwapTimedOut reuses the Cancelled status per product decision.
	mockTxRepo.EXPECT().
		UpdateTeleportStatusBySharedID(gomock.Any(), expectedSharedId, uint8(types.DvpSwapStateCancelled)).
		Return(nil)

	handler := &DvpTeleportEventHandler{txRepo: mockTxRepo, log: &testutil.StubLogger{}}

	log := core.ContractLog{
		RawEventData: testutil.MustMarshal(t, &DvpTeleport.DvpTeleportSwapTimedOut{SharedId: sharedIdBytes}),
	}

	err := handler.processSwapTimedOutEvent(context.Background(), log)

	assert.NoError(t, err)
}

func TestDvpTeleportEventHandler_ProcessERCDvpBalanceUpdatedEvent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTxRepo := mocks.NewMockTransactionRepository(ctrl)
	mockDecryptor := mocks.NewMockDecryptor(ctrl)

	mockDecryptor.EXPECT().
		DecryptPayloadBytes(gomock.Any(), gomock.Any(), gomock.Any(), types.AtomicSecret).
		DoAndReturn(func([]byte, uint64, interface{}, types.SecretType) ([]byte, error) {
			return json.Marshal(types.DvpBalanceUpdated{
				ErcId:              big.NewInt(42),
				TokenType:          uint8(types.AssetTypeDvpERC721),
				ResourceId:         "0xresource",
				Amount:             big.NewInt(500),
				UpdateType:         types.Mint,
				SharedId:           "shared-123",
				From:               "0xFrom",
				To:                 "0xTo",
				SourceChainId:      big.NewInt(1),
				DestinationChainId: big.NewInt(2),
			})
		})

	mockTxRepo.EXPECT().
		CreateTransaction(gomock.Any(), gomock.Any()).
		Return(nil)

	handler := &DvpTeleportEventHandler{
		txRepo:    mockTxRepo,
		decryptor: mockDecryptor,
		log:       &testutil.StubLogger{},
		pnData:    &core.PNodeDataAndSecrets{},
	}

	log := core.ContractLog{
		RawEventData: testutil.MustMarshal(t, &DvpTeleport.DvpTeleportERCDvpBalanceUpdated{
			EncryptedMessage: []byte("encrypted-data"),
		}),
		BlockNumber:     1000,
		TransactionHash: "0xabc123",
		LogIndex:        1,
	}

	err := handler.processERCDvpBalanceUpdatedEvent(context.Background(), log)

	assert.NoError(t, err)
}

func TestDvpTeleportEventHandler_ProcessERCDvpBalanceUpdatedEvent_DecryptionFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDecryptor := mocks.NewMockDecryptor(ctrl)
	mockDecryptor.EXPECT().
		DecryptPayloadBytes(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, errors.New("decryption failed"))

	handler := &DvpTeleportEventHandler{
		txRepo:    mocks.NewMockTransactionRepository(ctrl),
		decryptor: mockDecryptor,
		log:       &testutil.StubLogger{},
		pnData:    &core.PNodeDataAndSecrets{},
	}

	log := core.ContractLog{
		RawEventData: testutil.MustMarshal(t, &DvpTeleport.DvpTeleportERCDvpBalanceUpdated{
			EncryptedMessage: []byte("encrypted-data"),
		}),
		BlockNumber: 1000,
	}

	err := handler.processERCDvpBalanceUpdatedEvent(context.Background(), log)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "decryption failed")
}
