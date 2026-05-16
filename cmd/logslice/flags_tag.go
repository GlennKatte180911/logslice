package main

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/filter"
)

// buildTagFilter creates a TagFilter from a comma-separated tag string.
// It returns (nil, nil) when the flag value is empty, signalling that no
// tag filter should be added to the chain.
func buildTagFilter(raw string) (*filter.TagFilter, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}

	parts := strings.Split(raw, ",")
	tags := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			tags = append(tags, p)
		}
	}

	if len(tags) == 0 {
		return nil, fmt.Errorf("--tags: no valid tags found in %q", raw)
	}

	f, err := filter.NewTagFilter(tags)
	if err != nil {
		return nil, fmt.Errorf("--tags: %w", err)
	}
	return f, nil
}
