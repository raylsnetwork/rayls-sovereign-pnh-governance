package handlers

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/core"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/mocks"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/testutil"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/ParticipantCoreV1"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/domain"
)

func TestProcessParticipantRegistered(t *testing.T) {
	tests := []struct {
		name          string
		log           core.ContractLog
		setupMocks    func(*mocks.MockParticipantRepository)
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
			name: "chainId exceeds uint64 max - returns error",
			log: core.ContractLog{
				RawEventData: testutil.MustMarshal(t, &ParticipantCoreV1.ParticipantCoreV1ParticipantRegistered{
					Participant: ParticipantCoreV1.ParticipantStructsParticipant{
						ChainId: new(big.Int).Add(
							// set max uint64 + 1
							new(big.Int).SetUint64(^uint64(0)),
							big.NewInt(1),
						),
					},
				}),
			},
			expectedError: true,
			errorContains: "chainId exceeds uint64 max",
		},
		{
			name: "repository upsert fails - returns error",
			log: core.ContractLog{
				RawEventData: testutil.MustMarshal(t, &ParticipantCoreV1.ParticipantCoreV1ParticipantRegistered{
					Participant: ParticipantCoreV1.ParticipantStructsParticipant{
						ChainId:            big.NewInt(12345),
						OwnerId:            "0x123",
						Name:               "Node 1",
						Status:             1,
						Role:               1,
						AllowedToBroadcast: true,
						CreatedAt:          big.NewInt(1000),
						UpdatedAt:          big.NewInt(1000),
					},
				}),
			},
			setupMocks: func(mockParticipantRepo *mocks.MockParticipantRepository) {
				mockParticipantRepo.EXPECT().
					Upsert(gomock.Not(gomock.Nil()), gomock.Not(gomock.Nil())).
					Return(errors.New("db connection failed"))
			},
			expectedError: true,
			errorContains: "failed to upsert new participant registration",
		},
		{
			name: "success - participant registered with correct transformations",
			log: core.ContractLog{
				RawEventData: testutil.MustMarshal(t, &ParticipantCoreV1.ParticipantCoreV1ParticipantRegistered{
					Participant: ParticipantCoreV1.ParticipantStructsParticipant{
						ChainId:            big.NewInt(12345),
						OwnerId:            "0xABC123",
						Name:               "Auditor",
						Status:             1,
						Role:               2,
						AllowedToBroadcast: true,
						CreatedAt:          big.NewInt(1000),
					},
				}),
			},
			setupMocks: func(mockParticipantRepo *mocks.MockParticipantRepository) {
				mockParticipantRepo.EXPECT().
					Upsert(gomock.Not(gomock.Nil()), gomock.Not(gomock.Nil())).
					Do(func(ctx context.Context, p domain.Participant) {
						assert.NotNil(t, p.ChainId)
						assert.Equal(t, uint(12345), *p.ChainId)
						assert.Equal(t, time.Unix(1000, 0).UTC(), p.CreatedAt)
						assert.Equal(t, "active", p.StatusStr)
						assert.Equal(t, "auditor", p.RoleStr)
						assert.Equal(t, "0xABC123", p.OwnerId)
						assert.Equal(t, "Auditor", p.Name)
						assert.Equal(t, uint8(1), p.Status)
						assert.Equal(t, uint8(2), p.Role)
						assert.Equal(t, true, p.AllowedToBroadcast)
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
			mockParticipantRepo := mocks.NewMockParticipantRepository(ctrl)

			handler := &ParticipantCoreEventHandler{
				participantRepo: mockParticipantRepo,
				log:             &testutil.StubLogger{},
			}

			if tt.setupMocks != nil {
				tt.setupMocks(mockParticipantRepo)
			}

			// Execute
			err := handler.processParticipantRegistered(ctx, tt.log)

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

func TestProcessParticipantUpdated(t *testing.T) {
	tests := []struct {
		name          string
		log           core.ContractLog
		setupMocks    func(*mocks.MockParticipantRepository)
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
			name: "chainId exceeds uint64 max - returns error",
			log: core.ContractLog{
				RawEventData: testutil.MustMarshal(t, &ParticipantCoreV1.ParticipantCoreV1ParticipantUpdated{
					Participant: ParticipantCoreV1.ParticipantStructsParticipant{
						ChainId: new(big.Int).Add(
							// set max uint64 + 1
							new(big.Int).SetUint64(^uint64(0)),
							big.NewInt(1),
						),
					},
				}),
			},
			expectedError: true,
			errorContains: "chainId exceeds uint64 max",
		},
		{
			name: "repository upsert fails - returns error",
			log: core.ContractLog{
				RawEventData: testutil.MustMarshal(t, &ParticipantCoreV1.ParticipantCoreV1ParticipantUpdated{
					Participant: ParticipantCoreV1.ParticipantStructsParticipant{
						ChainId:            big.NewInt(12345),
						OwnerId:            "0x123",
						Name:               "Node A",
						Status:             2,
						Role:               1,
						AllowedToBroadcast: false,
						UpdatedAt:          big.NewInt(1640000200),
					},
				}),
			},
			setupMocks: func(mockParticipantRepo *mocks.MockParticipantRepository) {
				mockParticipantRepo.EXPECT().
					Upsert(gomock.Not(gomock.Nil()), gomock.Not(gomock.Nil())).
					Return(errors.New("db connection failed"))
			},
			expectedError: true,
			errorContains: "failed to upsert participant update",
		},
		{
			name: "success - participant updated",
			log: core.ContractLog{
				RawEventData: testutil.MustMarshal(t, &ParticipantCoreV1.ParticipantCoreV1ParticipantUpdated{
					Participant: ParticipantCoreV1.ParticipantStructsParticipant{
						ChainId:            big.NewInt(54321),
						OwnerId:            "0xDEF456",
						Name:               "Node A",
						Status:             2,
						Role:               1,
						AllowedToBroadcast: false,
						UpdatedAt:          big.NewInt(1640000200),
					},
				}),
			},
			setupMocks: func(mockParticipantRepo *mocks.MockParticipantRepository) {
				mockParticipantRepo.EXPECT().
					Upsert(gomock.Not(gomock.Nil()), gomock.Not(gomock.Nil())).
					Do(func(ctx context.Context, p domain.Participant) {
						assert.NotNil(t, p.ChainId)
						assert.Equal(t, uint(54321), *p.ChainId)
						assert.Equal(t, time.Unix(1640000200, 0).UTC(), p.UpdatedAt)
						assert.Equal(t, "0xDEF456", p.OwnerId)
						assert.Equal(t, "Node A", p.Name)
						assert.Equal(t, uint8(2), p.Status)
						assert.Equal(t, uint8(1), p.Role)
						assert.Equal(t, false, p.AllowedToBroadcast)
						assert.Equal(t, "inactive", p.StatusStr)
						assert.Equal(t, "issuer", p.RoleStr)
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
			mockParticipantRepo := mocks.NewMockParticipantRepository(ctrl)

			handler := &ParticipantCoreEventHandler{
				participantRepo: mockParticipantRepo,
				log:             &testutil.StubLogger{},
			}

			if tt.setupMocks != nil {
				tt.setupMocks(mockParticipantRepo)
			}

			// Execute
			err := handler.processParticipantUpdated(ctx, tt.log)

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

func TestParticipantCoreEventHandler_Name(t *testing.T) {
	handler := NewParticipantCoreEventHandler(nil, nil)

	name := handler.Name()

	assert.Equal(t, "ParticipantCoreHandler", name)
}
