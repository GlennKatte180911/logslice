package filter_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/parser"
)

func makeLevelEntry(level string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   "test message",
		Fields:    map[string]string{},
	}
}

func TestNewLevelFilter_EmptyList(t *testing.T) {
	_, err := filter.NewLevelFilter("")
	if err == nil {
		t.Fatal("expected error for empty level list, got nil")
	}
}

func TestNewLevelFilter_BlankEntry(t *testing.T) {
	_, err := filter.NewLevelFilter("ERROR,,WARN")
	if err == nil {
		t.Fatal("expected error for blank entry in list, got nil")
	}
}

func TestLevelFilter_Matches(t *testing.T) {
	f, err := filter.NewLevelFilter("ERROR,WARN")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		level string
		want  bool
	}{
		{"ERROR", true},
		{"WARN", true},
		{"error", true},  // case-insensitive
		{"warn", true},   // case-insensitive
		{"INFO", false},
		{"DEBUG", false},
		{"", false},
	}

	for _, tc := range tests {
		t.Run(tc.level, func(t *testing.T) {
			got := f.Matches(makeLevelEntry(tc.level))
			if got != tc.want {
				t.Errorf("Matches(%q) = %v, want %v", tc.level, got, tc.want)
			}
		})
	}
}

func TestLevelFilter_String(t *testing.T) {
	f, err := filter.NewLevelFilter("ERROR")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := f.String()
	if s == "" {
		t.Error("String() returned empty string")
	}
}
