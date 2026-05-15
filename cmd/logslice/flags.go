package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type config struct {
	from    string
	to      string
	pattern string
	input   io.Reader
}

func parseFlags(args []string) (*config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)

	var (
		from    = fs.String("from", "", "start of time range (RFC3339, e.g. 2024-01-01T00:00:00Z)")
		to      = fs.String("to", "", "end of time range (RFC3339, e.g. 2024-01-02T00:00:00Z)")
		pattern = fs.String("pattern", "", "regex pattern to match against log message")
		file    = fs.String("file", "", "path to log file (defaults to stdin)")
	)

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	cfg := &config{
		from:    strings.TrimSpace(*from),
		to:      strings.TrimSpace(*to),
		pattern: strings.TrimSpace(*pattern),
		input:   os.Stdin,
	}

	if *file != "" {
		f, err := os.Open(*file)
		if err != nil {
			return nil, fmt.Errorf("opening file %q: %w", *file, err)
		}
		cfg.input = f
	}

	return cfg, nil
}
