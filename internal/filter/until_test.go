package filter

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeUntilEntry(ts time.Time) parser.Entry {
	return parser.Entry{
		Timestamp: ts,
		Level:     "INFO",
		Message:   "until test",
		Fields:    map[string]string{},
	}
}

func TestNewUntilFilter_ZeroTime(t *testing.T) {
	_, err := NewUntilFilter(time.Time{})
	if err == nil {
		t.Fatal("expected error for zero cutoff, got nil")
	}
}

func TestNewUntilFilter_ValidCutoff(t *testing.T) {
	cutoff := time.Now()
	f, err := NewUntilFilter(cutoff)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestUntilFilter_Matches(t *testing.T) {
	cutoff := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	f, _ := NewUntilFilter(cutoff)

	before := makeUntilEntry(cutoff.Add(-time.Minute))
	at := makeUntilEntry(cutoff)
	after := makeUntilEntry(cutoff.Add(time.Minute))

	if !f.Matches(before) {
		t.Error("expected entry before cutoff to match")
	}
	if !f.Matches(at) {
		t.Error("expected entry at cutoff to match")
	}
	if f.Matches(after) {
		t.Error("expected entry after cutoff not to match")
	}
}

func TestUntilFilter_LatchesFalse(t *testing.T) {
	cutoff := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	f, _ := NewUntilFilter(cutoff)

	after := makeUntilEntry(cutoff.Add(time.Hour))
	before := makeUntilEntry(cutoff.Add(-time.Hour))

	// Trigger the latch.
	if f.Matches(after) {
		t.Error("expected first after-cutoff entry not to match")
	}
	// Even a before-cutoff entry should now be rejected.
	if f.Matches(before) {
		t.Error("expected entry to be rejected after latch is set")
	}
}

func TestUntilFilter_String(t *testing.T) {
	cutoff := time.Date(2024, 6, 15, 8, 30, 0, 0, time.UTC)
	f, _ := NewUntilFilter(cutoff)
	s := f.String()
	if s != "until(2024-06-15T08:30:00Z)" {
		t.Errorf("unexpected String() output: %q", s)
	}
}
