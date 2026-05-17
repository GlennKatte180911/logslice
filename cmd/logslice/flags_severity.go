package main

import (
	"fmt"

	"github.com/user/logslice/internal/filter"
)

// buildSeverityFilter constructs a SeverityFilter from the --min-level flag
// value. It returns nil, nil when the flag was not provided.
func buildSeverityFilter(minLevel string) (*filter.SeverityFilter, error) {
	if minLevel == "" {
		return nil, nil
	}
	f, err := filter.NewSeverityFilter(minLevel)
	if err != nil {
		return nil, fmt.Errorf("--min-level: %w", err)
	}
	return f, nil
}
