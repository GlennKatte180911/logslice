package filter

import (
	"fmt"
	"sync/atomic"

	"github.com/user/logslice/internal/parser"
)

// CountFilter wraps another filter and counts how many entries matched.
type CountFilter struct {
	inner   Matcher
	matched atomic.Int64
	total   atomic.Int64
}

// Matcher is the common interface implemented by all filters.
type Matcher interface {
	Matches(e parser.Entry) bool
	String() string
}

// NewCountFilter wraps inner so that match statistics are tracked.
// inner must not be nil.
func NewCountFilter(inner Matcher) (*CountFilter, error) {
	if inner == nil {
		return nil, fmt.Errorf("count filter: inner matcher must not be nil")
	}
	return &CountFilter{inner: inner}, nil
}

// Matches delegates to the inner filter and increments counters.
func (c *CountFilter) Matches(e parser.Entry) bool {
	c.total.Add(1)
	ok := c.inner.Matches(e)
	if ok {
		c.matched.Add(1)
	}
	return ok
}

// Matched returns the number of entries that passed the inner filter.
func (c *CountFilter) Matched() int64 { return c.matched.Load() }

// Total returns the total number of entries evaluated.
func (c *CountFilter) Total() int64 { return c.total.Load() }

// String returns a human-readable description.
func (c *CountFilter) String() string {
	return fmt.Sprintf("count(%s)", c.inner)
}
