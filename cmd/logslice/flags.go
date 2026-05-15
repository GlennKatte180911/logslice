package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// config holds all parsed CLI flags.
type config struct {
	From    string
	To      string
	Pattern string
	Levels  []string
	Field   string // "key=value" pair for field filter
	Format  string
	Input   string
}

// parseFlags parses os.Args and returns a populated config.
// On error it writes usage to stderr and returns a non-nil error.
func parseFlags(args []string) (*config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var (
		from    = fs.String("from", "", "start of time range (RFC3339)")
		to      = fs.String("to", "", "end of time range (RFC3339)")
		pattern = fs.String("pattern", "", "regex pattern to match against log message")
		levels  = fs.String("levels", "", "comma-separated list of log levels to include")
		field   = fs.String("field", "", "field filter as key=value")
		format  = fs.String("format", "text", "output format: text, json, csv")
		input   = fs.String("input", "-", "input file path; use - for stdin")
	)

	if err := fs.Parse(args); err != nil {
		return nil, fmt.Errorf("parseFlags: %w", err)
	}

	var lvlList []string
	if *levels != "" {
		for _, l := range strings.Split(*levels, ",") {
			if t := strings.TrimSpace(l); t != "" {
				lvlList = append(lvlList, t)
			}
		}
	}

	return &config{
		From:    *from,
		To:      *to,
		Pattern: *pattern,
		Levels:  lvlList,
		Field:   *field,
		Format:  *format,
		Input:   *input,
	}, nil
}
