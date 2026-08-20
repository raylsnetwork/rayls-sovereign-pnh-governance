package testutil

// StubLogger is a no-op logger implementation for testing
type StubLogger struct{}

func (l *StubLogger) Debug(msg string, args ...any) {}
func (l *StubLogger) Info(msg string, args ...any)  {}
func (l *StubLogger) Warn(msg string, args ...any)  {}
func (l *StubLogger) Error(msg string, args ...any) {}
