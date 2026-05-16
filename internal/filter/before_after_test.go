package filter

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeBAEntry(ts time.Time) parser.Entry {
	return parser.Entry{Timestamp: ts, Level: "info", Message: "test", Fields: map[string]string{}}
}

var (
	baBase = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	baBefore = baBase.Add(-time.Hour)
	baAfter  = baBase.Add(time.Hour)
)

func TestNewBeforeFilter_ZeroTime(t *testing.T) {
	_, err := NewBeforeFilter(time.Time{})
	if err == nil {
		t.Fatal("expected error for zero cutoff")
	}
}

func TestNewAfterFilter_ZeroTime(t *testing.T) {
	_, err := NewAfterFilter(time.Time{})
	if err == nil {
		t.Fatal("expected error for zero cutoff")
	}
}

func TestBeforeFilter_Matches(t *testing.T) {
	f, err := NewBeforeFilter(baBase)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		name  string
		ts    time.Time
		want  bool
	}{
		{"before cutoff", baBefore, true},
		{"at cutoff", baBase, false},
		{"after cutoff", baAfter, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := f.Matches(makeBAEntry(tc.ts)); got != tc.want {
				t.Errorf("Matches() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestAfterFilter_Matches(t *testing.T) {
	f, err := NewAfterFilter(baBase)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		name  string
		ts    time.Time
		want  bool
	}{
		{"before cutoff", baBefore, false},
		{"at cutoff", baBase, false},
		{"after cutoff", baAfter, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := f.Matches(makeBAEntry(tc.ts)); got != tc.want {
				t.Errorf("Matches() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestBeforeFilter_String(t *testing.T) {
	f, _ := NewBeforeFilter(baBase)
	s := f.String()
	if s == "" {
		t.Error("String() returned empty string")
	}
	if s[:7] != "before(" {
		t.Errorf("String() = %q, want prefix 'before('", s)
	}
}

func TestAfterFilter_String(t *testing.T) {
	f, _ := NewAfterFilter(baBase)
	s := f.String()
	if s == "" {
		t.Error("String() returned empty string")
	}
	if s[:6] != "after(" {
		t.Errorf("String() = %q, want prefix 'after('", s)
	}
}
