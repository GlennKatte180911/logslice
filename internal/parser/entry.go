package parser

import (
	"fmt"
	"time"
)

// LogLevel represents the severity level of a log entry.
type LogLevel string

const (
	LevelDebug LogLevel = "DEBUG"
	LevelInfo  LogLevel = "INFO"
	LevelWarn  LogLevel = "WARN"
	LevelError LogLevel = "ERROR"
	LevelFatal LogLevel = "FATAL"
)

// Entry represents a single parsed log entry.
type Entry struct {
	Timestamp time.Time
	Level     LogLevel
	Message   string
	Raw       string
}

// String returns a formatted representation of the log entry.
func (e Entry) String() string {
	return fmt.Sprintf("%s [%s] %s", e.Timestamp.Format(time.RFC3339), e.Level, e.Message)
}

// Filter holds criteria for filtering log entries.
type Filter struct {
	From    *time.Time
	To      *time.Time
	Level   LogLevel
	Pattern string
}

// Matches returns true if the entry satisfies all non-zero filter criteria.
func (f Filter) Matches(e Entry) bool {
	if f.From != nil && e.Timestamp.Before(*f.From) {
		return false
	}
	if f.To != nil && e.Timestamp.After(*f.To) {
		return false
	}
	if f.Level != "" && e.Level != f.Level {
		return false
	}
	return true
}
