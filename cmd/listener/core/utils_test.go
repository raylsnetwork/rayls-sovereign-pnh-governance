package core

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/mocks"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
)

// TestBuildTransactionMap tests the BuildTransactionMap function with various scenarios
func TestBuildTransactionMap(t *testing.T) {
	// Sample data structure for testing
	type TestData struct {
		ID   string
		Name string
	}

	// Extractor function for tests
	idExtractor := func(data TestData) string {
		return data.ID
	}

	tests := []struct {
		name          string
		decryptedData []TestData
		idType        TransactionIDType
		setupMock     func(*mocks.MockTransactionRepository)
		expectedError bool
		errorContains string
		validate      func(t *testing.T, txMap map[string]*domain.Transaction)
	}{
		{
			name:          "empty data returns empty map",
			decryptedData: []TestData{},
			idType:        MessageID,
			setupMock:     func(m *mocks.MockTransactionRepository) {},
			validate: func(t *testing.T, txMap map[string]*domain.Transaction) {
				assert.Empty(t, txMap)
			},
		},
		{
			name: "builds map by message_id successfully",
			decryptedData: []TestData{
				{ID: "msg-123", Name: "test1"},
				{ID: "msg-456", Name: "test2"},
			},
			idType: MessageID,
			setupMock: func(m *mocks.MockTransactionRepository) {
				m.EXPECT().
					GetTransactionsByMessageIDs(gomock.Any(), []string{"msg-123", "msg-456"}, false).
					Return([]domain.Transaction{
						{MessageId: "msg-123", SharedId: "shared-1"},
						{MessageId: "msg-456", SharedId: "shared-2"},
					}, nil)
			},
			validate: func(t *testing.T, txMap map[string]*domain.Transaction) {
				// Validate map structure
				assert.Len(t, txMap, 2)
				assert.Contains(t, txMap, "msg-123")
				assert.Contains(t, txMap, "msg-456")

				// Explicitly validate transaction field values
				assert.Equal(t, "msg-123", txMap["msg-123"].MessageId)
				assert.Equal(t, "shared-1", txMap["msg-123"].SharedId)
				assert.Equal(t, "msg-456", txMap["msg-456"].MessageId)
				assert.Equal(t, "shared-2", txMap["msg-456"].SharedId)
			},
		},
		{
			name: "builds map by shared_id successfully",
			decryptedData: []TestData{
				{ID: "shared-abc", Name: "test1"},
				{ID: "shared-def", Name: "test2"},
			},
			idType: SharedID,
			setupMock: func(m *mocks.MockTransactionRepository) {
				m.EXPECT().
					GetTransactionsBySharedIDs(gomock.Any(), []string{"shared-abc", "shared-def"}, false).
					Return([]domain.Transaction{
						{MessageId: "msg-1", SharedId: "shared-abc"},
						{MessageId: "msg-2", SharedId: "shared-def"},
					}, nil)
			},
			validate: func(t *testing.T, txMap map[string]*domain.Transaction) {
				// Validate map structure
				assert.Len(t, txMap, 2)
				assert.Contains(t, txMap, "shared-abc")
				assert.Contains(t, txMap, "shared-def")

				// Explicitly validate transaction field values
				assert.Equal(t, "shared-abc", txMap["shared-abc"].SharedId)
				assert.Equal(t, "msg-1", txMap["shared-abc"].MessageId)
				assert.Equal(t, "shared-def", txMap["shared-def"].SharedId)
				assert.Equal(t, "msg-2", txMap["shared-def"].MessageId)
			},
		},
		{
			name: "unsupported identifier returns error",
			decryptedData: []TestData{
				{ID: "test-123", Name: "test1"},
			},
			idType:        TransactionIDType(999),
			setupMock:     func(m *mocks.MockTransactionRepository) {},
			expectedError: true,
			errorContains: "unsupported identifier type",
		},
		{
			name: "repository error for message_id",
			decryptedData: []TestData{
				{ID: "msg-123", Name: "test1"},
			},
			idType: MessageID,
			setupMock: func(m *mocks.MockTransactionRepository) {
				m.EXPECT().
					GetTransactionsByMessageIDs(gomock.Any(), []string{"msg-123"}, false).
					Return(nil, errors.New("database connection failed"))
			},
			expectedError: true,
			errorContains: "failed to fetch transactions: database connection failed",
		},
		{
			name: "repository error for shared_id",
			decryptedData: []TestData{
				{ID: "shared-123", Name: "test1"},
			},
			idType: SharedID,
			setupMock: func(m *mocks.MockTransactionRepository) {
				m.EXPECT().
					GetTransactionsBySharedIDs(gomock.Any(), []string{"shared-123"}, false).
					Return(nil, errors.New("query timeout"))
			},
			expectedError: true,
			errorContains: "failed to fetch transactions: query timeout",
		},
		{
			name: "repository returns an empty slice",
			decryptedData: []TestData{
				{ID: "msg-123", Name: "test1"},
			},
			idType: MessageID,
			setupMock: func(m *mocks.MockTransactionRepository) {
				m.EXPECT().
					GetTransactionsByMessageIDs(gomock.Any(), []string{"msg-123"}, false).
					Return([]domain.Transaction{}, nil)
			},
			validate: func(t *testing.T, txMap map[string]*domain.Transaction) {
				assert.Empty(t, txMap)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockTransactionRepository(ctrl)
			tt.setupMock(mockRepo)

			ctx := context.Background()
			txMap, err := BuildTransactionMap(ctx, mockRepo, tt.decryptedData, idExtractor, tt.idType)

			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				assert.Nil(t, txMap)
			} else {
				require.NoError(t, err)
				require.NotNil(t, txMap)
				if tt.validate != nil {
					tt.validate(t, txMap)
				}
			}
		})
	}
}
