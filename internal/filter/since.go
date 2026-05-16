package filter

import (
	"fmt"
	"time"

	"github.com/user/logslice/internal/parser"
)

// SinceFilter passes entries whose timestamp is within a duration from now.
type SinceFilter struct {
	duration time.Duration
	now      func() time.Time
}

// NewSinceFilter creates a SinceFilter that passes entries within d from now.
// Returns an error if d is zero or negative.
func NewSinceFilter(d time.Duration) (*SinceFilter, error) {
	if d <= 0 {
		return nil, fmt.Errorf("since: duration must be positive, got %v", d)
	}
	return &SinceFilter{
		duration: d,
		now:      time.Now,
	}, nil
}

// Matches returns true if the entry's timestamp is within the filter's duration
// looking back from the current time.
func (f *SinceFilter) Matches(e parser.Entry) bool {
	cutoff := f.now().Add(-f.duration)
	return !e.Timestamp.Before(cutoff)
}

// String returns a human-readable description of the filter.
func (f *SinceFilter) String() string {
	return fmt.Sprintf("since(%v)", f.duration)
}
