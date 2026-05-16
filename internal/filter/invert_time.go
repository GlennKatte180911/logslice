package filter

import (
	"fmt"
	"time"

	"github.com/user/logslice/internal/parser"
)

// InvertTimeFilter matches entries whose timestamp falls OUTSIDE a given range.
// It is the logical complement of the TimeRange filter.
type InvertTimeFilter struct {
	from time.Time
	to   time.Time
}

// NewInvertTimeFilter returns a filter that passes entries whose timestamp is
// strictly before `from` or strictly after `to`. Both times must be non-zero
// and `from` must not be after `to`.
func NewInvertTimeFilter(from, to time.Time) (*InvertTimeFilter, error) {
	if from.IsZero() {
		return nil, fmt.Errorf("invert time filter: from time must not be zero")
	}
	if to.IsZero() {
		return nil, fmt.Errorf("invert time filter: to time must not be zero")
	}
	if to.Before(from) {
		return nil, fmt.Errorf("invert time filter: to (%s) is before from (%s)", to.Format(time.RFC3339), from.Format(time.RFC3339))
	}
	return &InvertTimeFilter{from: from, to: to}, nil
}

// Matches returns true when the entry's timestamp is outside [from, to].
func (f *InvertTimeFilter) Matches(e parser.Entry) bool {
	return e.Timestamp.Before(f.from) || e.Timestamp.After(f.to)
}

// String returns a human-readable description of the filter.
func (f *InvertTimeFilter) String() string {
	return fmt.Sprintf("InvertTime(before=%s OR after=%s)",
		f.from.Format(time.RFC3339),
		f.to.Format(time.RFC3339),
	)
}
