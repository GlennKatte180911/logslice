package filter

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makePrefixEntry(msg string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   msg,
		Fields:    map[string]string{},
	}
}

func TestNewPrefixFilter_EmptyPrefix(t *testing.T) {
	_, err := NewPrefixFilter("", true)
	if err == nil {
		t.Fatal("expected error for empty prefix, got nil")
	}
}

func TestNewPrefixFilter_ValidPrefix(t *testing.T) {
	f, err := NewPrefixFilter("hello", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestPrefixFilter_Matches_CaseSensitive(t *testing.T) {
	f, _ := NewPrefixFilter("ERROR", true)

	tests := []struct {
		msg  string
		want bool
	}{
		{"ERROR: disk full", true},
		{"error: disk full", false},
		{"INFO: all good", false},
		{"ERROR", true},
	}
	for _, tt := range tests {
		got := f.Matches(makePrefixEntry(tt.msg))
		if got != tt.want {
			t.Errorf("Matches(%q) = %v, want %v", tt.msg, got, tt.want)
		}
	}
}

func TestPrefixFilter_Matches_CaseInsensitive(t *testing.T) {
	f, _ := NewPrefixFilter("error", false)

	tests := []struct {
		msg  string
		want bool
	}{
		{"ERROR: disk full", true},
		{"error: disk full", true},
		{"Error: something", true},
		{"INFO: all good", false},
	}
	for _, tt := range tests {
		got := f.Matches(makePrefixEntry(tt.msg))
		if got != tt.want {
			t.Errorf("Matches(%q) = %v, want %v", tt.msg, got, tt.want)
		}
	}
}

func TestPrefixFilter_String(t *testing.T) {
	f, _ := NewPrefixFilter("hello", true)
	s := f.String()
	if s == "" {
		t.Fatal("expected non-empty string representation")
	}
	expected := `prefix("hello", case-sensitive)`
	if s != expected {
		t.Errorf("String() = %q, want %q", s, expected)
	}
}
