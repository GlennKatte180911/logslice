package filter

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// PrefixFilter matches log entries whose message starts with a given prefix.
type PrefixFilter struct {
	prefix          string
	caseSensitive   bool
	normalizedPrefix string
}

// NewPrefixFilter creates a PrefixFilter that matches entries whose message
// starts with prefix. If caseSensitive is false the comparison is done in
// lower-case.
func NewPrefixFilter(prefix string, caseSensitive bool) (*PrefixFilter, error) {
	if prefix == "" {
		return nil, fmt.Errorf("prefix filter: prefix must not be empty")
	}
	normalized := prefix
	if !caseSensitive {
		normalized = strings.ToLower(prefix)
	}
	return &PrefixFilter{
		prefix:           prefix,
		caseSensitive:    caseSensitive,
		normalizedPrefix: normalized,
	}, nil
}

// Matches returns true when the entry message starts with the configured prefix.
func (f *PrefixFilter) Matches(e parser.Entry) bool {
	msg := e.Message
	if !f.caseSensitive {
		msg = strings.ToLower(msg)
	}
	return strings.HasPrefix(msg, f.normalizedPrefix)
}

// String returns a human-readable description of the filter.
func (f *PrefixFilter) String() string {
	cs := "case-insensitive"
	if f.caseSensitive {
		cs = "case-sensitive"
	}
	return fmt.Sprintf("prefix(%q, %s)", f.prefix, cs)
}
