package main

import (
	"flag"
	"fmt"
	"os"
)

// Config holds all parsed command-line options.
type Config struct {
	From    string
	To      string
	Pattern string
	Level   string
	Field   string
	Not     bool
	Sample  int
	Dedup   bool
	Limit   int
	Format  string
	Stats   bool
	Args    []string
}

// parseFlags parses os.Args and returns a populated Config.
func parseFlags() (Config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var cfg Config
	fs.StringVar(&cfg.From, "from", "", "start of time range (RFC3339)")
	fs.StringVar(&cfg.To, "to", "", "end of time range (RFC3339)")
	fs.StringVar(&cfg.Pattern, "pattern", "", "regex pattern to match against message")
	fs.StringVar(&cfg.Level, "level", "", "comma-separated list of log levels to include")
	fs.StringVar(&cfg.Field, "field", "", "key=value field filter")
	fs.BoolVar(&cfg.Not, "not", false, "invert the last filter")
	fs.IntVar(&cfg.Sample, "sample", 0, "emit every Nth matching entry (0 = disabled)")
	fs.BoolVar(&cfg.Dedup, "dedup", false, "suppress consecutive duplicate messages")
	fs.IntVar(&cfg.Limit, "limit", 0, "maximum number of entries to output (0 = unlimited)")
	fs.StringVar(&cfg.Format, "format", "text", "output format: text, json, csv")
	fs.BoolVar(&cfg.Stats, "stats", false, "print match statistics to stderr after processing")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return Config{}, fmt.Errorf("parseFlags: %w", err)
	}
	cfg.Args = fs.Args()
	return cfg, nil
}
