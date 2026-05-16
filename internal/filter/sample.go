package filter

import (
	"fmt"

	"github.com/iand/logslice/internal/parser"
)

// SampleFilter passes every Nth matching entry from an inner filter.
type SampleFilter struct {
	inner  Filter
	n      int
	count  int
}

// NewSampleFilter creates a SampleFilter that passes every nth entry accepted
// by inner. n must be >= 1.
func NewSampleFilter(inner Filter, n int) (*SampleFilter, error) {
	if inner == nil {
		return nil, fmt.Errorf("sample: inner filter must not be nil")
	}
	if n < 1 {
		return nil, fmt.Errorf("sample: n must be >= 1, got %d", n)
	}
	return &SampleFilter{inner: inner, n: n}, nil
}

// Matches returns true for every nth entry that the inner filter accepts.
func (s *SampleFilter) Matches(e parser.Entry) bool {
	if !s.inner.Matches(e) {
		return false
	}
	s.count++
	return s.count%s.n == 0
}

// String returns a human-readable description of the filter.
func (s *SampleFilter) String() string {
	return fmt.Sprintf("sample(every %d, inner=%s)", s.n, s.inner)
}

// Reset resets the internal counter.
func (s *SampleFilter) Reset() {
	s.count = 0
}
