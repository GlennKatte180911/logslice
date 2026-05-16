package filter

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeSuffixEntry(msg string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   msg,
		Fields:    map[string]string{},
	}
}

func TestNewSuffixFilter_EmptySuffix(t *testing.T) {
	_, err := NewSuffixFilter("", false)
	if err == nil {
		t.Fatal("expected error for empty suffix, got nil")
	}
}

func TestNewSuffixFilter_ValidSuffix(t *testing.T) {
	f, err := NewSuffixFilter("done", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestSuffixFilter_Matches_CaseSensitive(t *testing.T) {
	f, _ := NewSuffixFilter("done", false)

	cases := []struct {
		msg  string
		want bool
	}{
		{"task done", true},
		{"task Done", false},
		{"done task", false},
		{"done", true},
		{"", false},
	}
	for _, tc := range cases {
		got := f.Matches(makeSuffixEntry(tc.msg))
		if got != tc.want {
			t.Errorf("Matches(%q) = %v, want %v", tc.msg, got, tc.want)
		}
	}
}

func TestSuffixFilter_Matches_CaseInsensitive(t *testing.T) {
	f, _ := NewSuffixFilter("DONE", true)

	cases := []struct {
		msg  string
		want bool
	}{
		{"task done", true},
		{"task Done", true},
		{"task DONE", true},
		{"done task", false},
	}
	for _, tc := range cases {
		got := f.Matches(makeSuffixEntry(tc.msg))
		if got != tc.want {
			t.Errorf("Matches(%q) = %v, want %v", tc.msg, got, tc.want)
		}
	}
}

func TestSuffixFilter_String(t *testing.T) {
	f1, _ := NewSuffixFilter("end", false)
	if s := f1.String(); s != `suffix("end")` {
		t.Errorf("unexpected String(): %q", s)
	}

	f2, _ := NewSuffixFilter("end", true)
	if s := f2.String(); s != `suffix("end", case-insensitive)` {
		t.Errorf("unexpected String(): %q", s)
	}
}
