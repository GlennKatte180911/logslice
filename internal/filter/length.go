package filter

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// LengthFilter matches log entries whose message length satisfies a comparison.
type LengthFilter struct {
	op  string
	val int
}

// NewLengthFilter returns a LengthFilter that applies op ("lt", "lte", "gt", "gte", "eq")
// against the byte-length of each entry's message.
func NewLengthFilter(op string, val int) (*LengthFilter, error) {
	op = strings.ToLower(strings.TrimSpace(op))
	switch op {
	case "lt", "lte", "gt", "gte", "eq":
	default:
		return nil, fmt.Errorf("length filter: unknown op %q (want lt|lte|gt|gte|eq)", op)
	}
	if val < 0 {
		return nil, fmt.Errorf("length filter: val must be >= 0, got %d", val)
	}
	return &LengthFilter{op: op, val: val}, nil
}

// Matches returns true when the entry's message length satisfies the configured comparison.
func (f *LengthFilter) Matches(e parser.Entry) bool {
	l := len(e.Message)
	switch f.op {
	case "lt":
		return l < f.val
	case "lte":
		return l <= f.val
	case "gt":
		return l > f.val
	case "gte":
		return l >= f.val
	case "eq":
		return l == f.val
	}
	return false
}

// String returns a human-readable description of the filter.
func (f *LengthFilter) String() string {
	return fmt.Sprintf("LengthFilter(message len %s %d)", f.op, f.val)
}
