package filter_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/parser"
)

func makeThrottleEntry(ts time.Time, msg string) parser.Entry {
	return parser.Entry{Timestamp: ts, Message: msg, Fields: map[string]string{}}
}

var alwaysMatchThrottle = &alwaysMatchFilter{}

func TestNewThrottleFilter_NilInner(t *testing.T) {
	_, err := filter.NewThrottleFilter(nil, time.Second)
	if err == nil {
		t.Fatal("expected error for nil inner, got nil")
	}
}

func TestNewThrottleFilter_ZeroWindow(t *testing.T) {
	_, err := filter.NewThrottleFilter(alwaysMatchThrottle, 0)
	if err == nil {
		t.Fatal("expected error for zero window, got nil")
	}
}

func TestNewThrottleFilter_NegativeWindow(t *testing.T) {
	_, err := filter.NewThrottleFilter(alwaysMatchThrottle, -time.Second)
	if err == nil {
		t.Fatal("expected error for negative window, got nil")
	}
}

func TestThrottleFilter_FirstEntryPasses(t *testing.T) {
	f, err := filter.NewThrottleFilter(alwaysMatchThrottle, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if !f.Matches(makeThrottleEntry(base, "first")) {
		t.Error("expected first entry to pass")
	}
}

func TestThrottleFilter_SecondEntryWithinWindowDropped(t *testing.T) {
	f, _ := filter.NewThrottleFilter(alwaysMatchThrottle, time.Minute)
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	f.Matches(makeThrottleEntry(base, "first"))
	if f.Matches(makeThrottleEntry(base.Add(30*time.Second), "second")) {
		t.Error("expected entry within window to be dropped")
	}
}

func TestThrottleFilter_EntryAfterWindowPasses(t *testing.T) {
	f, _ := filter.NewThrottleFilter(alwaysMatchThrottle, time.Minute)
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	f.Matches(makeThrottleEntry(base, "first"))
	if !f.Matches(makeThrottleEntry(base.Add(time.Minute), "second")) {
		t.Error("expected entry exactly at window boundary to pass")
	}
}

func TestThrottleFilter_InnerFilterRespected(t *testing.T) {
	neverMatch, _ := filter.NewPatternFilter("^$") // matches empty string only
	f, _ := filter.NewThrottleFilter(neverMatch, time.Second)
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if f.Matches(makeThrottleEntry(base, "non-empty message")) {
		t.Error("expected inner filter rejection to propagate")
	}
}

func TestThrottleFilter_String(t *testing.T) {
	f, _ := filter.NewThrottleFilter(alwaysMatchThrottle, 5*time.Second)
	s := f.String()
	if s == "" {
		t.Error("expected non-empty String()")
	}
}
