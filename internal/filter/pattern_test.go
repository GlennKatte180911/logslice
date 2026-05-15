package filter

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makePatternEntry(msg string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   msg,
	}
}

func TestNewPatternFilter_EmptyPattern(t *testing.T) {
	_, err := NewPatternFilter("")
	if err == nil {
		t.Fatal("expected error for empty pattern, got nil")
	}
}

func TestNewPatternFilter_InvalidRegex(t *testing.T) {
	_, err := NewPatternFilter("[invalid")
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestPatternFilter_Matches(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		msg     string
		want    bool
	}{
		{"exact match", "error", "an error occurred", true},
		{"case sensitive no match", "ERROR", "an error occurred", false},
		{"case insensitive flag", "(?i)error", "an ERROR occurred", true},
		{"no match", "timeout", "connection refused", false},
		{"anchored match", "^started", "started service", true},
		{"anchored no match", "^started", "service started", false},
		{"digit pattern", `\d+`, "request took 42ms", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := NewPatternFilter(tc.pattern)
			if err != nil {
				t.Fatalf("NewPatternFilter(%q) error: %v", tc.pattern, err)
			}
			got := f.Matches(makePatternEntry(tc.msg))
			if got != tc.want {
				t.Errorf("Matches(%q) = %v, want %v", tc.msg, got, tc.want)
			}
		})
	}
}

func TestPatternFilter_String(t *testing.T) {
	f, err := NewPatternFilter(`\d+`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := f.String()
	want := `PatternFilter(\d+)`
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
