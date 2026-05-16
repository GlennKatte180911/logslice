package filter

import (
	"fmt"
	"time"

	"github.com/nicholasgasior/logslice/internal/parser"
)

// BurstFilter drops entries that arrive in a burst — more than maxCount
// entries within the given window duration. Once the burst threshold is
// exceeded the filter stops matching until the window resets.
type BurstFilter struct {
	inner    Filter
	maxCount int
	window   time.Duration
	windowStart time.Time
	count    int
}

// NewBurstFilter returns a BurstFilter that wraps inner and allows at most
// maxCount matching entries per window. Returns an error if inner is nil,
// maxCount < 1, or window <= 0.
func NewBurstFilter(inner Filter, maxCount int, window time.Duration) (*BurstFilter, error) {
	if inner == nil {
		return nil, fmt.Errorf("burst: inner filter must not be nil")
	}
	if maxCount < 1 {
		return nil, fmt.Errorf("burst: maxCount must be >= 1, got %d", maxCount)
	}
	if window <= 0 {
		return nil, fmt.Errorf("burst: window must be positive, got %s", window)
	}
	return &BurstFilter{
		inner:    inner,
		maxCount: maxCount,
		window:   window,
	}, nil
}

// Matches returns true when the inner filter matches and the burst threshold
// has not been exceeded within the current window.
func (f *BurstFilter) Matches(e parser.Entry) bool {
	if !f.inner.Matches(e) {
		return false
	}
	if f.windowStart.IsZero() || e.Timestamp.Sub(f.windowStart) >= f.window {
		f.windowStart = e.Timestamp
		f.count = 0
	}
	f.count++
	return f.count <= f.maxCount
}

// String returns a human-readable description of the filter.
func (f *BurstFilter) String() string {
	return fmt.Sprintf("burst(max=%d per %s, inner=%s)", f.maxCount, f.window, f.inner)
}
