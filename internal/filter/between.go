package filter

import (
	"fmt"

	"github.com/user/logslice/internal/parser"
)

// BetweenFilter passes entries whose message or field value falls
// lexicographically between two boundary strings (inclusive).
type BetweenFilter struct {
	key   string // empty string means match against Message
	low   string
	high  string
}

// NewBetweenFilter creates a BetweenFilter for the given key (or message when
// key is empty) and [low, high] bounds. Returns an error when low > high.
func NewBetweenFilter(key, low, high string) (*BetweenFilter, error) {
	if low == "" {
		return nil, fmt.Errorf("between filter: low bound must not be empty")
	}
	if high == "" {
		return nil, fmt.Errorf("between filter: high bound must not be empty")
	}
	if low > high {
		return nil, fmt.Errorf("between filter: low %q must not be greater than high %q", low, high)
	}
	return &BetweenFilter{key: key, low: low, high: high}, nil
}

// Matches returns true when the target value is within [low, high].
func (f *BetweenFilter) Matches(e parser.Entry) bool {
	var val string
	if f.key == "" {
		val = e.Message
	} else {
		v, ok := e.Fields[f.key]
		if !ok {
			return false
		}
		val = v
	}
	return val >= f.low && val <= f.high
}

// String returns a human-readable description of the filter.
func (f *BetweenFilter) String() string {
	if f.key == "" {
		return fmt.Sprintf("between(message, %q, %q)", f.low, f.high)
	}
	return fmt.Sprintf("between(%s, %q, %q)", f.key, f.low, f.high)
}
