package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/core"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/events"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/mocks"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/testutil"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/contracts/EnygmaTokenManagerV1"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
)

func TestEnygmaTokenManagerEventHandler_ProcessTokenRegistered(t *testing.T) {
	// Helper to build a token
	buildToken := func(resourceID string, issuerID string) *domain.Token {
		return &domain.Token{
			ResourceId: resourceID,
			IssuerId:   issuerID,
			Name:       "Test Token",
			Symbol:     "TST",
		}
	}

	tests := []struct {
		name             string
		eventResourceId  string
		eventIssuerId    string
		setupMocks       func(*mocks.MockTokenService, *mocks.MockTokenRepository)
		validateBehavior func(*testing.T, *domain.Token, error)
	}{
		{
			name:            "success - token found, issuer matches, upserted",
			eventResourceId: "abc123",
			eventIssuerId:   "1",
			setupMocks: func(tokenService *mocks.MockTokenService, tokenRepo *mocks.MockTokenRepository) {
				token := buildToken("abc123", "1")

				tokenService.EXPECT().
					GetTokenByResourceId(gomock.Any(), "abc123").
					Return(token, nil)

				tokenRepo.EXPECT().
					Upsert(gomock.Any(), token).
					Return(nil)
			},
			validateBehavior: func(t *testing.T, token *domain.Token, err error) {
				require.NoError(t, err)
				require.NotNil(t, token)
				assert.Equal(t, "abc123", token.ResourceId)
				assert.Equal(t, "1", token.IssuerId)
				assert.Equal(t, "Test Token", token.Name)
			},
		},
		{
			name:            "blockchain returns nil - creates and returns stub token",
			eventResourceId: "nonexistent",
			eventIssuerId:   "1",
			setupMocks: func(tokenService *mocks.MockTokenService, tokenRepo *mocks.MockTokenRepository) {
				tokenService.EXPECT().
					GetTokenByResourceId(gomock.Any(), "nonexistent").
					Return(nil, nil)
				// Stub token must be upserted to keep the tokens table consistent
				tokenRepo.EXPECT().
					Upsert(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			validateBehavior: func(t *testing.T, token *domain.Token, err error) {
				require.NoError(t, err)
				require.NotNil(t, token, "should return a stub token when blockchain has no data")
				assert.Equal(t, "nonexistent", token.ResourceId)
				assert.Equal(t, "1", token.IssuerId)
				assert.Equal(t, "Unknown", token.Name)
			},
		},
		{
			name:            "issuer mismatch - returns ErrIssuerMismatch",
			eventResourceId: "abc123",
			eventIssuerId:   "1",
			setupMocks: func(tokenService *mocks.MockTokenService, tokenRepo *mocks.MockTokenRepository) {
				token := buildToken("abc123", "2")

				tokenService.EXPECT().
					GetTokenByResourceId(gomock.Any(), "abc123").
					Return(token, nil)
			},
			validateBehavior: func(t *testing.T, token *domain.Token, err error) {
				require.ErrorIs(t, err, core.ErrIssuerMismatch)
				assert.Nil(t, token)
			},
		},
		{
			name:            "error - GetTokenByResourceId fails",
			eventResourceId: "abc123",
			eventIssuerId:   "1",
			setupMocks: func(tokenService *mocks.MockTokenService, tokenRepo *mocks.MockTokenRepository) {
				// Mock GetTokenByResourceId to return error
				tokenService.EXPECT().
					GetTokenByResourceId(gomock.Any(), "abc123").
					Return(nil, errors.New("db connection failed"))
			},
			validateBehavior: func(t *testing.T, token *domain.Token, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "db connection failed")
				assert.Nil(t, token)
			},
		},
		{
			name:            "error - Upsert fails",
			eventResourceId: "abc123",
			eventIssuerId:   "1",
			setupMocks: func(tokenService *mocks.MockTokenService, tokenRepo *mocks.MockTokenRepository) {
				token := buildToken("abc123", "1")
				tokenService.EXPECT().
					GetTokenByResourceId(gomock.Any(), "abc123").
					Return(token, nil)

				// Mock Upsert to fail
				tokenRepo.EXPECT().
					Upsert(gomock.Any(), token).
					Return(errors.New("database write failed"))
			},
			validateBehavior: func(t *testing.T, token *domain.Token, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "database write failed")
				assert.Nil(t, token)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTokenService := mocks.NewMockTokenService(ctrl)
			mockTokenRepo := mocks.NewMockTokenRepository(ctrl)
			logger := &testutil.StubLogger{}

			tc.setupMocks(mockTokenService, mockTokenRepo)

			handler := NewEnygmaTokenManagerEventHandler(mockTokenRepo, mockTokenService, logger)

			result, err := handler.processTokenRegistered(context.Background(), tc.eventResourceId, tc.eventIssuerId)

			tc.validateBehavior(t, result, err)
		})
	}
}

