package filter

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// TagFilter matches log entries whose message or fields contain all of the
// specified tags (case-insensitive substring match on field values and message).
type TagFilter struct {
	tags []string
}

// NewTagFilter returns a TagFilter that passes entries containing every tag in
// the supplied list. It returns an error if the list is empty or any tag is
// blank.
func NewTagFilter(tags []string) (*TagFilter, error) {
	if len(tags) == 0 {
		return nil, fmt.Errorf("tag filter: tag list must not be empty")
	}
	normalized := make([]string, 0, len(tags))
	for _, t := range tags {
		t = strings.TrimSpace(t)
		if t == "" {
			return nil, fmt.Errorf("tag filter: blank tag is not allowed")
		}
		normalized = append(normalized, strings.ToLower(t))
	}
	return &TagFilter{tags: normalized}, nil
}

// Matches returns true when the entry contains every configured tag.
func (f *TagFilter) Matches(e parser.Entry) bool {
	// Build a searchable corpus from the message and all field values.
	parts := make([]string, 0, 1+len(e.Fields))
	parts = append(parts, strings.ToLower(e.Message))
	for _, v := range e.Fields {
		parts = append(parts, strings.ToLower(v))
	}
	corpus := strings.Join(parts, " ")

	for _, tag := range f.tags {
		if !strings.Contains(corpus, tag) {
			return false
		}
	}
	return true
}

// String returns a human-readable description of the filter.
func (f *TagFilter) String() string {
	return fmt.Sprintf("tag(%s)", strings.Join(f.tags, ","))
}
