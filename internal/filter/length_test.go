package filter

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeLengthEntry(msg string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   msg,
		Fields:    map[string]string{},
	}
}

func TestNewLengthFilter_UnknownOp(t *testing.T) {
	_, err := NewLengthFilter("neq", 5)
	if err == nil {
		t.Fatal("expected error for unknown op")
	}
}

func TestNewLengthFilter_NegativeVal(t *testing.T) {
	_, err := NewLengthFilter("gt", -1)
	if err == nil {
		t.Fatal("expected error for negative val")
	}
}

func TestNewLengthFilter_ValidOps(t *testing.T) {
	ops := []string{"lt", "lte", "gt", "gte", "eq"}
	for _, op := range ops {
		_, err := NewLengthFilter(op, 10)
		if err != nil {
			t.Errorf("op %q: unexpected error: %v", op, err)
		}
	}
}

func TestLengthFilter_Matches(t *testing.T) {
	tests := []struct {
		op      string
		val     int
		msg     string
		want    bool
	}{
		{"lt", 10, "hello", true},
		{"lt", 5, "hello", false},
		{"lte", 5, "hello", true},
		{"gt", 4, "hello", true},
		{"gt", 5, "hello", false},
		{"gte", 5, "hello", true},
		{"gte", 6, "hello", false},
		{"eq", 5, "hello", true},
		{"eq", 4, "hello", false},
	}
	for _, tc := range tests {
		f, err := NewLengthFilter(tc.op, tc.val)
		if err != nil {
			t.Fatalf("%s %d: unexpected error: %v", tc.op, tc.val, err)
		}
		got := f.Matches(makeLengthEntry(tc.msg))
		if got != tc.want {
			t.Errorf("op=%s val=%d msg=%q: got %v, want %v", tc.op, tc.val, tc.msg, got, tc.want)
		}
	}
}

func TestLengthFilter_String(t *testing.T) {
	f, _ := NewLengthFilter("gt", 20)
	s := f.String()
	if s == "" {
		t.Fatal("expected non-empty String()")
	}
}
