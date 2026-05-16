package filter

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeFieldExistsEntry(fields map[string]string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   "test message",
		Fields:    fields,
	}
}

func TestNewFieldExistsFilter_EmptyKey(t *testing.T) {
	_, err := NewFieldExistsFilter("", false)
	if err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestNewFieldExistsFilter_Valid(t *testing.T) {
	f, err := NewFieldExistsFilter("request_id", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestFieldExistsFilter_Matches(t *testing.T) {
	tests := []struct {
		name   string
		key    string
		negate bool
		fields map[string]string
		want   bool
	}{
		{"key present, not negated", "request_id", false, map[string]string{"request_id": "abc"}, true},
		{"key absent, not negated", "request_id", false, map[string]string{"other": "val"}, false},
		{"key present, negated", "request_id", true, map[string]string{"request_id": "abc"}, false},
		{"key absent, negated", "request_id", true, map[string]string{"other": "val"}, true},
		{"nil fields, not negated", "request_id", false, nil, false},
		{"nil fields, negated", "request_id", true, nil, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := NewFieldExistsFilter(tc.key, tc.negate)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			entry := makeFieldExistsEntry(tc.fields)
			if got := f.Matches(entry); got != tc.want {
				t.Errorf("Matches() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestFieldExistsFilter_String(t *testing.T) {
	f, _ := NewFieldExistsFilter("user_id", false)
	if s := f.String(); s != `field_exists("user_id")` {
		t.Errorf("unexpected String(): %q", s)
	}

	fn, _ := NewFieldExistsFilter("user_id", true)
	if s := fn.String(); s != `field_missing("user_id")` {
		t.Errorf("unexpected String() for negated: %q", s)
	}
}
