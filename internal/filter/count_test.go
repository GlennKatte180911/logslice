package filter_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/parser"
)

func makeCountEntry(level, msg string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
		Fields:    map[string]string{},
	}
}

func TestNewCountFilter_NilInner(t *testing.T) {
	_, err := filter.NewCountFilter(nil)
	if err == nil {
		t.Fatal("expected error for nil inner, got nil")
	}
}

func TestCountFilter_TracksMatches(t *testing.T) {
	lf, err := filter.NewLevelFilter([]string{"error"})
	if err != nil {
		t.Fatalf("NewLevelFilter: %v", err)
	}
	cf, err := filter.NewCountFilter(lf)
	if err != nil {
		t.Fatalf("NewCountFilter: %v", err)
	}

	entries := []parser.Entry{
		makeCountEntry("error", "boom"),
		makeCountEntry("info", "ok"),
		makeCountEntry("error", "again"),
		makeCountEntry("debug", "verbose"),
	}
	for _, e := range entries {
		cf.Matches(e)
	}

	if got := cf.Total(); got != 4 {
		t.Errorf("Total: want 4, got %d", got)
	}
	if got := cf.Matched(); got != 2 {
		t.Errorf("Matched: want 2, got %d", got)
	}
}

func TestCountFilter_String(t *testing.T) {
	lf, _ := filter.NewLevelFilter([]string{"warn"})
	cf, _ := filter.NewCountFilter(lf)
	s := cf.String()
	if s == "" {
		t.Error("String() returned empty string")
	}
}

func TestCountFilter_NoMatches(t *testing.T) {
	lf, _ := filter.NewLevelFilter([]string{"fatal"})
	cf, _ := filter.NewCountFilter(lf)

	cf.Matches(makeCountEntry("info", "hello"))
	cf.Matches(makeCountEntry("debug", "world"))

	if cf.Matched() != 0 {
		t.Errorf("expected 0 matches, got %d", cf.Matched())
	}
	if cf.Total() != 2 {
		t.Errorf("expected 2 total, got %d", cf.Total())
	}
}
