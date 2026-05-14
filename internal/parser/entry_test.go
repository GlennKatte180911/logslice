package parser

import (
	"testing"
	"time"
)

func mustTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestEntryString(t *testing.T) {
	e := Entry{
		Timestamp: mustTime("2024-01-15T10:00:00Z"),
		Level:     LevelInfo,
		Message:   "server started",
	}
	got := e.String()
	want := "2024-01-15T10:00:00Z [INFO] server started"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestFilterMatches(t *testing.T) {
	ts := mustTime("2024-01-15T12:00:00Z")
	from := mustTime("2024-01-15T10:00:00Z")
	to := mustTime("2024-01-15T14:00:00Z")
	before := mustTime("2024-01-15T08:00:00Z")
	after := mustTime("2024-01-15T16:00:00Z")

	tests := []struct {
		name   string
		filter Filter
		entry  Entry
		want   bool
	}{
		{"no filter", Filter{}, Entry{Timestamp: ts, Level: LevelInfo}, true},
		{"within range", Filter{From: &from, To: &to}, Entry{Timestamp: ts, Level: LevelInfo}, true},
		{"before from", Filter{From: &from}, Entry{Timestamp: before, Level: LevelInfo}, false},
		{"after to", Filter{To: &to}, Entry{Timestamp: after, Level: LevelInfo}, false},
		{"level match", Filter{Level: LevelError}, Entry{Timestamp: ts, Level: LevelError}, true},
		{"level mismatch", Filter{Level: LevelError}, Entry{Timestamp: ts, Level: LevelInfo}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.filter.Matches(tt.entry); got != tt.want {
				t.Errorf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}
