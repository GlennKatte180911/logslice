package filter_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/parser"
)

func makeBetweenEntry(msg string, fields map[string]string) parser.Entry {
	e := parser.Entry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   msg,
		Fields:    make(map[string]string),
	}
	for k, v := range fields {
		e.Fields[k] = v
	}
	return e
}

func TestNewBetweenFilter_EmptyLow(t *testing.T) {
	_, err := filter.NewBetweenFilter("", "", "z")
	if err == nil {
		t.Fatal("expected error for empty low bound")
	}
}

func TestNewBetweenFilter_EmptyHigh(t *testing.T) {
	_, err := filter.NewBetweenFilter("", "a", "")
	if err == nil {
		t.Fatal("expected error for empty high bound")
	}
}

func TestNewBetweenFilter_LowGreaterThanHigh(t *testing.T) {
	_, err := filter.NewBetweenFilter("", "z", "a")
	if err == nil {
		t.Fatal("expected error when low > high")
	}
}

func TestBetweenFilter_Matches_Message(t *testing.T) {
	f, err := filter.NewBetweenFilter("", "beta", "delta")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cases := []struct {
		msg  string
		want bool
	}{
		{"alpha", false},
		{"beta", true},
		{"charlie", true},
		{"delta", true},
		{"echo", false},
	}
	for _, tc := range cases {
		got := f.Matches(makeBetweenEntry(tc.msg, nil))
		if got != tc.want {
			t.Errorf("msg=%q: got %v, want %v", tc.msg, got, tc.want)
		}
	}
}

func TestBetweenFilter_Matches_Field(t *testing.T) {
	f, err := filter.NewBetweenFilter("code", "200", "299")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cases := []struct {
		fields map[string]string
		want   bool
	}{
		{map[string]string{"code": "199"}, false},
		{map[string]string{"code": "200"}, true},
		{map[string]string{"code": "250"}, true},
		{map[string]string{"code": "299"}, true},
		{map[string]string{"code": "300"}, false},
		{map[string]string{"status": "200"}, false}, // wrong key
	}
	for _, tc := range cases {
		got := f.Matches(makeBetweenEntry("msg", tc.fields))
		if got != tc.want {
			t.Errorf("fields=%v: got %v, want %v", tc.fields, got, tc.want)
		}
	}
}

func TestBetweenFilter_String(t *testing.T) {
	f1, _ := filter.NewBetweenFilter("", "a", "z")
	if s := f1.String(); s != `between(message, "a", "z")` {
		t.Errorf("unexpected String: %s", s)
	}

	f2, _ := filter.NewBetweenFilter("env", "prod", "staging")
	if s := f2.String(); s != `between(env, "prod", "staging")` {
		t.Errorf("unexpected String: %s", s)
	}
}
