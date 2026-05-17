package filter

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeWindowEntry(ts time.Time) *parser.Entry {
	return &parser.Entry{Timestamp: ts, Message: "msg", Fields: map[string]string{}}
}

var windowAnchor = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func TestNewWindowFilter_ZeroAnchor(t *testing.T) {
	_, err := NewWindowFilter(time.Time{}, time.Minute, time.Minute)
	if err == nil {
		t.Fatal("expected error for zero anchor")
	}
}

func TestNewWindowFilter_NegativeBefore(t *testing.T) {
	_, err := NewWindowFilter(windowAnchor, -time.Second, time.Minute)
	if err == nil {
		t.Fatal("expected error for negative before")
	}
}

func TestNewWindowFilter_NegativeAfter(t *testing.T) {
	_, err := NewWindowFilter(windowAnchor, time.Minute, -time.Second)
	if err == nil {
		t.Fatal("expected error for negative after")
	}
}

func TestWindowFilter_Matches(t *testing.T) {
	f, err := NewWindowFilter(windowAnchor, 5*time.Minute, 5*time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cases := []struct {
		offset time.Duration
		want   bool
	}{
		{0, true},
		{-5 * time.Minute, true},
		{5 * time.Minute, true},
		{-5*time.Minute - time.Second, false},
		{5*time.Minute + time.Second, false},
	}

	for _, tc := range cases {
		e := makeWindowEntry(windowAnchor.Add(tc.offset))
		if got := f.Matches(e); got != tc.want {
			t.Errorf("offset %s: got %v, want %v", tc.offset, got, tc.want)
		}
	}
}

func TestWindowFilter_String(t *testing.T) {
	f, _ := NewWindowFilter(windowAnchor, 2*time.Minute, 3*time.Minute)
	s := f.String()
	if s == "" {
		t.Fatal("expected non-empty string")
	}
	for _, sub := range []string{"window", "before", "after"} {
		if !containsSubstr(s, sub) {
			t.Errorf("String() missing %q: %s", sub, s)
		}
	}
}

func containsSubstr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && stringContains(s, sub))
}

func stringContains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
