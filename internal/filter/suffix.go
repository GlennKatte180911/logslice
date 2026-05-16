package filter

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// suffixFilter matches log entries whose message ends with a given suffix.
type suffixFilter struct {
	suffix          string
	caseInsensitive bool
}

// NewSuffixFilter returns a Filter that matches entries whose message ends
// with suffix. If caseInsensitive is true the comparison is folded to lower
// case before matching. Returns an error when suffix is empty.
func NewSuffixFilter(suffix string, caseInsensitive bool) (Filter, error) {
	if suffix == "" {
		return nil, fmt.Errorf("suffix filter: suffix must not be empty")
	}
	return &suffixFilter{
		suffix:          suffix,
		caseInsensitive: caseInsensitive,
	}, nil
}

// Matches returns true when the entry's message ends with the configured suffix.
func (f *suffixFilter) Matches(e parser.Entry) bool {
	msg := e.Message
	suffix := f.suffix
	if f.caseInsensitive {
		msg = strings.ToLower(msg)
		suffix = strings.ToLower(suffix)
	}
	return strings.HasSuffix(msg, suffix)
}

// String returns a human-readable description of the filter.
func (f *suffixFilter) String() string {
	if f.caseInsensitive {
		return fmt.Sprintf("suffix(%q, case-insensitive)", f.suffix)
	}
	return fmt.Sprintf("suffix(%q)", f.suffix)
}
