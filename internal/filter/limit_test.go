package filter_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/parser"
)

func makeLimitEntry(msg string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   msg,
		Fields:    map[string]string{},
	}
}

func alwaysMatch() filter.Filter {
	p, _ := filter.NewPatternFilter(".*")
	return p
}

func TestNewLimitFilter_NilInner(t *testing.T) {
	_, err := filter.NewLimitFilter(nil, 5)
	if err == nil {
		t.Fatal("expected error for nil inner filter")
	}
}

func TestNewLimitFilter_ZeroMax(t *testing.T) {
	_, err := filter.NewLimitFilter(alwaysMatch(), 0)
	if err == nil {
		t.Fatal("expected error for max=0")
	}
}

func TestNewLimitFilter_NegativeMax(t *testing.T) {
	_, err := filter.NewLimitFilter(alwaysMatch(), -3)
	if err == nil {
		t.Fatal("expected error for negative max")
	}
}

func TestLimitFilter_StopsAfterMax(t *testing.T) {
	f, err := filter.NewLimitFilter(alwaysMatch(), 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	results := make([]bool, 6)
	for i := range results {
		results[i] = f.Matches(makeLimitEntry("msg"))
	}

	expected := []bool{true, true, true, false, false, false}
	for i, got := range results {
		if got != expected[i] {
			t.Errorf("entry %d: got %v, want %v", i, got, expected[i])
		}
	}
	if f.Matched() != 3 {
		t.Errorf("Matched() = %d, want 3", f.Matched())
	}
}

func TestLimitFilter_Reset(t *testing.T) {
	f, _ := filter.NewLimitFilter(alwaysMatch(), 2)
	f.Matches(makeLimitEntry("a"))
	f.Matches(makeLimitEntry("b"))

	if f.Matches(makeLimitEntry("c")) {
		t.Fatal("expected false after limit reached")
	}

	f.Reset()
	if !f.Matches(makeLimitEntry("d")) {
		t.Fatal("expected true after reset")
	}
	if f.Matched() != 1 {
		t.Errorf("Matched() after reset = %d, want 1", f.Matched())
	}
}

func TestLimitFilter_String(t *testing.T) {
	f, _ := filter.NewLimitFilter(alwaysMatch(), 10)
	s := f.String()
	if s == "" {
		t.Fatal("String() returned empty string")
	}
}
