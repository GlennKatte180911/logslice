package filter

import (
	"fmt"
	"strconv"

	"github.com/user/logslice/internal/parser"
)

// NumericOp represents a numeric comparison operator.
type NumericOp string

const (
	OpEq  NumericOp = "eq"
	OpLt  NumericOp = "lt"
	OpGt  NumericOp = "gt"
	OpLte NumericOp = "lte"
	OpGte NumericOp = "gte"
)

// numericFieldFilter matches log entries where a numeric field satisfies
// a comparison against a threshold value.
type numericFieldFilter struct {
	key       string
	op        NumericOp
	threshold float64
}

// NewNumericFieldFilter creates a filter that compares a numeric field value
// against threshold using the given operator (eq, lt, gt, lte, gte).
func NewNumericFieldFilter(key string, op NumericOp, threshold float64) (*numericFieldFilter, error) {
	if key == "" {
		return nil, fmt.Errorf("numeric field filter: key must not be empty")
	}
	switch op {
	case OpEq, OpLt, OpGt, OpLte, OpGte:
		// valid
	default:
		return nil, fmt.Errorf("numeric field filter: unknown operator %q", op)
	}
	return &numericFieldFilter{key: key, op: op, threshold: threshold}, nil
}

// Matches returns true if the entry's field value satisfies the comparison.
func (f *numericFieldFilter) Matches(e parser.Entry) bool {
	raw, ok := e.Fields[f.key]
	if !ok {
		return false
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return false
	}
	switch f.op {
	case OpEq:
		return v == f.threshold
	case OpLt:
		return v < f.threshold
	case OpGt:
		return v > f.threshold
	case OpLte:
		return v <= f.threshold
	case OpGte:
		return v >= f.threshold
	}
	return false
}

// String returns a human-readable description of the filter.
func (f *numericFieldFilter) String() string {
	return fmt.Sprintf("numeric_field(%s %s %g)", f.key, f.op, f.threshold)
}
