package filter

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeNotEntry(level, msg string, fields map[string]string) parser.Entry {
	if fields == nil {
		fields = map[string]string{}
	}
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
		Fields:    fields,
	}
}

func TestNewNotFilter_NilInner(t *testing.T) {
	_, err := NewNotFilter(nil)
	if err == nil {
		t.Fatal("expected error for nil inner filter, got nil")
	}
}

func TestNotFilter_InvertsMatch(t *testing.T) {
	inner, err := NewLevelFilter([]string{"error"})
	if err != nil {
		t.Fatalf("unexpected error building level filter: %v", err)
	}

	nf, err := NewNotFilter(inner)
	if err != nil {
		t.Fatalf("unexpected error building not filter: %v", err)
	}

	errorEntry := makeNotEntry("error", "boom", nil)
	infoEntry := makeNotEntry("info", "ok", nil)

	if nf.Matches(errorEntry) {
		t.Error("expected NOT(error) to reject an error-level entry")
	}
	if !nf.Matches(infoEntry) {
		t.Error("expected NOT(error) to accept an info-level entry")
	}
}

func TestNotFilter_String(t *testing.T) {
	inner, _ := NewLevelFilter([]string{"warn"})
	nf, _ := NewNotFilter(inner)

	got := nf.String()
	want := "NOT(level in [warn])"
	if got != want {
		t.Errorf("String() = %q; want %q", got, want)
	}
}

func TestNotFilter_WithPatternFilter(t *testing.T) {
	inner, err := NewPatternFilter("timeout")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	nf, err := NewNotFilter(inner)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	match := makeNotEntry("error", "connection timeout occurred", nil)
	nomatch := makeNotEntry("info", "all systems nominal", nil)

	if nf.Matches(match) {
		t.Error("expected NOT(pattern) to reject entry matching the pattern")
	}
	if !nf.Matches(nomatch) {
		t.Error("expected NOT(pattern) to accept entry not matching the pattern")
	}
}
