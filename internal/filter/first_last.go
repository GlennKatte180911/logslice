package filter

import (
	"fmt"

	"github.com/user/logslice/internal/parser"
)

// FirstFilter passes only the first N entries that match the inner filter.
type FirstFilter struct {
	inner   Filter
	max     int
	matched int
}

// NewFirstFilter returns a filter that passes only the first n entries
// matched by inner. n must be >= 1.
func NewFirstFilter(inner Filter, n int) (*FirstFilter, error) {
	if inner == nil {
		return nil, fmt.Errorf("first: inner filter must not be nil")
	}
	if n < 1 {
		return nil, fmt.Errorf("first: n must be >= 1, got %d", n)
	}
	return &FirstFilter{inner: inner, max: n}, nil
}

// Matches returns true for the first n entries accepted by the inner filter.
func (f *FirstFilter) Matches(e parser.Entry) bool {
	if f.matched >= f.max {
		return false
	}
	if f.inner.Matches(e) {
		f.matched++
		return true
	}
	return false
}

func (f *FirstFilter) String() string {
	return fmt.Sprintf("first(%d, %s)", f.max, f.inner)
}

// LastFilter buffers all entries from inner and only emits the last N on
// demand. Because log processing is streaming, Matches always returns false;
// callers should use Collect to retrieve the tail entries after processing.
type LastFilter struct {
	inner  Filter
	max    int
	buffer []parser.Entry
}

// NewLastFilter returns a filter that retains the last n entries matched
// by inner. n must be >= 1.
func NewLastFilter(inner Filter, n int) (*LastFilter, error) {
	if inner == nil {
		return nil, fmt.Errorf("last: inner filter must not be nil")
	}
	if n < 1 {
		return nil, fmt.Errorf("last: n must be >= 1, got %d", n)
	}
	return &LastFilter{inner: inner, max: n}, nil
}

// Matches buffers matching entries and always returns false; use Collect.
func (f *LastFilter) Matches(e parser.Entry) bool {
	if f.inner.Matches(e) {
		f.buffer = append(f.buffer, e)
		if len(f.buffer) > f.max {
			f.buffer = f.buffer[len(f.buffer)-f.max:]
		}
	}
	return false
}

// Collect returns the retained tail entries in arrival order.
func (f *LastFilter) Collect() []parser.Entry {
	out := make([]parser.Entry, len(f.buffer))
	copy(out, f.buffer)
	return out
}

func (f *LastFilter) String() string {
	return fmt.Sprintf("last(%d, %s)", f.max, f.inner)
}
