package filter

import (
	"fmt"
	"regexp"

	"github.com/user/logslice/internal/parser"
)

// RegexFieldFilter matches log entries where a specific field value matches a regular expression.
type RegexFieldFilter struct {
	key     string
	pattern *regexp.Regexp
}

// NewRegexFieldFilter creates a new RegexFieldFilter for the given field key and regex pattern.
// Returns an error if key or pattern is empty, or if the pattern is not valid regex.
func NewRegexFieldFilter(key, pattern string) (*RegexFieldFilter, error) {
	if key == "" {
		return nil, fmt.Errorf("field key must not be empty")
	}
	if pattern == "" {
		return nil, fmt.Errorf("regex pattern must not be empty")
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern %q: %w", pattern, err)
	}
	return &RegexFieldFilter{key: key, pattern: re}, nil
}

// Matches returns true if the entry's field identified by key exists and its
// value matches the compiled regular expression.
func (f *RegexFieldFilter) Matches(e parser.Entry) bool {
	val, ok := e.Fields[f.key]
	if !ok {
		return false
	}
	return f.pattern.MatchString(val)
}

// String returns a human-readable description of the filter.
func (f *RegexFieldFilter) String() string {
	return fmt.Sprintf("regex_field(%s=~/%s/)", f.key, f.pattern.String())
}
