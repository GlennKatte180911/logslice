package filter

import (
	"fmt"

	"github.com/user/logslice/internal/parser"
)

// NotFilter wraps another filter and inverts its result.
type NotFilter struct {
	inner Filter
}

// Filter is the interface satisfied by all filters in this package.
// It is declared here to avoid circular imports if not already declared.
type Filter interface {
	Matches(e parser.Entry) bool
	String() string
}

// NewNotFilter returns a NotFilter that negates the given inner filter.
// It returns an error if inner is nil.
func NewNotFilter(inner Filter) (*NotFilter, error) {
	if inner == nil {
		return nil, fmt.Errorf("not filter: inner filter must not be nil")
	}
	return &NotFilter{inner: inner}, nil
}

// Matches returns true when the inner filter does NOT match the entry.
func (n *NotFilter) Matches(e parser.Entry) bool {
	return !n.inner.Matches(e)
}

// String returns a human-readable description of the filter.
func (n *NotFilter) String() string {
	return fmt.Sprintf("NOT(%s)", n.inner)
}
