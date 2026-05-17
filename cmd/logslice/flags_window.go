package main

import (
	"fmt"
	"time"

	"github.com/user/logslice/internal/filter"
)

// buildWindowFilter constructs a WindowFilter from the parsed CLI flags.
// windowAnchor is an RFC3339 timestamp string; windowBefore and windowAfter
// are duration strings (e.g. "5m", "1h30m").
func buildWindowFilter(cfg *config) (*filter.WindowFilter, error) {
	if cfg.windowAnchor == "" {
		return nil, nil
	}

	anchor, err := time.Parse(time.RFC3339, cfg.windowAnchor)
	if err != nil {
		return nil, fmt.Errorf("--window-anchor: %w", err)
	}

	before, err := time.ParseDuration(cfg.windowBefore)
	if err != nil {
		return nil, fmt.Errorf("--window-before: %w", err)
	}

	after, err := time.ParseDuration(cfg.windowAfter)
	if err != nil {
		return nil, fmt.Errorf("--window-after: %w", err)
	}

	wf, err := filter.NewWindowFilter(anchor, before, after)
	if err != nil {
		return nil, fmt.Errorf("window filter: %w", err)
	}
	return wf, nil
}
