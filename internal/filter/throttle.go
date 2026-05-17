package filter

import (
	"fmt"
	"time"

	"github.com/user/logslice/internal/parser"
)

// ThrottleFilter passes at most one entry per window duration.
// Subsequent entries within the same window are dropped.
type ThrottleFilter struct {
	inner    Filter
	window   time.Duration
	lastSeen time.Time
}

// NewThrottleFilter creates a ThrottleFilter that passes at most one entry
// from inner per window duration. window must be positive and inner non-nil.
func NewThrottleFilter(inner Filter, window time.Duration) (*ThrottleFilter, error) {
	if inner == nil {
		return nil, fmt.Errorf("throttle: inner filter must not be nil")
	}
	if window <= 0 {
		return nil, fmt.Errorf("throttle: window must be positive, got %v", window)
	}
	return &ThrottleFilter{inner: inner, window: window}, nil
}

// Matches returns true if the inner filter matches and at least one full window
// has elapsed since the last accepted entry.
func (f *ThrottleFilter) Matches(e parser.Entry) bool {
	if !f.inner.Matches(e) {
		return false
	}
	if f.lastSeen.IsZero() || e.Timestamp.Sub(f.lastSeen) >= f.window {
		f.lastSeen = e.Timestamp
		return true
	}
	return false
}

// String returns a human-readable description of the filter.
func (f *ThrottleFilter) String() string {
	return fmt.Sprintf("throttle(%v, %s)", f.window, f.inner)
}
