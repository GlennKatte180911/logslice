package filter

import (
	"strings"

	"github.com/user/logslice/internal/parser"
)

// Filter is the common interface for all log entry filters.
type Filter interface {
	Matches(parser.Entry) bool
	String() string
}

// Chain applies a sequence of filters using AND semantics: an entry must
// satisfy every filter to be considered a match.
type Chain struct {
	filters []Filter
}

// NewChain returns an empty Chain.
func NewChain() *Chain {
	return &Chain{}
}

// Add appends a filter to the chain.
func (c *Chain) Add(f Filter) {
	c.filters = append(c.filters, f)
}

// Len returns the number of filters currently in the chain.
func (c *Chain) Len() int {
	return len(c.filters)
}

// Matches returns true when the entry satisfies all filters in the chain.
// An empty chain matches every entry.
func (c *Chain) Matches(e parser.Entry) bool {
	for _, f := range c.filters {
		if !f.Matches(e) {
			return false
		}
	}
	return true
}

// String returns a human-readable summary of all chained filters.
func (c *Chain) String() string {
	if len(c.filters) == 0 {
		return "Chain(empty)"
	}
	parts := make([]string, len(c.filters))
	for i, f := range c.filters {
		parts[i] = f.String()
	}
	return "Chain(" + strings.Join(parts, " AND ") + ")"
}
