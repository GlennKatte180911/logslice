package filter

import (
	"fmt"
	"time"

	"github.com/user/logslice/internal/parser"
)

// UntilFilter passes entries whose timestamp is before or equal to a cutoff,
// then stops matching once the cutoff is exceeded (i.e. it is a one-way gate).
type UntilFilter struct {
	cutoff time.Time
	done   bool
}

// NewUntilFilter creates a UntilFilter that passes entries up to and including
// the given cutoff time. Once an entry exceeds the cutoff the filter stops
// matching for all subsequent entries.
func NewUntilFilter(cutoff time.Time) (*UntilFilter, error) {
	if cutoff.IsZero() {
		return nil, fmt.Errorf("until: cutoff time must not be zero")
	}
	return &UntilFilter{cutoff: cutoff}, nil
}

// Matches returns true while the entry timestamp is not after the cutoff.
// Once an entry exceeds the cutoff the filter latches to false.
func (f *UntilFilter) Matches(e parser.Entry) bool {
	if f.done {
		return false
	}
	if e.Timestamp.After(f.cutoff) {
		f.done = true
		return false
	}
	return true
}

// String returns a human-readable description of the filter.
func (f *UntilFilter) String() string {
	return fmt.Sprintf("until(%s)", f.cutoff.Format(time.RFC3339))
}
