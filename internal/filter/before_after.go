package filter

import (
	"fmt"
	"time"

	"github.com/user/logslice/internal/parser"
)

// BeforeFilter passes entries whose timestamp is strictly before a cutoff.
type BeforeFilter struct {
	cutoff time.Time
}

// NewBeforeFilter returns a BeforeFilter for the given cutoff time.
func NewBeforeFilter(cutoff time.Time) (*BeforeFilter, error) {
	if cutoff.IsZero() {
		return nil, fmt.Errorf("before filter: cutoff time must not be zero")
	}
	return &BeforeFilter{cutoff: cutoff}, nil
}

// Matches returns true when the entry timestamp is before the cutoff.
func (f *BeforeFilter) Matches(e parser.Entry) bool {
	return e.Timestamp.Before(f.cutoff)
}

// String returns a human-readable description of the filter.
func (f *BeforeFilter) String() string {
	return fmt.Sprintf("before(%s)", f.cutoff.Format(time.RFC3339))
}

// AfterFilter passes entries whose timestamp is strictly after a cutoff.
type AfterFilter struct {
	cutoff time.Time
}

// NewAfterFilter returns an AfterFilter for the given cutoff time.
func NewAfterFilter(cutoff time.Time) (*AfterFilter, error) {
	if cutoff.IsZero() {
		return nil, fmt.Errorf("after filter: cutoff time must not be zero")
	}
	return &AfterFilter{cutoff: cutoff}, nil
}

// Matches returns true when the entry timestamp is after the cutoff.
func (f *AfterFilter) Matches(e parser.Entry) bool {
	return e.Timestamp.After(f.cutoff)
}

// String returns a human-readable description of the filter.
func (f *AfterFilter) String() string {
	return fmt.Sprintf("after(%s)", f.cutoff.Format(time.RFC3339))
}
