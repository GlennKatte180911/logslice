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
	Levels  string
	Format  string
	Input   string
}

// parseFlags parses os.Args and returns a Config.
// It writes usage to stderr and exits on error.
func parseFlags(args []string) (Config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var cfg Config
	fs.StringVar(&cfg.From, "from", "", "start of time range (RFC3339)")
	fs.StringVar(&cfg.To, "to", "", "end of time range (RFC3339)")
	fs.StringVar(&cfg.Pattern, "pattern", "", "regex pattern to match against log message")
	fs.StringVar(&cfg.Levels, "levels", "", "comma-separated list of log levels to include (e.g. ERROR,WARN)")
	fs.StringVar(&cfg.Format, "format", "text", "output format: text, json, or csv")
	fs.StringVar(&cfg.Input, "input", "", "path to log file (defaults to stdin)")

	if err := fs.Parse(args); err != nil {
		return Config{}, fmt.Errorf("parseFlags: %w", err)
	}

	return cfg, nil
}
