package filter

import (
	"testing"
	"time"

	"github.com/nicholasgasior/logslice/internal/parser"
)

func makeBurstEntry(ts time.Time, msg string) parser.Entry {
	return parser.Entry{Timestamp: ts, Message: msg, Fields: map[string]string{}}
}

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func alwaysMatchFilter() Filter { return &patternAlwaysMatch{} }

type patternAlwaysMatch struct{}

func (p *patternAlwaysMatch) Matches(_ parser.Entry) bool { return true }
func (p *patternAlwaysMatch) String() string               { return "always" }

func TestNewBurstFilter_NilInner(t *testing.T) {
	_, err := NewBurstFilter(nil, 3, time.Second)
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewBurstFilter_ZeroMax(t *testing.T) {
	_, err := NewBurstFilter(alwaysMatchFilter(), 0, time.Second)
	if err == nil {
		t.Fatal("expected error for maxCount=0")
	}
}

func TestNewBurstFilter_NonPositiveWindow(t *testing.T) {
	_, err := NewBurstFilter(alwaysMatchFilter(), 3, 0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestBurstFilter_AllowsUpToMax(t *testing.T) {
	f, err := NewBurstFilter(alwaysMatchFilter(), 3, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 3; i++ {
		e := makeBurstEntry(epoch.Add(time.Duration(i)*time.Second), "msg")
		if !f.Matches(e) {
			t.Errorf("entry %d should match", i)
		}
	}
}

func TestBurstFilter_DropsAfterMax(t *testing.T) {
	f, err := NewBurstFilter(alwaysMatchFilter(), 2, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	f.Matches(makeBurstEntry(epoch, "a"))
	f.Matches(makeBurstEntry(epoch.Add(time.Second), "b"))
	if f.Matches(makeBurstEntry(epoch.Add(2*time.Second), "c")) {
		t.Error("third entry within window should be dropped")
	}
}

func TestBurstFilter_ResetsAfterWindow(t *testing.T) {
	f, err := NewBurstFilter(alwaysMatchFilter(), 2, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	f.Matches(makeBurstEntry(epoch, "a"))
	f.Matches(makeBurstEntry(epoch.Add(time.Second), "b"))
	// New window starts after 1 minute
	if !f.Matches(makeBurstEntry(epoch.Add(2*time.Minute), "c")) {
		t.Error("entry in new window should match")
	}
}

func TestBurstFilter_String(t *testing.T) {
	f, _ := NewBurstFilter(alwaysMatchFilter(), 5, 30*time.Second)
	s := f.String()
	if s == "" {
		t.Error("String() should not be empty")
	}
}
