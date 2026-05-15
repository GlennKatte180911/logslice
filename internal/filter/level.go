package filter

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// LevelFilter matches log entries whose level is in the allowed set.
type LevelFilter struct {
	levels map[string]struct{}
}

// NewLevelFilter creates a LevelFilter from a comma-separated list of levels
// (e.g. "ERROR,WARN"). Returns an error if the list is empty or any level is
// blank after trimming.
func NewLevelFilter(raw string) (*LevelFilter, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, fmt.Errorf("level filter: level list must not be empty")
	}

	parts := strings.Split(raw, ",")
	levels := make(map[string]struct{}, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(strings.ToUpper(p))
		if trimmed == "" {
			return nil, fmt.Errorf("level filter: blank level in list %q", raw)
		}
		levels[trimmed] = struct{}{}
	}

	return &LevelFilter{levels: levels}, nil
}

// Matches returns true when the entry's level is in the allowed set.
func (f *LevelFilter) Matches(e parser.Entry) bool {
	_, ok := f.levels[strings.ToUpper(e.Level)]
	return ok
}

// String returns a human-readable description of the filter.
func (f *LevelFilter) String() string {
	keys := make([]string, 0, len(f.levels))
	for k := range f.levels {
		keys = append(keys, k)
	}
	return fmt.Sprintf("level in [%s]", strings.Join(keys, ","))
}
