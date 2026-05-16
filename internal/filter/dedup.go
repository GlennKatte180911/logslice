package filter

import (
	"fmt"

	"github.com/iand/logslice/internal/parser"
)

// DedupFilter suppresses consecutive duplicate log messages.
type DedupFilter struct {
	lastMessage string
	seen        bool
}

// NewDedupFilter creates a DedupFilter that drops entries whose message is
// identical to the immediately preceding matched entry.
func NewDedupFilter() *DedupFilter {
	return &DedupFilter{}
}

// Matches returns false when the entry's message equals the previous entry's
// message, effectively deduplicating consecutive identical lines.
func (d *DedupFilter) Matches(e parser.Entry) bool {
	if d.seen && e.Message == d.lastMessage {
		return false
	}
	d.lastMessage = e.Message
	d.seen = true
	return true
}

// String returns a human-readable description of the filter.
func (d *DedupFilter) String() string {
	return fmt.Sprintf("dedup()")
}

// Reset clears the remembered last message.
func (d *DedupFilter) Reset() {
	d.lastMessage = ""
	d.seen = false
}
