package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"time"
)

// config holds all parsed CLI flags.
type config struct {
	from    string
	to      string
	before  string
	after   string
	pattern string
	level   string
	field   string
	format  string
	showStats bool
}

func parseFlags(args []string, stderr io.Writer) (*config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)
	fs.SetOutput(stderr)

	cfg := &config{}

	fs.StringVar(&cfg.from, "from", "", "include entries at or after this time (RFC3339)")
	fs.StringVar(&cfg.to, "to", "", "include entries at or before this time (RFC3339)")
	fs.StringVar(&cfg.before, "before", "", "include entries strictly before this time (RFC3339)")
	fs.StringVar(&cfg.after, "after", "", "include entries strictly after this time (RFC3339)")
	fs.StringVar(&cfg.pattern, "pattern", "", "filter entries whose message matches this regex")
	fs.StringVar(&cfg.level, "level", "", "comma-separated list of log levels to include")
	fs.StringVar(&cfg.field, "field", "", "key=value pair to match in entry fields")
	fs.StringVar(&cfg.format, "format", "text", "output format: text, json, csv")
	fs.BoolVar(&cfg.showStats, "stats", false, "print match statistics to stderr after processing")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if cfg.from != "" && cfg.after != "" {
		return nil, errors.New("flags -from and -after are mutually exclusive")
	}
	if cfg.to != "" && cfg.before != "" {
		return nil, errors.New("flags -to and -before are mutually exclusive")
	}

	for _, pair := range []struct{ name, val string }{
		{"from", cfg.from}, {"to", cfg.to},
		{"before", cfg.before}, {"after", cfg.after},
	} {
		if pair.val == "" {
			continue
		}
		if _, err := time.Parse(time.RFC3339, pair.val); err != nil {
			return nil, fmt.Errorf("flag -%s: invalid RFC3339 time %q", pair.name, pair.val)
		}
	}

	return cfg, nil
}
