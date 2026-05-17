package filter_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/parser"
)

func makeSeverityEntry(level string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   "test message",
		Fields:    map[string]string{},
	}
}

func TestNewSeverityFilter_EmptyLevel(t *testing.T) {
	_, err := filter.NewSeverityFilter("")
	if err == nil {
		t.Fatal("expected error for empty level, got nil")
	}
}

func TestNewSeverityFilter_UnknownLevel(t *testing.T) {
	_, err := filter.NewSeverityFilter("verbose")
	if err == nil {
		t.Fatal("expected error for unknown level, got nil")
	}
}

func TestNewSeverityFilter_ValidLevel(t *testing.T) {
	f, err := filter.NewSeverityFilter("warn")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestSeverityFilter_Matches(t *testing.T) {
	tests := []struct {
		minLevel   string
		entryLevel string
		want       bool
	}{
		{"info", "trace", false},
		{"info", "debug", false},
		{"info", "info", true},
		{"info", "warn", true},
		{"info", "error", true},
		{"info", "fatal", true},
		{"warn", "info", false},
		{"warn", "warn", true},
		{"warn", "error", true},
		{"error", "warn", false},
		{"error", "error", true},
		{"debug", "trace", false},
		{"debug", "debug", true},
		{"info", "WARN", true}, // case-insensitive
		{"info", "unknown", false},
	}

	for _, tc := range tests {
		f, err := filter.NewSeverityFilter(tc.minLevel)
		if err != nil {
			t.Fatalf("NewSeverityFilter(%q): %v", tc.minLevel, err)
		}
		got := f.Matches(makeSeverityEntry(tc.entryLevel))
		if got != tc.want {
			t.Errorf("min=%q entry=%q: got %v, want %v", tc.minLevel, tc.entryLevel, got, tc.want)
		}
	}
}

func TestSeverityFilter_String(t *testing.T) {
	f, _ := filter.NewSeverityFilter("error")
	got := f.String()
	expected := "severity(>=error)"
	if got != expected {
		t.Errorf("String() = %q, want %q", got, expected)
	}
}

func TestSeverityFilter_AllValidLevels(t *testing.T) {
	validLevels := []string{"trace", "debug", "info", "warn", "error", "fatal"}
	for _, level := range validLevels {
		f, err := filter.NewSeverityFilter(level)
		if err != nil {
			t.Errorf("NewSeverityFilter(%q) returned unexpected error: %v", level, err)
			continue
		}
		if f == nil {
			t.Errorf("NewSeverityFilter(%q) returned nil filter", level)
		}
	}
}
