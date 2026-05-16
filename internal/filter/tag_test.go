package filter

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeTagEntry(msg string, fields map[string]string) parser.Entry {
	if fields == nil {
		fields = map[string]string{}
	}
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   msg,
		Fields:    fields,
	}
}

func TestNewTagFilter_EmptyList(t *testing.T) {
	_, err := NewTagFilter([]string{})
	if err == nil {
		t.Fatal("expected error for empty tag list")
	}
}

func TestNewTagFilter_BlankTag(t *testing.T) {
	_, err := NewTagFilter([]string{"valid", ""})
	if err == nil {
		t.Fatal("expected error for blank tag")
	}
}

func TestTagFilter_Matches_MessageHit(t *testing.T) {
	f, err := NewTagFilter([]string{"error"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := makeTagEntry("an ERROR occurred", nil)
	if !f.Matches(e) {
		t.Error("expected match on message substring")
	}
}

func TestTagFilter_Matches_FieldHit(t *testing.T) {
	f, err := NewTagFilter([]string{"database"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := makeTagEntry("connection failed", map[string]string{"component": "Database"})
	if !f.Matches(e) {
		t.Error("expected match on field value")
	}
}

func TestTagFilter_Matches_MultipleTagsAllPresent(t *testing.T) {
	f, err := NewTagFilter([]string{"auth", "timeout"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := makeTagEntry("auth timeout exceeded", nil)
	if !f.Matches(e) {
		t.Error("expected match when all tags present")
	}
}

func TestTagFilter_Matches_MultipleTagsMissing(t *testing.T) {
	f, err := NewTagFilter([]string{"auth", "timeout"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := makeTagEntry("auth failure", nil)
	if f.Matches(e) {
		t.Error("expected no match when not all tags present")
	}
}

func TestTagFilter_String(t *testing.T) {
	f, _ := NewTagFilter([]string{"foo", "bar"})
	got := f.String()
	expected := "tag(foo,bar)"
	if got != expected {
		t.Errorf("String() = %q, want %q", got, expected)
	}
}
