package filter

import (
	"fmt"
	"regexp"

	"github.com/user/logslice/internal/parser"
)

// PatternFilter filters log entries by matching a regular expression
// against the entry's message field.
type PatternFilter struct {
	re *regexp.Regexp
}

// NewPatternFilter compiles the given pattern and returns a PatternFilter.
// Returns an error if the pattern is not a valid regular expression.
func NewPatternFilter(pattern string) (*PatternFilter, error) {
	if pattern == "" {
		return nil, fmt.Errorf("pattern must not be empty")
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern %q: %w", pattern, err)
	}
	return &PatternFilter{re: re}, nil
}

// Matches reports whether the entry's message matches the compiled pattern.
func (p *PatternFilter) Matches(e parser.Entry) bool {
	return p.re.MatchString(e.Message)
}

// String returns a human-readable description of the filter.
func (p *PatternFilter) String() string {
	return fmt.Sprintf("PatternFilter(%s)", p.re.String())
}
