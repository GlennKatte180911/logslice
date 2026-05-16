package filter_test

import (
	"testing"
	"time"

	"github.com/iand/logslice/internal/filter"
	"github.com/iand/logslice/internal/parser"
)

func makeDedupEntry(msg string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   msg,
		Fields:    map[string]string{},
	}
}

func TestDedupFilter_FirstEntryAlwaysPasses(t *testing.T) {
	d := filter.NewDedupFilter()
	if !d.Matches(makeDedupEntry("hello")) {
		t.Error("first entry should always pass")
	}
}

func TestDedupFilter_ConsecutiveDuplicateDropped(t *testing.T) {
	d := filter.NewDedupFilter()
	d.Matches(makeDedupEntry("hello"))
	if d.Matches(makeDedupEntry("hello")) {
		t.Error("consecutive duplicate should be dropped")
	}
}

func TestDedupFilter_DifferentMessagePasses(t *testing.T) {
	d := filter.NewDedupFilter()
	d.Matches(makeDedupEntry("hello"))
	if !d.Matches(makeDedupEntry("world")) {
		t.Error("different message should pass")
	}
}

func TestDedupFilter_NonConsecutiveDuplicatePasses(t *testing.T) {
	d := filter.NewDedupFilter()
	d.Matches(makeDedupEntry("hello"))
	d.Matches(makeDedupEntry("world"))
	if !d.Matches(makeDedupEntry("hello")) {
		t.Error("non-consecutive duplicate should pass")
	}
}

func TestDedupFilter_Reset(t *testing.T) {
	d := filter.NewDedupFilter()
	d.Matches(makeDedupEntry("hello"))
	d.Reset()
	if !d.Matches(makeDedupEntry("hello")) {
		t.Error("after reset, same message should pass again")
	}
}

func TestDedupFilter_String(t *testing.T) {
	d := filter.NewDedupFilter()
	if d.String() == "" {
		t.Error("String() should not be empty")
	}
}

func TestDedupFilter_MultipleConsecutiveDuplicates(t *testing.T) {
	d := filter.NewDedupFilter()
	msgs := []string{"a", "a", "a", "b", "b", "c"}
	expected := []bool{true, false, false, true, false, true}
	for i, msg := range msgs {
		got := d.Matches(makeDedupEntry(msg))
		if got != expected[i] {
			t.Errorf("entry %d (%q): got %v, want %v", i, msg, got, expected[i])
		}
	}
}
