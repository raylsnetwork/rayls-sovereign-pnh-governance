package testutil

import (
	"context"
	"errors"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/dto"
)

// StubLogger is a no-op logger implementation for testing
type StubLogger struct{}

func (l *StubLogger) Debug(msg string, args ...any) {}
func (l *StubLogger) Info(msg string, args ...any)  {}
func (l *StubLogger) Warn(msg string, args ...any)  {}
func (l *StubLogger) Error(msg string, args ...any) {}

// NoopLogger is an alias for StubLogger (for backwards compatibility)
type NoopLogger = StubLogger

func NewNoopLogger() *NoopLogger { return &NoopLogger{} }

// StubTokenMetadataService is a configurable metadata service for testing
type StubTokenMetadataService struct {
	Metadata    *dto.TokenMetadataInfoDto // Metadata to return
	ShouldError bool                      // If true, returns an error
}

func (s *StubTokenMetadataService) GetMetadata(
	ctx context.Context,
	baseURL, ercId string,
) (*dto.TokenMetadataInfoDto, error) {
	if s.ShouldError {
		return nil, errors.New("metadata fetch failed")
	}
	if s.Metadata != nil {
		return s.Metadata, nil
	}
	return &dto.TokenMetadataInfoDto{}, nil
}
