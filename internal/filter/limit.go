package filter

import (
	"fmt"

	"github.com/user/logslice/internal/parser"
)

// LimitFilter passes only the first N matching entries from its inner filter.
type LimitFilter struct {
	inner   Filter
	max     int
	matched int
}

// NewLimitFilter creates a LimitFilter that forwards at most max entries
// accepted by inner. max must be greater than zero.
func NewLimitFilter(inner Filter, max int) (*LimitFilter, error) {
	if inner == nil {
		return nil, fmt.Errorf("limit: inner filter must not be nil")
	}
	if max <= 0 {
		return nil, fmt.Errorf("limit: max must be greater than zero, got %d", max)
	}
	return &LimitFilter{inner: inner, max: max}, nil
}

// Matches returns true for the first max entries that the inner filter accepts.
// Once the limit is reached every subsequent entry is rejected.
func (f *LimitFilter) Matches(e parser.Entry) bool {
	if f.matched >= f.max {
		return false
	}
	if f.inner.Matches(e) {
		f.matched++
		return true
	}
	return false
}

// Reset sets the internal match counter back to zero so the filter can be
// reused across multiple passes.
func (f *LimitFilter) Reset() {
	f.matched = 0
}

// Matched returns the number of entries accepted so far.
func (f *LimitFilter) Matched() int {
	return f.matched
}

// String returns a human-readable description of the filter.
func (f *LimitFilter) String() string {
	return fmt.Sprintf("limit(%d, %s)", f.max, f.inner)
}
