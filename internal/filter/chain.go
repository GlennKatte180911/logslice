package filter

import (
	"strings"

	"github.com/user/logslice/internal/parser"
)

// Filter is the interface implemented by all log entry filters.
type Filter interface {
	Matches(parser.Entry) bool
	String() string
}

// Chain combines multiple filters with AND semantics: an entry must satisfy
// every filter in the chain to be considered a match.
type Chain struct {
	filters []Filter
}

// NewChain creates a Chain from the provided filters. Nil filters are ignored.
func NewChain(filters ...Filter) *Chain {
	var valid []Filter
	for _, f := range filters {
		if f != nil {
			valid = append(valid, f)
		}
	}
	return &Chain{filters: valid}
}

// Add appends a filter to the chain.
func (c *Chain) Add(f Filter) {
	if f != nil {
		c.filters = append(c.filters, f)
	}
}

// Matches returns true only when every filter in the chain matches the entry.
// An empty chain matches all entries.
func (c *Chain) Matches(e parser.Entry) bool {
	for _, f := range c.filters {
		if !f.Matches(e) {
			return false
		}
	}
	return true
}

// String returns a comma-separated list of filter descriptions.
func (c *Chain) String() string {
	parts := make([]string, len(c.filters))
	for i, f := range c.filters {
		parts[i] = f.String()
	}
	return "chain[" + strings.Join(parts, ", ") + "]"
}
