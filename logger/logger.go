package logger

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/config"
)

var logger *slog.Logger

func init() {
	logger = slog.Default()
}

// InitializeLogger sets up the global logger based on configuration
func InitializeLogger(conf *config.Config) {
	level := getLogLevel(conf.Logging)

	var handler slog.Handler
	if conf.Logging == "production" {
		// JSON handler for production
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   true,
			Level:       level,
			ReplaceAttr: replaceAttr,
		})
	} else {
		// Colored text handler for development
		handler = NewColoredHandler(os.Stdout, &Options{Level: level})
	}

	logger = slog.New(handler)
	slog.SetDefault(logger)
}

// getLogLevel converts config string to slog.Level
func getLogLevel(levelStr string) slog.Level {
	switch levelStr {
	case "development", "debug":
		return slog.LevelDebug
	default:
		return slog.LevelInfo
	}
}

// log is the internal logging function that handles caller depth correctly
func log(level slog.Level, msg string, args ...any) {
	if !logger.Enabled(context.Background(), level) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // Skip log, public function, and runtime.Callers
	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.Add(args...)
	_ = logger.Handler().Handle(context.Background(), r)
}

// Info logs at INFO level
func Info(msg string, args ...any) {
	log(slog.LevelInfo, msg, args...)
}

// Debug logs at DEBUG level
func Debug(msg string, args ...any) {
	log(slog.LevelDebug, msg, args...)
}

// Error logs at ERROR level
func Error(msg string, args ...any) {
	log(slog.LevelError, msg, args...)
}

// Warn logs at WARN level
func Warn(msg string, args ...any) {
	log(slog.LevelWarn, msg, args...)
}

// Logger defines the interface for logging operations across all services
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// loggerImpl implements Logger using package-level functions
type loggerImpl struct{}

// NewLogger creates a new Logger instance
func NewLogger() Logger {
	return &loggerImpl{}
}

// logWithDepth is used by loggerImpl to log with correct caller depth
func logWithDepth(level slog.Level, msg string, args ...any) {
	if !logger.Enabled(context.Background(), level) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // Skip runtime.Callers, logWithDepth, loggerImpl method
	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.Add(args...)
	_ = logger.Handler().Handle(context.Background(), r)
}

func (l *loggerImpl) Debug(msg string, args ...any) { logWithDepth(slog.LevelDebug, msg, args...) }
func (l *loggerImpl) Info(msg string, args ...any)  { logWithDepth(slog.LevelInfo, msg, args...) }
func (l *loggerImpl) Warn(msg string, args ...any)  { logWithDepth(slog.LevelWarn, msg, args...) }
func (l *loggerImpl) Error(msg string, args ...any) { logWithDepth(slog.LevelError, msg, args...) }

// replaceAttr customizes attribute formatting for JSON output
func replaceAttr(_ []string, a slog.Attr) slog.Attr {
	if a.Value.Kind() == slog.KindAny {
		if v, ok := a.Value.Any().(error); ok {
			a.Value = fmtErr(v)
		}
	}
	if a.Key == slog.SourceKey {
		if s, ok := a.Value.Any().(*slog.Source); ok {
			// Show relative path in JSON output too
			if projectRoot, err := getProjectRoot(); err == nil {
				if rel, err := filepath.Rel(projectRoot, s.File); err == nil {
					s.File = rel
				}
			}
		}
	}

	return a
}

// fmtErr formats an error for JSON logging.
// For withstack.WithStackError, uses .Error() which includes the stack trace.
// For other errors, %+v extracts stack traces if available.
func fmtErr(err error) slog.Value {
	return slog.StringValue(err.Error())
}
