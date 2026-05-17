package filter

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeFirstLastEntry(msg string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   msg,
		Fields:    map[string]string{},
	}
}

// alwaysMatchFL is a local always-match helper to avoid cross-test pollution.
type alwaysMatchFL struct{}

func (a *alwaysMatchFL) Matches(_ parser.Entry) bool { return true }
func (a *alwaysMatchFL) String() string               { return "always" }

// --- FirstFilter ---

func TestNewFirstFilter_NilInner(t *testing.T) {
	_, err := NewFirstFilter(nil, 3)
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewFirstFilter_ZeroN(t *testing.T) {
	_, err := NewFirstFilter(&alwaysMatchFL{}, 0)
	if err == nil {
		t.Fatal("expected error for n=0")
	}
}

func TestFirstFilter_PassesFirstN(t *testing.T) {
	f, err := NewFirstFilter(&alwaysMatchFL{}, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 3; i++ {
		if !f.Matches(makeFirstLastEntry("msg")) {
			t.Errorf("entry %d should match", i)
		}
	}
	if f.Matches(makeFirstLastEntry("msg")) {
		t.Error("4th entry should not match")
	}
}

func TestFirstFilter_String(t *testing.T) {
	f, _ := NewFirstFilter(&alwaysMatchFL{}, 5)
	s := f.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}

// --- LastFilter ---

func TestNewLastFilter_NilInner(t *testing.T) {
	_, err := NewLastFilter(nil, 3)
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewLastFilter_ZeroN(t *testing.T) {
	_, err := NewLastFilter(&alwaysMatchFL{}, 0)
	if err == nil {
		t.Fatal("expected error for n=0")
	}
}

func TestLastFilter_CollectRetainsTail(t *testing.T) {
	f, err := NewLastFilter(&alwaysMatchFL{}, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	msgs := []string{"a", "b", "c", "d", "e"}
	for _, m := range msgs {
		f.Matches(makeFirstLastEntry(m))
	}
	tail := f.Collect()
	if len(tail) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(tail))
	}
	expected := []string{"c", "d", "e"}
	for i, e := range tail {
		if e.Message != expected[i] {
			t.Errorf("entry %d: want %q, got %q", i, expected[i], e.Message)
		}
	}
}

func TestLastFilter_MatchesAlwaysFalse(t *testing.T) {
	f, _ := NewLastFilter(&alwaysMatchFL{}, 2)
	if f.Matches(makeFirstLastEntry("x")) {
		t.Error("Matches should always return false for LastFilter")
	}
}

func TestLastFilter_String(t *testing.T) {
	f, _ := NewLastFilter(&alwaysMatchFL{}, 2)
	if f.String() == "" {
		t.Error("expected non-empty string")
	}
}
