package main

import (
	"fmt"
	"os"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/parser"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	cfg, err := parseFlags(args)
	if err != nil {
		return err
	}

	chain, err := buildChain(cfg)
	if err != nil {
		return err
	}

	entries, err := parser.Parse(cfg.input)
	if err != nil {
		return fmt.Errorf("parsing input: %w", err)
	}

	for _, e := range entries {
		if chain.Matches(e) {
			fmt.Println(e)
		}
	}
	return nil
}

func buildChain(cfg *config) (*filter.Chain, error) {
	var filters []filter.Filter

	if cfg.from != "" || cfg.to != "" {
		tr, err := filter.NewTimeRange(cfg.from, cfg.to)
		if err != nil {
			return nil, fmt.Errorf("time range filter: %w", err)
		}
		filters = append(filters, tr)
	}

	if cfg.pattern != "" {
		pf, err := filter.NewPatternFilter(cfg.pattern)
		if err != nil {
			return nil, fmt.Errorf("pattern filter: %w", err)
		}
		filters = append(filters, pf)
	}

	return filter.NewChain(filters...), nil
}
