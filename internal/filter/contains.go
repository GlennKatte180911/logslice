package filter

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// containsFilter matches entries whose message contains a given substring.
type containsFilter struct {
	substring string
	caseSensitive bool
}

// NewContainsFilter returns a Filter that matches log entries whose message
// contains the given substring. If caseSensitive is false the comparison is
// done in lower-case.
func NewContainsFilter(substring string, caseSensitive bool) (*containsFilter, error) {
	if substring == "" {
		return nil, fmt.Errorf("contains filter: substring must not be empty")
	}
	return &containsFilter{substring: substring, caseSensitive: caseSensitive}, nil
}

// Matches returns true when the entry's message contains the configured
// substring.
func (f *containsFilter) Matches(e parser.Entry) bool {
	msg := e.Message
	sub := f.substring
	if !f.caseSensitive {
		msg = strings.ToLower(msg)
		sub = strings.ToLower(sub)
	}
	return strings.Contains(msg, sub)
}

// String returns a human-readable description of the filter.
func (f *containsFilter) String() string {
	sensitivity := "case-insensitive"
	if f.caseSensitive {
		sensitivity = "case-sensitive"
	}
	return fmt.Sprintf("contains(%q, %s)", f.substring, sensitivity)
}