func TestEnygmaTokenManagerEventHandler_ProcessEnygmaTokenRegistered(t *testing.T) {
	tests := []struct {
		name             string
		rawEventData     json.RawMessage
		setupMocks       func(*mocks.MockTokenService, *mocks.MockTokenRepository)
		validateBehavior func(*testing.T, error)
	}{
		{
			name: "success - valid event, token found and upserted",
			rawEventData: testutil.MustMarshal(t, &EnygmaTokenManagerV1.EnygmaTokenManagerV1EnygmaTokenRegistered{
				IssuerChainId: big.NewInt(1),
				BlockNumber:   big.NewInt(100),
				Name:          "Test Token",
				InitialSupply: big.NewInt(1000),
			}),
			setupMocks: func(tokenService *mocks.MockTokenService, tokenRepo *mocks.MockTokenRepository) {
				token := &domain.Token{
					IssuerId: "1",
					Name:     "Test Token",
					Symbol:   "TST",
				}

				tokenService.EXPECT().
					GetTokenByResourceId(gomock.Any(), gomock.Any()).
					Return(token, nil)

				tokenRepo.EXPECT().
					Upsert(gomock.Any(), token).
					Return(nil)
			},
			validateBehavior: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "blockchain returns nil - creates stub token",
			rawEventData: testutil.MustMarshal(t, &EnygmaTokenManagerV1.EnygmaTokenManagerV1EnygmaTokenRegistered{
				IssuerChainId: big.NewInt(1),
			}),
			setupMocks: func(tokenService *mocks.MockTokenService, tokenRepo *mocks.MockTokenRepository) {
				tokenService.EXPECT().
					GetTokenByResourceId(gomock.Any(), gomock.Any()).
					Return(nil, nil)
				// Stub token must be upserted to maintain consistency with transaction records
				tokenRepo.EXPECT().
					Upsert(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			validateBehavior: func(t *testing.T, err error) {
				require.NoError(t, err, "should not error when blockchain has no token data")
			},
		},
		{
			name: "issuer mismatch - skips without error",
			rawEventData: testutil.MustMarshal(t, &EnygmaTokenManagerV1.EnygmaTokenManagerV1EnygmaTokenRegistered{
				IssuerChainId: big.NewInt(1),
			}),
			setupMocks: func(tokenService *mocks.MockTokenService, tokenRepo *mocks.MockTokenRepository) {
				token := &domain.Token{
					ResourceId: "616263313233000000000000000000000000000000000000000000000000",
					IssuerId:   "2",
				}

				tokenService.EXPECT().
					GetTokenByResourceId(gomock.Any(), gomock.Any()).
					Return(token, nil)
			},
			validateBehavior: func(t *testing.T, err error) {
				require.NoError(t, err, "issuer mismatch is handled internally and should not propagate")
			},
		},
		{
			name:         "error - type assertion fails",
			rawEventData: nil,
			setupMocks: func(tokenService *mocks.MockTokenService, tokenRepo *mocks.MockTokenRepository) {
			},
			validateBehavior: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "failed to unmarshal event data")
			},
		},
		{
			name: "error - GetTokenByResourceId fails",
			rawEventData: testutil.MustMarshal(t, &EnygmaTokenManagerV1.EnygmaTokenManagerV1EnygmaTokenRegistered{
				IssuerChainId: big.NewInt(1),
			}),
			setupMocks: func(tokenService *mocks.MockTokenService, tokenRepo *mocks.MockTokenRepository) {
				tokenService.EXPECT().
					GetTokenByResourceId(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("db connection failed"))
			},
			validateBehavior: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "db connection failed")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTokenService := mocks.NewMockTokenService(ctrl)
			mockTokenRepo := mocks.NewMockTokenRepository(ctrl)
			logger := &testutil.StubLogger{}

			tc.setupMocks(mockTokenService, mockTokenRepo)

			handler := NewEnygmaTokenManagerEventHandler(mockTokenRepo, mockTokenService, logger)

			// Build the log
			log := core.ContractLog{
				RawEventData: tc.rawEventData,
			}

			err := handler.processEnygmaTokenRegistered(context.Background(), log)

			tc.validateBehavior(t, err)
		})
	}
}

