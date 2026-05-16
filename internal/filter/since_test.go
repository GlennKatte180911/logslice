package filter

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeSinceEntry(ts time.Time) parser.Entry {
	return parser.Entry{
		Timestamp: ts,
		Level:     "info",
		Message:   "since test",
		Fields:    map[string]string{},
	}
}

func TestNewSinceFilter_ZeroDuration(t *testing.T) {
	_, err := NewSinceFilter(0)
	if err == nil {
		t.Fatal("expected error for zero duration, got nil")
	}
}

func TestNewSinceFilter_NegativeDuration(t *testing.T) {
	_, err := NewSinceFilter(-5 * time.Minute)
	if err == nil {
		t.Fatal("expected error for negative duration, got nil")
	}
}

func TestNewSinceFilter_ValidDuration(t *testing.T) {
	f, err := NewSinceFilter(10 * time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestSinceFilter_Matches(t *testing.T) {
	now := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	f, err := NewSinceFilter(10 * time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	f.now = func() time.Time { return now }

	tests := []struct {
		name    string
		ts      time.Time
		want    bool
	}{
		{"exactly at cutoff", now.Add(-10 * time.Minute), true},
		{"within window", now.Add(-5 * time.Minute), true},
		{"at now", now, true},
		{"just before cutoff", now.Add(-10*time.Minute - time.Second), false},
		{"well before cutoff", now.Add(-1 * time.Hour), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := f.Matches(makeSinceEntry(tc.ts))
			if got != tc.want {
				t.Errorf("Matches() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestSinceFilter_String(t *testing.T) {
	f, _ := NewSinceFilter(30 * time.Minute)
	got := f.String()
	want := "since(30m0s)"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
