package filter_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/parser"
)

func mustTime(t *testing.T, s string) time.Time {
	t.Helper()
	v, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Fatalf("mustTime: %v", err)
	}
	return v
}

func makeEntry(t *testing.T, ts string) parser.Entry {
	t.Helper()
	return parser.Entry{Timestamp: mustTime(t, ts), Level: "INFO", Message: "test"}
}

func TestNewTimeRange_InvalidFrom(t *testing.T) {
	_, err := filter.NewTimeRange("not-a-date", "")
	if err == nil {
		t.Fatal("expected error for invalid from timestamp")
	}
}

func TestNewTimeRange_ToBeforeFrom(t *testing.T) {
	_, err := filter.NewTimeRange("2024-01-02T00:00:00Z", "2024-01-01T00:00:00Z")
	if err == nil {
		t.Fatal("expected error when to is before from")
	}
}

func TestTimeRange_Contains(t *testing.T) {
	tr, err := filter.NewTimeRange("2024-06-01T00:00:00Z", "2024-06-30T23:59:59Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cases := []struct {
		ts   string
		want bool
	}{
		{"2024-05-31T23:59:59Z", false},
		{"2024-06-01T00:00:00Z", true},
		{"2024-06-15T12:00:00Z", true},
		{"2024-06-30T23:59:59Z", true},
		{"2024-07-01T00:00:00Z", false},
	}

	for _, tc := range cases {
		e := makeEntry(t, tc.ts)
		if got := tr.Contains(e); got != tc.want {
			t.Errorf("Contains(%s) = %v, want %v", tc.ts, got, tc.want)
		}
	}
}

func TestTimeRange_Apply(t *testing.T) {
	tr, _ := filter.NewTimeRange("2024-06-10T00:00:00Z", "2024-06-20T00:00:00Z")

	entries := []parser.Entry{
		makeEntry(t, "2024-06-09T00:00:00Z"),
		makeEntry(t, "2024-06-10T00:00:00Z"),
		makeEntry(t, "2024-06-15T06:00:00Z"),
		makeEntry(t, "2024-06-20T00:00:00Z"),
		makeEntry(t, "2024-06-21T00:00:00Z"),
	}

	got := tr.Apply(entries)
	if len(got) != 3 {
		t.Fatalf("Apply returned %d entries, want 3", len(got))
	}
}
