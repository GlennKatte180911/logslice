package main

import (
	"fmt"

	"github.com/user/logslice/internal/filter"
)

// buildLengthFilter constructs a LengthFilter from the parsed flags and adds
// it to the chain. It is a no-op when lengthOp is empty.
func buildLengthFilter(chain *filter.Chain, lengthOp string, lengthVal int) error {
	if lengthOp == "" {
		return nil
	}
	f, err := filter.NewLengthFilter(lengthOp, lengthVal)
	if err != nil {
		return fmt.Errorf("--length-op / --length-val: %w", err)
	}
	chain.Add(f)
	return nil
}
