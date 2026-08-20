package handlers

import (
	"context"
	"errors"
	"math/big"
	"testing"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/core"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/events"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/mocks"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/testutil"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/contracts/TokenFreezeManagerV1"
)

func TestTokenFreezeManagerEventHandler_Name(t *testing.T) {
	handler := NewTokenFreezeManagerEventHandler(nil, nil, &testutil.StubLogger{})
	assert.Equal(t, "TokenFreezeManagerHandler", handler.Name())
}

func TestTokenFreezeManagerEventHandler_ContractName(t *testing.T) {
	handler := NewTokenFreezeManagerEventHandler(nil, nil, &testutil.StubLogger{})
	assert.Equal(t, events.ContractTokenFreezeManager, handler.ContractName())
}

func TestTokenFreezeManagerEventHandler_Handle_UnknownEventReturnsNil(t *testing.T) {
	// An unrecognised event name is silently ignored and returns nil
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler := NewTokenFreezeManagerEventHandler(
		mocks.NewMockTokenFreezeRepository(ctrl),
		mocks.NewMockProvider(ctrl),
		&testutil.StubLogger{},
	)

	err := handler.Handle(context.Background(), core.ContractLog{EventName: "UnknownEvent"})
	assert.NoError(t, err)
}

func TestTokenFreezeManagerEventHandler_Handle_InvalidDataReturnsError(t *testing.T) {
	// A TokenFreezeStatusChanged event with nil raw data returns an unmarshal error
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler := NewTokenFreezeManagerEventHandler(
		mocks.NewMockTokenFreezeRepository(ctrl),
		mocks.NewMockProvider(ctrl),
		&testutil.StubLogger{},
	)

	err := handler.Handle(context.Background(), core.ContractLog{
		EventName:    events.TokenFreezeStatusChanged,
		RawEventData: nil,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal event data")
}

func TestTokenFreezeManagerEventHandler_Handle_SuccessCallsRepo(t *testing.T) {
	// A valid TokenFreezeStatusChanged event fetches the block and updates the freeze status
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTokenFreezeRepository(ctrl)
	mockProvider := mocks.NewMockProvider(ctrl)

	block := ethTypes.NewBlock(&ethTypes.Header{Time: 1700000000}, &ethTypes.Body{}, nil, nil)
	mockProvider.EXPECT().GetBlockByNumber(gomock.Any(), uint64(1000)).Return(block, nil)
	mockRepo.EXPECT().
		UpdateTokenFreezeStatus(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	handler := NewTokenFreezeManagerEventHandler(mockRepo, mockProvider, &testutil.StubLogger{})

	err := handler.Handle(context.Background(), core.ContractLog{
		EventName:       events.TokenFreezeStatusChanged,
		BlockNumber:     1000,
		TransactionHash: "0xabc123",
		RawEventData: testutil.MustMarshal(t, &TokenFreezeManagerV1.TokenFreezeManagerV1TokenFreezeStatusChanged{
			ResourceId: [32]byte{1, 2, 3},
			ChainIds:   []*big.Int{big.NewInt(10), big.NewInt(20)},
			Action:     1,
		}),
	})
	assert.NoError(t, err)
}

func TestTokenFreezeManagerEventHandler_ProcessStatusChanged_InvalidDataReturnsError(t *testing.T) {
	// Nil raw event data returns an unmarshal error with a descriptive message
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler := &TokenFreezeManagerEventHandler{
		tokenFreezeRepo: mocks.NewMockTokenFreezeRepository(ctrl),
		provider:        mocks.NewMockProvider(ctrl),
		log:             &testutil.StubLogger{},
	}

	err := handler.processTokenFreezeStatusChanged(context.Background(), core.ContractLog{RawEventData: nil})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal event data for TokenFreezeStatusChanged")
}

func TestTokenFreezeManagerEventHandler_ProcessStatusChanged_GetBlockFailsReturnsError(t *testing.T) {
	// A provider error when fetching the block is propagated as the handler error
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProvider := mocks.NewMockProvider(ctrl)
	mockProvider.EXPECT().GetBlockByNumber(gomock.Any(), uint64(500)).Return(nil, errors.New("rpc error"))

	handler := &TokenFreezeManagerEventHandler{
		tokenFreezeRepo: mocks.NewMockTokenFreezeRepository(ctrl),
		provider:        mockProvider,
		log:             &testutil.StubLogger{},
	}

	err := handler.processTokenFreezeStatusChanged(context.Background(), core.ContractLog{
		BlockNumber:     500,
		TransactionHash: "0xdead",
		RawEventData: testutil.MustMarshal(t, &TokenFreezeManagerV1.TokenFreezeManagerV1TokenFreezeStatusChanged{
			ResourceId: [32]byte{4, 5, 6},
			ChainIds:   []*big.Int{big.NewInt(99)},
			Action:     1,
		}),
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rpc error")
}

func TestTokenFreezeManagerEventHandler_ProcessStatusChanged_RepoFailsReturnsError(t *testing.T) {
	// A repository error when persisting freeze status is propagated as the handler error
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTokenFreezeRepository(ctrl)
	mockProvider := mocks.NewMockProvider(ctrl)

	block := ethTypes.NewBlock(&ethTypes.Header{Time: 1700000001}, &ethTypes.Body{}, nil, nil)
	mockProvider.EXPECT().GetBlockByNumber(gomock.Any(), uint64(600)).Return(block, nil)
	mockRepo.EXPECT().
		UpdateTokenFreezeStatus(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(errors.New("db write error"))

	handler := &TokenFreezeManagerEventHandler{
		tokenFreezeRepo: mockRepo,
		provider:        mockProvider,
		log:             &testutil.StubLogger{},
	}

	err := handler.processTokenFreezeStatusChanged(context.Background(), core.ContractLog{
		BlockNumber:     600,
		TransactionHash: "0xbeef",
		RawEventData: testutil.MustMarshal(t, &TokenFreezeManagerV1.TokenFreezeManagerV1TokenFreezeStatusChanged{
			ResourceId: [32]byte{7, 8, 9},
			ChainIds:   []*big.Int{big.NewInt(11), big.NewInt(22)},
			Action:     1,
		}),
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db write error")
}

func TestTokenFreezeManagerEventHandler_ProcessStatusChanged_FreezeSucceeds(t *testing.T) {
	// A freeze action (action=1) calls UpdateTokenFreezeStatus with the correct resource ID, chain IDs, and block metadata
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTokenFreezeRepository(ctrl)
	mockProvider := mocks.NewMockProvider(ctrl)

	block := ethTypes.NewBlock(&ethTypes.Header{Time: 1700000002}, &ethTypes.Body{}, nil, nil)
	mockProvider.EXPECT().GetBlockByNumber(gomock.Any(), uint64(1000)).Return(block, nil)
	mockRepo.EXPECT().
		UpdateTokenFreezeStatus(
			gomock.Any(),
			"0a0b0c0000000000000000000000000000000000000000000000000000000000",
			[]*big.Int{big.NewInt(100), big.NewInt(200)},
			uint8(1),
			new(big.Int).SetUint64(1000),
			"0xcafe",
			gomock.Any(),
		).
		Return(nil)

	handler := &TokenFreezeManagerEventHandler{
		tokenFreezeRepo: mockRepo,
		provider:        mockProvider,
		log:             &testutil.StubLogger{},
	}

	err := handler.processTokenFreezeStatusChanged(context.Background(), core.ContractLog{
		BlockNumber:     1000,
		TransactionHash: "0xcafe",
		RawEventData: testutil.MustMarshal(t, &TokenFreezeManagerV1.TokenFreezeManagerV1TokenFreezeStatusChanged{
			ResourceId: [32]byte{10, 11, 12},
			ChainIds:   []*big.Int{big.NewInt(100), big.NewInt(200)},
			Action:     1,
		}),
	})
	assert.NoError(t, err)
}

func TestTokenFreezeManagerEventHandler_ProcessStatusChanged_UnfreezeSucceeds(t *testing.T) {
	// An unfreeze action (action=0) calls UpdateTokenFreezeStatus with the correct resource ID, chain IDs, and block metadata
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTokenFreezeRepository(ctrl)
	mockProvider := mocks.NewMockProvider(ctrl)

	block := ethTypes.NewBlock(&ethTypes.Header{Time: 1700000003}, &ethTypes.Body{}, nil, nil)
	mockProvider.EXPECT().GetBlockByNumber(gomock.Any(), uint64(2000)).Return(block, nil)
	mockRepo.EXPECT().
		UpdateTokenFreezeStatus(
			gomock.Any(),
			"0d0e0f0000000000000000000000000000000000000000000000000000000000",
			[]*big.Int{big.NewInt(300)},
			uint8(0),
			new(big.Int).SetUint64(2000),
			"0xf00d",
			gomock.Any(),
		).
		Return(nil)

	handler := &TokenFreezeManagerEventHandler{
		tokenFreezeRepo: mockRepo,
		provider:        mockProvider,
		log:             &testutil.StubLogger{},
	}

	err := handler.processTokenFreezeStatusChanged(context.Background(), core.ContractLog{
		BlockNumber:     2000,
		TransactionHash: "0xf00d",
		RawEventData: testutil.MustMarshal(t, &TokenFreezeManagerV1.TokenFreezeManagerV1TokenFreezeStatusChanged{
			ResourceId: [32]byte{13, 14, 15},
			ChainIds:   []*big.Int{big.NewInt(300)},
			Action:     0,
		}),
	})
	assert.NoError(t, err)
}
