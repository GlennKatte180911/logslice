package filter_test

import (
	"testing"
	"time"

	"github.com/iand/logslice/internal/filter"
	"github.com/iand/logslice/internal/parser"
)

func makeSampleEntry(msg string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   msg,
		Fields:    map[string]string{},
	}
}

func TestNewSampleFilter_NilInner(t *testing.T) {
	_, err := filter.NewSampleFilter(nil, 2)
	if err == nil {
		t.Fatal("expected error for nil inner filter")
	}
}

func TestNewSampleFilter_ZeroN(t *testing.T) {
	pf, _ := filter.NewPatternFilter(".*")
	_, err := filter.NewSampleFilter(pf, 0)
	if err == nil {
		t.Fatal("expected error for n=0")
	}
}

func TestNewSampleFilter_NegativeN(t *testing.T) {
	pf, _ := filter.NewPatternFilter(".*")
	_, err := filter.NewSampleFilter(pf, -3)
	if err == nil {
		t.Fatal("expected error for negative n")
	}
}

func TestSampleFilter_EveryN(t *testing.T) {
	pf, _ := filter.NewPatternFilter(".*")
	sf, err := filter.NewSampleFilter(pf, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	results := make([]bool, 9)
	for i := range results {
		results[i] = sf.Matches(makeSampleEntry("msg"))
	}

	expected := []bool{false, false, true, false, false, true, false, false, true}
	for i, got := range results {
		if got != expected[i] {
			t.Errorf("entry %d: got %v, want %v", i+1, got, expected[i])
		}
	}
}

func TestSampleFilter_N1PassesAll(t *testing.T) {
	pf, _ := filter.NewPatternFilter(".*")
	sf, _ := filter.NewSampleFilter(pf, 1)

	for i := 0; i < 5; i++ {
		if !sf.Matches(makeSampleEntry("msg")) {
			t.Errorf("entry %d: expected match with n=1", i+1)
		}
	}
}

func TestSampleFilter_InnerFilterRespected(t *testing.T) {
	pf, _ := filter.NewPatternFilter("keep")
	sf, _ := filter.NewSampleFilter(pf, 2)

	// only "keep" entries count toward the sample counter
	if sf.Matches(makeSampleEntry("skip")) {
		t.Error("non-matching inner entry should not pass")
	}
	if sf.Matches(makeSampleEntry("keep")) {
		t.Error("1st matching inner entry should not pass (n=2)")
	}
	if !sf.Matches(makeSampleEntry("keep")) {
		t.Error("2nd matching inner entry should pass (n=2)")
	}
}

func TestSampleFilter_Reset(t *testing.T) {
	pf, _ := filter.NewPatternFilter(".*")
	sf, _ := filter.NewSampleFilter(pf, 2)

	sf.Matches(makeSampleEntry("a")) // count=1
	sf.Matches(makeSampleEntry("b")) // count=2, passes
	sf.Reset()

	if sf.Matches(makeSampleEntry("c")) {
		t.Error("after reset first entry should not pass")
	}
	if !sf.Matches(makeSampleEntry("d")) {
		t.Error("after reset second entry should pass")
	}
}

func TestSampleFilter_String(t *testing.T) {
	pf, _ := filter.NewPatternFilter("foo")
	sf, _ := filter.NewSampleFilter(pf, 5)
	s := sf.String()
	if s == "" {
		t.Error("String() should not be empty")
	}
}
