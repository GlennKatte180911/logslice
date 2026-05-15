package filter

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeNumericEntry(key, value string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   "test",
		Fields:    map[string]string{key: value},
	}
}

func TestNewNumericFieldFilter_EmptyKey(t *testing.T) {
	_, err := NewNumericFieldFilter("", OpEq, 1.0)
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestNewNumericFieldFilter_UnknownOp(t *testing.T) {
	_, err := NewNumericFieldFilter("latency", NumericOp("neq"), 1.0)
	if err == nil {
		t.Fatal("expected error for unknown operator")
	}
}

func TestNumericFieldFilter_Matches(t *testing.T) {
	tests := []struct {
		name      string
		op        NumericOp
		threshold float64
		fieldVal  string
		want      bool
	}{
		{"eq match", OpEq, 42.0, "42", true},
		{"eq no match", OpEq, 42.0, "43", false},
		{"lt match", OpLt, 10.0, "9.5", true},
		{"lt no match", OpLt, 10.0, "10", false},
		{"gt match", OpGt, 5.0, "6", true},
		{"gt no match", OpGt, 5.0, "5", false},
		{"lte match equal", OpLte, 7.0, "7", true},
		{"lte match less", OpLte, 7.0, "6", true},
		{"lte no match", OpLte, 7.0, "8", false},
		{"gte match equal", OpGte, 3.0, "3", true},
		{"gte no match", OpGte, 3.0, "2", false},
		{"missing field", OpEq, 1.0, "", false},
		{"non-numeric value", OpEq, 1.0, "abc", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := NewNumericFieldFilter("val", tc.op, tc.threshold)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			var e parser.Entry
			if tc.fieldVal != "" {
				e = makeNumericEntry("val", tc.fieldVal)
			} else {
				e = parser.Entry{Fields: map[string]string{}}
			}
			if got := f.Matches(e); got != tc.want {
				t.Errorf("Matches() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestNumericFieldFilter_String(t *testing.T) {
	f, _ := NewNumericFieldFilter("latency", OpGt, 100.0)
	got := f.String()
	want := "numeric_field(latency gt 100)"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
