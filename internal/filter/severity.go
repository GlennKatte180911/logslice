package filter

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// severityRank maps level strings to a numeric rank for comparison.
var severityRank = map[string]int{
	"trace": 0,
	"debug": 1,
	"info":  2,
	"warn":  3,
	"error": 4,
	"fatal": 5,
}

// SeverityFilter passes entries whose level meets a minimum severity threshold.
type SeverityFilter struct {
	minLevel string
	minRank  int
}

// NewSeverityFilter creates a SeverityFilter that passes entries at or above
// the given minimum level (e.g. "warn" passes warn, error, and fatal).
func NewSeverityFilter(minLevel string) (*SeverityFilter, error) {
	norm := strings.ToLower(strings.TrimSpace(minLevel))
	if norm == "" {
		return nil, fmt.Errorf("severity: minimum level must not be empty")
	}
	rank, ok := severityRank[norm]
	if !ok {
		return nil, fmt.Errorf("severity: unknown level %q; valid levels: trace, debug, info, warn, error, fatal", minLevel)
	}
	return &SeverityFilter{minLevel: norm, minRank: rank}, nil
}

// Matches returns true when the entry's level is at or above the minimum severity.
func (f *SeverityFilter) Matches(e parser.Entry) bool {
	norm := strings.ToLower(strings.TrimSpace(e.Level))
	rank, ok := severityRank[norm]
	if !ok {
		// Unknown levels do not satisfy any threshold.
		return false
	}
	return rank >= f.minRank
}

// String returns a human-readable description of the filter.
func (f *SeverityFilter) String() string {
	return fmt.Sprintf("severity(>=%s)", f.minLevel)
}
