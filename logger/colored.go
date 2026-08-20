package logger

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"sync"

	"github.com/fatih/color"
)

// ColoredHandler is a slog.Handler that outputs colored text logs
type ColoredHandler struct {
	opts               Options
	preformattedGroups string
	preformattedAttrs  []byte

	projectRoot string

	mu  *sync.Mutex
	out io.Writer
}

// Options for the ColoredHandler
type Options struct {
	Level slog.Leveler
}

// NewColoredHandler creates a new ColoredHandler
func NewColoredHandler(out io.Writer, opts *Options) *ColoredHandler {
	projectRoot, err := getProjectRoot()
	if err != nil {
		projectRoot = "/"
	}

	h := &ColoredHandler{
		out:         out,
		mu:          new(sync.Mutex),
		projectRoot: projectRoot,
	}

	if opts != nil {
		h.opts = *opts
	}
	if h.opts.Level == nil {
		h.opts.Level = slog.LevelInfo
	}
	return h
}

// Enabled reports whether the handler handles records at the given level
func (h *ColoredHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

// Handle handles the Record
func (h *ColoredHandler) Handle(_ context.Context, r slog.Record) error {
	buf := make([]byte, 0, 1024)

	// Time
	if !r.Time.IsZero() {
		buf = r.Time.AppendFormat(buf, "15:04:05 2006-01-02")
		buf = append(buf, ' ')
	}

	// Source location
	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		if f.File != "" {
			src := slog.Source{
				Function: f.Function,
				File:     f.File,
				Line:     f.Line,
			}
			buf = h.appendSource(buf, src)
			buf = append(buf, ' ')
		}
	}

	// Level (colored)
	buf = h.appendLevel(buf, r.Level)
	buf = append(buf, ' ')

	// Message
	buf = append(buf, r.Message...)
	buf = append(buf, " | "...)

	// Preformatted attrs (from WithAttrs)
	if len(h.preformattedAttrs) != 0 {
		buf = append(buf, h.preformattedAttrs...)
	}

	// Record attrs
	r.Attrs(func(a slog.Attr) bool {
		buf = h.appendAttr(buf, a, h.preformattedGroups)
		return true
	})

	buf = append(buf, '\n')

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.out.Write(buf)
	return err
}

// WithGroup returns a new Handler with the given group appended to the receiver's existing groups
func (h *ColoredHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	h2 := *h
	h2.mu = new(sync.Mutex) // New mutex for the copy

	pre := fmt.Sprintf("%s%s.", h.preformattedGroups, name)
	h2.preformattedGroups = pre

	return &h2
}

// WithAttrs returns a new Handler whose attributes consist of both the receiver's attributes and the arguments
func (h *ColoredHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	h2 := *h
	h2.mu = new(sync.Mutex) // New mutex for the copy

	pre := slices.Clip(h.preformattedAttrs)
	for _, attr := range attrs {
		pre = h2.appendAttr(pre, attr, h.preformattedGroups)
	}
	h2.preformattedAttrs = pre

	return &h2
}

// appendAttr appends a formatted attribute to the buffer
func (h *ColoredHandler) appendAttr(buf []byte, a slog.Attr, prefix string) []byte {
	if a.Equal(slog.Attr{}) {
		return buf
	}
	if len(buf) != 0 {
		buf = append(buf, ' ')
	}

	a.Value = a.Value.Resolve()

	// Special handling for errors - format with %+v to get stack traces
	if a.Value.Kind() == slog.KindAny {
		if err, ok := a.Value.Any().(error); ok {
			formatted := fmt.Sprintf("%+v", err)
			buf = fmt.Appendf(buf, "%s=%s", color.CyanString(prefix+a.Key), color.RedString(formatted))
			return buf
		}
	}

	switch a.Value.Kind() {
	case slog.KindTime:
		buf = fmt.Appendf(buf, "%s=%s", color.CyanString(prefix+a.Key), a.Value.Time().Format("15:04:05 2006-01-02"))
	case slog.KindString:
		val := a.Value.String()
		if a.Key == "error" {
			buf = fmt.Appendf(buf, "%s=%s", color.CyanString(prefix+a.Key), color.RedString("%q", val))
		} else {
			buf = fmt.Appendf(buf, "%s=%q", color.CyanString(prefix+a.Key), val)
		}
	case slog.KindGroup:
		attrs := a.Value.Group()
		if len(attrs) == 0 {
			return buf
		}
		groupPrefix := prefix
		if a.Key != "" {
			groupPrefix = fmt.Sprintf("%s%s.", prefix, a.Key)
		}
		for _, ga := range attrs {
			buf = h.appendAttr(buf, ga, groupPrefix)
		}
	default:
		buf = fmt.Appendf(buf, "%s=%s", color.CyanString(prefix+a.Key), a.Value)
	}

	return buf
}

// appendSource appends the source location to the buffer
func (h *ColoredHandler) appendSource(buf []byte, src slog.Source) []byte {
	relativePath, err := filepath.Rel(h.projectRoot, src.File)
	if err != nil {
		relativePath = src.File
	}
	fileBuf := fmt.Sprintf("%s:%d", relativePath, src.Line)

	faintWhite := color.New(color.FgWhite, color.Faint)
	buf = append(buf, faintWhite.Sprintf("%s", fileBuf)...)
	return buf
}

// appendLevel appends the colored level to the buffer
func (h *ColoredHandler) appendLevel(buf []byte, level slog.Level) []byte {
	switch level {
	case slog.LevelDebug:
		buf = append(buf, color.MagentaString(level.String())...)
	case slog.LevelInfo:
		buf = append(buf, color.BlueString(level.String())...)
	case slog.LevelWarn:
		buf = append(buf, color.YellowString(level.String())...)
	case slog.LevelError:
		buf = append(buf, color.RedString(level.String())...)
	default:
		buf = append(buf, level.String()...)
	}
	return buf
}

// getProjectRoot finds the project root by looking for go.mod
func getProjectRoot() (string, error) {
	path, err := filepath.Abs(".")
	if err != nil {
		return "", err
	}

	for path != "/" {
		entries, err := os.ReadDir(path)
		if err != nil {
			return "", err
		}

		for _, e := range entries {
			if e.Name() == "go.mod" {
				return path, nil
			}
		}
		path = filepath.Dir(path)
	}

	return "", errors.New("didn't find go.mod file in any parent directory")
}