func TestEnygmaTokenManagerEventHandler_Handle(t *testing.T) {
	tests := []struct {
		name             string
		log              core.ContractLog
		setupMocks       func(*mocks.MockTokenService, *mocks.MockTokenRepository)
		validateBehavior func(*testing.T, error)
	}{
		{
			name: "routes EnygmaTokenRegistered to handler",
			log: core.ContractLog{
				EventName: events.EnygmaTokenRegistered,
				RawEventData: testutil.MustMarshal(t, &EnygmaTokenManagerV1.EnygmaTokenManagerV1EnygmaTokenRegistered{
					IssuerChainId: big.NewInt(12345),
				}),
			},
			setupMocks: func(tokenService *mocks.MockTokenService, tokenRepo *mocks.MockTokenRepository) {
				token := &domain.Token{
					IssuerId: "12345",
				}
				tokenService.EXPECT().
					GetTokenByResourceId(gomock.Any(), gomock.Any()).
					Return(token, nil)

				tokenRepo.EXPECT().
					Upsert(gomock.Any(), token).
					Return(nil)
			},
			validateBehavior: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "returns error when processEnygmaTokenRegistered fails",
			log: core.ContractLog{
				EventName: events.EnygmaTokenRegistered,
				RawEventData: testutil.MustMarshal(t, &EnygmaTokenManagerV1.EnygmaTokenManagerV1EnygmaTokenRegistered{
					IssuerChainId: big.NewInt(999),
				}),
			},
			setupMocks: func(tokenService *mocks.MockTokenService, tokenRepo *mocks.MockTokenRepository) {
				// Make GetTokenByResourceId fail to trigger error in processEnygmaTokenRegistered
				tokenService.EXPECT().
					GetTokenByResourceId(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("connection failed"))
			},
			validateBehavior: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "connection failed")
			},
		},
		{
			name: "don't throw error if it is an unknown event",
			log: core.ContractLog{
				EventName:    "UnknownEvent",
				RawEventData: nil,
			},
			setupMocks: func(tokenService *mocks.MockTokenService, tokenRepo *mocks.MockTokenRepository) {},
			validateBehavior: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTokenService := mocks.NewMockTokenService(ctrl)
			mockTokenRepo := mocks.NewMockTokenRepository(ctrl)
			logger := &testutil.StubLogger{}

			tc.setupMocks(mockTokenService, mockTokenRepo)

			handler := NewEnygmaTokenManagerEventHandler(mockTokenRepo, mockTokenService, logger)

			err := handler.Handle(context.Background(), tc.log)

			tc.validateBehavior(t, err)
		})
	}
}

func TestEnygmaTokenManagerEventHandler_Name(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTokenRepo := mocks.NewMockTokenRepository(ctrl)
	mockTokenService := mocks.NewMockTokenService(ctrl)
	logger := &testutil.StubLogger{}

	handler := NewEnygmaTokenManagerEventHandler(mockTokenRepo, mockTokenService, logger)

	name := handler.Name()

	assert.Equal(t, "EnygmaTokenManagerHandler", name)
}

func TestNewEnygmaTokenManagerEventHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTokenRepo := mocks.NewMockTokenRepository(ctrl)
	mockTokenService := mocks.NewMockTokenService(ctrl)
	logger := &testutil.StubLogger{}

	handler := NewEnygmaTokenManagerEventHandler(mockTokenRepo, mockTokenService, logger)

	require.NotNil(t, handler)
	assert.Equal(t, "EnygmaTokenManagerHandler", handler.Name())
}
