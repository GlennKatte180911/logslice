package filter

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// FieldFilter matches log entries where a specific field equals a given value.
type FieldFilter struct {
	key   string
	value string
}

// NewFieldFilter creates a FieldFilter that matches entries where fields[key] == value.
// Both key and value must be non-empty.
func NewFieldFilter(key, value string) (*FieldFilter, error) {
	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)
	if key == "" {
		return nil, fmt.Errorf("field filter: key must not be empty")
	}
	if value == "" {
		return nil, fmt.Errorf("field filter: value must not be empty")
	}
	return &FieldFilter{key: key, value: value}, nil
}

// Matches returns true when the entry's Fields map contains key with the expected value.
func (f *FieldFilter) Matches(e parser.Entry) bool {
	v, ok := e.Fields[f.key]
	if !ok {
		return false
	}
	return v == f.value
}

// String returns a human-readable description of the filter.
func (f *FieldFilter) String() string {
	return fmt.Sprintf("field(%s=%s)", f.key, f.value)
}
