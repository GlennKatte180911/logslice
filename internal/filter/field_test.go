package filter_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/parser"
)

func makeFieldEntry(fields map[string]string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   "test message",
		Fields:    fields,
	}
}

func TestNewFieldFilter_EmptyKey(t *testing.T) {
	_, err := filter.NewFieldFilter("", "value")
	if err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestNewFieldFilter_EmptyValue(t *testing.T) {
	_, err := filter.NewFieldFilter("service", "")
	if err == nil {
		t.Fatal("expected error for empty value, got nil")
	}
}

func TestFieldFilter_Matches(t *testing.T) {
	f, err := filter.NewFieldFilter("service", "auth")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		name   string
		fields map[string]string
		want   bool
	}{
		{"matching field", map[string]string{"service": "auth"}, true},
		{"wrong value", map[string]string{"service": "billing"}, false},
		{"missing key", map[string]string{"host": "localhost"}, false},
		{"nil fields", nil, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := makeFieldEntry(tc.fields)
			if got := f.Matches(e); got != tc.want {
				t.Errorf("Matches() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestFieldFilter_String(t *testing.T) {
	f, _ := filter.NewFieldFilter("env", "production")
	got := f.String()
	want := "field(env=production)"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
