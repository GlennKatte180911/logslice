package filter

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeRegexFieldEntry(key, value string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   "test message",
		Fields:    map[string]string{key: value},
	}
}

func TestNewRegexFieldFilter_EmptyKey(t *testing.T) {
	_, err := NewRegexFieldFilter("", "foo.*")
	if err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestNewRegexFieldFilter_EmptyPattern(t *testing.T) {
	_, err := NewRegexFieldFilter("service", "")
	if err == nil {
		t.Fatal("expected error for empty pattern, got nil")
	}
}

func TestNewRegexFieldFilter_InvalidRegex(t *testing.T) {
	_, err := NewRegexFieldFilter("service", "[invalid")
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestRegexFieldFilter_Matches(t *testing.T) {
	f, err := NewRegexFieldFilter("service", "^auth.*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		name    string
		entry   parser.Entry
		want    bool
	}{
		{"matches prefix", makeRegexFieldEntry("service", "auth-service"), true},
		{"no match", makeRegexFieldEntry("service", "billing-service"), false},
		{"missing field", makeRegexFieldEntry("host", "auth-service"), false},
		{"empty value", makeRegexFieldEntry("service", ""), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := f.Matches(tc.entry)
			if got != tc.want {
				t.Errorf("Matches() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestRegexFieldFilter_String(t *testing.T) {
	f, err := NewRegexFieldFilter("service", "^auth.*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := f.String()
	want := "regex_field(service=~/^auth.*/)"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
