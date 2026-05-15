package filter

import (
	"github.com/user/logslice/internal/parser"
)

// Filter is the common interface implemented by all log entry filters.
type Filter interface {
	Matches(e parser.Entry) bool
}

// Chain combines multiple filters with AND semantics: an entry must satisfy
// every filter in the chain to be considered a match.
type Chain struct {
	filters []Filter
}

// NewChain returns a Chain that applies each provided filter in order.
func NewChain(filters ...Filter) *Chain {
	return &Chain{filters: filters}
}

// Add appends a filter to the chain.
func (c *Chain) Add(f Filter) {
	c.filters = append(c.filters, f)
}

// Matches returns true only if every filter in the chain matches the entry.
// An empty chain matches all entries.
func (c *Chain) Matches(e parser.Entry) bool {
	for _, f := range c.filters {
		if !f.Matches(e) {
			return false
		}
	}
	return true
}

// Apply filters the provided slice of entries, returning only those that
// satisfy the chain.
func (c *Chain) Apply(entries []parser.Entry) []parser.Entry {
	result := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		if c.Matches(e) {
			result = append(result, e)
		}
	}
	return result
}
