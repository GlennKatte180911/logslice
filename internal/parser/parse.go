package parser

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"time"
)

// defaultPattern matches lines like: 2024-01-15T10:00:00Z [INFO] message text
var defaultPattern = regexp.MustCompile(
	`^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:Z|[+-]\d{2}:\d{2}))\s+\[(DEBUG|INFO|WARN|ERROR|FATAL)\]\s+(.+)$`,
)

// ParseResult holds successfully parsed entries and any lines that failed parsing.
type ParseResult struct {
	Entries  []Entry
	Skipped  int
}

// Parse reads log lines from r, parses each one, and returns all entries
// that match the given filter. Lines that do not match the expected format
// are counted as skipped.
func Parse(r io.Reader, f Filter) (ParseResult, error) {
	var result ParseResult
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		entry, err := parseLine(line)
		if err != nil {
			result.Skipped++
			continue
		}

		if f.Matches(entry) {
			result.Entries = append(result.Entries, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return result, fmt.Errorf("scanner error: %w", err)
	}

	return result, nil
}

// parseLine attempts to parse a single log line into an Entry.
func parseLine(line string) (Entry, error) {
	matches := defaultPattern.FindStringSubmatch(line)
	if matches == nil {
		return Entry{}, fmt.Errorf("line does not match log format: %q", line)
	}

	ts, err := time.Parse(time.RFC3339, matches[1])
	if err != nil {
		return Entry{}, fmt.Errorf("invalid timestamp %q: %w", matches[1], err)
	}

	return Entry{
		Timestamp: ts,
		Level:     LogLevel(matches[2]),
		Message:   matches[3],
		Raw:       line,
	}, nil
}
