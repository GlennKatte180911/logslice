package filter

import (
	"fmt"

	"github.com/user/logslice/internal/parser"
)

// FieldExistsFilter matches entries that contain (or do not contain) a specific field key.
type FieldExistsFilter struct {
	key    string
	negate bool // if true, matches entries where the field is absent
}

// NewFieldExistsFilter returns a filter that matches entries containing the given field key.
// If negate is true, it matches entries where the field is absent.
func NewFieldExistsFilter(key string, negate bool) (*FieldExistsFilter, error) {
	if key == "" {
		return nil, fmt.Errorf("field_exists: key must not be empty")
	}
	return &FieldExistsFilter{key: key, negate: negate}, nil
}

// Matches returns true when the entry's Fields map contains (or lacks, if negated) the key.
func (f *FieldExistsFilter) Matches(e parser.Entry) bool {
	_, ok := e.Fields[f.key]
	if f.negate {
		return !ok
	}
	return ok
}

// String returns a human-readable description of the filter.
func (f *FieldExistsFilter) String() string {
	if f.negate {
		return fmt.Sprintf("field_missing(%q)", f.key)
	}
	return fmt.Sprintf("field_exists(%q)", f.key)
}
