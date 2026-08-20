package utils

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// ParseTime parses a timestamp string that can be either:
// - Unix ts: "1762202349"
// - Date only: "2025-12-16"
// - ISO8601 with time: "2025-12-16T10:30:00Z"
func ParseTime(input string) (time.Time, error) {
	// Try Unix timestamp first
	if unix, err := strconv.ParseInt(input, 10, 64); err == nil {
		return time.Unix(unix, 0).UTC(), nil
	}

	// Try date only format (YYYY-MM-DD)
	if t, err := time.Parse("2006-01-02", input); err == nil {
		return t.UTC(), nil
	}

	// Try ISO8601 with time (RFC3339)
	return time.Parse(time.RFC3339, input)
}

// ValidateTimestampRange validates timestamp string filters (after/before)
// Returns parsed times and error if validation fails
func ValidateTimestampRange(after, before string) (afterTime, beforeTime time.Time, err error) {
	now := time.Now().UTC()

	// Validate 'after' timestamp if provided
	if after != "" {
		afterTime, err = ParseTime(after)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf(
				"invalid 'after' timestamp format. Use Unix seconds, date only (YYYY-MM-DD) or ISO8601 (YYYY-MM-DDTHH:MM:SSZ): %w",
				err,
			)
		}
		if afterTime.After(now) {
			return time.Time{}, time.Time{}, errors.New("'after' timestamp cannot be in the future")
		}
	}

	// Validate 'before' timestamp if provided
	if before != "" {
		beforeTime, err = ParseTime(before)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf(
				"invalid 'before' timestamp format. Use Unix seconds, date only (YYYY-MM-DD) or ISO8601 (YYYY-MM-DDTHH:MM:SSZ): %w",
				err,
			)
		}
		if beforeTime.After(now) {
			return time.Time{}, time.Time{}, errors.New("'before' timestamp cannot be in the future")
		}
	}

	// Validate date range is valid when both are provided
	if after != "" && before != "" {
		if !afterTime.Before(beforeTime) {
			return time.Time{}, time.Time{}, errors.New("'after' timestamp must be earlier than 'before' timestamp")
		}
	}

	return afterTime, beforeTime, nil
}
