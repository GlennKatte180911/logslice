package filter_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/parser"
)

func makeChainEntry(level, msg string, fields map[string]string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
		Fields:    fields,
	}
}

func TestChain_EmptyMatchesAll(t *testing.T) {
	c := filter.NewChain()
	e := makeChainEntry("info", "hello", nil)
	if !c.Matches(e) {
		t.Error("empty chain should match all entries")
	}
}

func TestChain_SingleFilter(t *testing.T) {
	lf, _ := filter.NewLevelFilter([]string{"error"})
	c := filter.NewChain(lf)

	if c.Matches(makeChainEntry("info", "msg", nil)) {
		t.Error("chain should not match info entry when filtering for error")
	}
	if !c.Matches(makeChainEntry("error", "msg", nil)) {
		t.Error("chain should match error entry")
	}
}

func TestChain_MultipleFilters_AND(t *testing.T) {
	lf, _ := filter.NewLevelFilter([]string{"warn"})
	ff, _ := filter.NewFieldFilter("service", "api")
	c := filter.NewChain(lf, ff)

	tests := []struct {
		name   string
		entry  parser.Entry
		want   bool
	}{
		{"both match", makeChainEntry("warn", "m", map[string]string{"service": "api"}), true},
		{"only level", makeChainEntry("warn", "m", map[string]string{"service": "db"}), false},
		{"only field", makeChainEntry("info", "m", map[string]string{"service": "api"}), false},
		{"neither", makeChainEntry("info", "m", nil), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := c.Matches(tc.entry); got != tc.want {
				t.Errorf("Matches() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestChain_Add(t *testing.T) {
	c := filter.NewChain()
	lf, _ := filter.NewLevelFilter([]string{"debug"})
	c.Add(lf)

	if c.Matches(makeChainEntry("info", "msg", nil)) {
		t.Error("chain should not match after adding level filter for debug")
	}
}

func TestChain_String(t *testing.T) {
	lf, _ := filter.NewLevelFilter([]string{"error"})
	ff, _ := filter.NewFieldFilter("env", "prod")
	c := filter.NewChain(lf, ff)
	got := c.String()
	if got == "" {
		t.Error("String() should not be empty")
	}
}
