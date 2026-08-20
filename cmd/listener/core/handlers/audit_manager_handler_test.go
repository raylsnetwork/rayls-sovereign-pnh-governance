package handlers

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/core"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/mocks"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/testutil"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/config"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/AuditManagerV1"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/types"
)

func TestProcessNewAuditOrChainInfo(t *testing.T) {
	tests := []struct {
		name           string
		log            core.ContractLog
		setupMocks     func(*testutil.StubConfigProvider, *mocks.MockDecryptor, *AuditManagerEventHandler)
		expectedError  bool
		errorContains  string
		verifyPNData   bool
		expectedPNData core.PNodeDataAndSecrets
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
			name: "GetConfig fails - returns error",
			log: core.ContractLog{
				RawEventData: testutil.MustMarshal(t, &AuditManagerV1.AuditManagerV1NewAuditOrChainInfo{}),
			},
			setupMocks: func(stubConfigProvider *testutil.StubConfigProvider, mockDecryptor *mocks.MockDecryptor, handler *AuditManagerEventHandler) {
				stubConfigProvider.ConfigErr = errors.New("failed to get config")
			},
			expectedError: true,
			errorContains: "failed to get config",
		},
		{
			name: "GatherParticipantsData fails - returns error",
			log: core.ContractLog{
				RawEventData: testutil.MustMarshal(t, &AuditManagerV1.AuditManagerV1NewAuditOrChainInfo{}),
			},
			setupMocks: func(stubConfigProvider *testutil.StubConfigProvider, mockDecryptor *mocks.MockDecryptor, handler *AuditManagerEventHandler) {
				// GetConfig succeeds
				stubConfigProvider.Config = &config.Config{}

				// GatherParticipantsData fails
				mockDecryptor.EXPECT().
					GatherParticipantsData(gomock.Any(), gomock.Not(gomock.Nil())).
					Return(nil, errors.New("failed to gather participants"))
			},
			expectedError: true,
			errorContains: "failed to gather participants",
		},
		{
			name: "success - participant data gathered and stored",
			log: core.ContractLog{
				RawEventData: testutil.MustMarshal(t, &AuditManagerV1.AuditManagerV1NewAuditOrChainInfo{}),
			},
			setupMocks: func(stubConfigProvider *testutil.StubConfigProvider, mockDecryptor *mocks.MockDecryptor, handler *AuditManagerEventHandler) {
				chainId1 := "12345"
				chainId2 := "12346"
				venSecret1 := []byte{1, 2, 3, 4, 5}
				venSecret2 := []byte{6, 7, 8, 9, 10}

				expectedPNData := core.PNodeDataAndSecrets{
					chainId1: map[string]*types.IPNodeDataAndSecrets{
						"100": {
							ChainId:         big.NewInt(12345),
							BlockNumber:     big.NewInt(100),
							HubSharedSecret: venSecret1,
						},
					},
					chainId2: map[string]*types.IPNodeDataAndSecrets{
						"110": {
							ChainId:         big.NewInt(12346),
							BlockNumber:     big.NewInt(110),
							HubSharedSecret: venSecret2,
						},
					},
				}

				// GetConfig succeeds
				stubConfigProvider.Config = &config.Config{}

				mockDecryptor.EXPECT().
					GatherParticipantsData(gomock.Any(), gomock.Not(gomock.Nil())).
					Return(expectedPNData, nil)
			},
			expectedError: false,
			verifyPNData:  true,
			expectedPNData: core.PNodeDataAndSecrets{
				"12345": map[string]*types.IPNodeDataAndSecrets{
					"100": {
						ChainId:         big.NewInt(12345),
						BlockNumber:     big.NewInt(100),
						HubSharedSecret: []byte{1, 2, 3, 4, 5},
					},
				},
				"12346": map[string]*types.IPNodeDataAndSecrets{
					"110": {
						ChainId:         big.NewInt(12346),
						BlockNumber:     big.NewInt(110),
						HubSharedSecret: []byte{6, 7, 8, 9, 10},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			stubConfigProvider := testutil.NewStubConfigProvider()
			mockDecryptor := mocks.NewMockDecryptor(ctrl)

			pnData := make(core.PNodeDataAndSecrets)
			handler := &AuditManagerEventHandler{
				configProvider: stubConfigProvider,
				decryptor:      mockDecryptor,
				log:            &testutil.StubLogger{},
				pnData:         &pnData,
			}

			if tt.setupMocks != nil {
				tt.setupMocks(stubConfigProvider, mockDecryptor, handler)
			}

			// Execute
			err := handler.processNewAuditOrChainInfo(ctx, tt.log)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				if tt.verifyPNData {
					assert.Equal(t, tt.expectedPNData, *handler.pnData)
				}
			}
		})
	}
}

func TestAuditManagerEventHandler_Name(t *testing.T) {
	handler := NewAuditManagerEventHandler(nil, nil, nil, nil)

	name := handler.Name()

	assert.Equal(t, "AuditManagerHandler", name)
}
