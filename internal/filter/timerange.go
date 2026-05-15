package filter

import (
	"fmt"
	"time"

	"github.com/user/logslice/internal/parser"
)

// TimeRange represents an inclusive time window for filtering log entries.
type TimeRange struct {
	From time.Time
	To   time.Time
}

// NewTimeRange constructs a TimeRange from two RFC3339 timestamp strings.
// Either bound may be empty to indicate an open-ended range.
func NewTimeRange(from, to string) (TimeRange, error) {
	var tr TimeRange
	var err error

	if from != "" {
		tr.From, err = time.Parse(time.RFC3339, from)
		if err != nil {
			return TimeRange{}, fmt.Errorf("invalid --from timestamp: %w", err)
		}
	}

	if to != "" {
		tr.To, err = time.Parse(time.RFC3339, to)
		if err != nil {
			return TimeRange{}, fmt.Errorf("invalid --to timestamp: %w", err)
		}
	}

	if !tr.From.IsZero() && !tr.To.IsZero() && tr.To.Before(tr.From) {
		return TimeRange{}, fmt.Errorf("--to must not be before --from")
	}

	return tr, nil
}

// Contains reports whether the given log entry falls within the time range.
// A zero From or To bound is treated as unbounded on that side.
func (tr TimeRange) Contains(e parser.Entry) bool {
	if !tr.From.IsZero() && e.Timestamp.Before(tr.From) {
		return false
	}
	if !tr.To.IsZero() && e.Timestamp.After(tr.To) {
		return false
	}
	return true
}

// Apply filters a slice of entries, returning only those within the range.
func (tr TimeRange) Apply(entries []parser.Entry) []parser.Entry {
	out := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		if tr.Contains(e) {
			out = append(out, e)
		}
	}
	return out
}
