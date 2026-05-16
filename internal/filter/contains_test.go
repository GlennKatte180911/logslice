package filter

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeContainsEntry(msg string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   msg,
		Fields:    map[string]string{},
	}
}

func TestNewContainsFilter_EmptySubstring(t *testing.T) {
	_, err := NewContainsFilter("", true)
	if err == nil {
		t.Fatal("expected error for empty substring, got nil")
	}
}

func TestNewContainsFilter_ValidSubstring(t *testing.T) {
	f, err := NewContainsFilter("hello", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestContainsFilter_Matches_CaseSensitive(t *testing.T) {
	f, _ := NewContainsFilter("Error", true)

	tests := []struct {
		msg  string
		want bool
	}{
		{"An Error occurred", true},
		{"an error occurred", false},
		{"no match here", false},
		{"Error", true},
	}
	for _, tc := range tests {
		got := f.Matches(makeContainsEntry(tc.msg))
		if got != tc.want {
			t.Errorf("Matches(%q) = %v, want %v", tc.msg, got, tc.want)
		}
	}
}

func TestContainsFilter_Matches_CaseInsensitive(t *testing.T) {
	f, _ := NewContainsFilter("error", false)

	tests := []struct {
		msg  string
		want bool
	}{
		{"An Error occurred", true},
		{"an error occurred", true},
		{"no match here", false},
		{"ERROR: disk full", true},
	}
	for _, tc := range tests {
		got := f.Matches(makeContainsEntry(tc.msg))
		if got != tc.want {
			t.Errorf("Matches(%q) = %v, want %v", tc.msg, got, tc.want)
		}
	}
}

func TestContainsFilter_String(t *testing.T) {
	f, _ := NewContainsFilter("hello", true)
	got := f.String()
	if got != `contains("hello", case-sensitive)` {
		t.Errorf("unexpected String(): %q", got)
	}

	f2, _ := NewContainsFilter("world", false)
	got2 := f2.String()
	if got2 != `contains("world", case-insensitive)` {
		t.Errorf("unexpected String(): %q", got2)
	}
}
