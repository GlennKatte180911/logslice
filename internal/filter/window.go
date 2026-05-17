package filter

import (
	"fmt"
	"time"

	"github.com/user/logslice/internal/parser"
)

// WindowFilter passes entries whose timestamp falls within a sliding window
// of [anchor - before, anchor + after]. The anchor defaults to now.
type WindowFilter struct {
	anchor time.Time
	before time.Duration
	after  time.Duration
	from   time.Time
	to     time.Time
}

// NewWindowFilter creates a WindowFilter centred on anchor spanning
// [anchor-before, anchor+after]. Both before and after must be >= 0.
func NewWindowFilter(anchor time.Time, before, after time.Duration) (*WindowFilter, error) {
	if anchor.IsZero() {
		return nil, fmt.Errorf("window: anchor time must not be zero")
	}
	if before < 0 {
		return nil, fmt.Errorf("window: before duration must be non-negative, got %s", before)
	}
	if after < 0 {
		return nil, fmt.Errorf("window: after duration must be non-negative, got %s", after)
	}
	return &WindowFilter{
		anchor: anchor,
		before: before,
		after:  after,
		from:   anchor.Add(-before),
		to:     anchor.Add(after),
	}, nil
}

// Matches returns true when the entry's timestamp is within the window.
func (w *WindowFilter) Matches(e *parser.Entry) bool {
	t := e.Timestamp
	return !t.Before(w.from) && !t.After(w.to)
}

func (w *WindowFilter) String() string {
	return fmt.Sprintf("window(anchor=%s before=%s after=%s)",
		w.anchor.Format(time.RFC3339), w.before, w.after)
}
