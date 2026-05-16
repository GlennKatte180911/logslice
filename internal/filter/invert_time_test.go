package filter

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeInvertTimeEntry(ts time.Time) parser.Entry {
	return parser.Entry{Timestamp: ts, Level: "info", Message: "test"}
}

func TestNewInvertTimeFilter_ZeroFrom(t *testing.T) {
	_, err := NewInvertTimeFilter(time.Time{}, time.Now())
	if err == nil {
		t.Fatal("expected error for zero from time")
	}
}

func TestNewInvertTimeFilter_ZeroTo(t *testing.T) {
	_, err := NewInvertTimeFilter(time.Now(), time.Time{})
	if err == nil {
		t.Fatal("expected error for zero to time")
	}
}

func TestNewInvertTimeFilter_ToBeforeFrom(t *testing.T) {
	now := time.Now()
	_, err := NewInvertTimeFilter(now, now.Add(-time.Hour))
	if err == nil {
		t.Fatal("expected error when to is before from")
	}
}

func TestInvertTimeFilter_Matches(t *testing.T) {
	base := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	from := base
	to := base.Add(2 * time.Hour)

	f, err := NewInvertTimeFilter(from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		name  string
		ts    time.Time
		want  bool
	}{
		{"before range", base.Add(-time.Minute), true},
		{"at from boundary", from, false},
		{"inside range", base.Add(time.Hour), false},
		{"at to boundary", to, false},
		{"after range", to.Add(time.Minute), true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := f.Matches(makeInvertTimeEntry(tc.ts))
			if got != tc.want {
				t.Errorf("Matches() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestInvertTimeFilter_String(t *testing.T) {
	base := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	f, _ := NewInvertTimeFilter(base, base.Add(time.Hour))
	s := f.String()
	if s == "" {
		t.Fatal("expected non-empty string representation")
	}
	if len(s) < 10 {
		t.Errorf("string too short: %q", s)
	}
}
