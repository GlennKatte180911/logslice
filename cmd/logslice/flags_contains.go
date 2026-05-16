package main

import (
	"fmt"

	"github.com/user/logslice/internal/filter"
)

// buildContainsFilter constructs a containsFilter from the parsed CLI flags and
// adds it to the chain. It is a no-op when the --contains flag is empty.
func buildContainsFilter(chain interface {
	Add(filter.Filter)
}, substring string, caseSensitive bool) error {
	if substring == "" {
		return nil
	}
	f, err := filter.NewContainsFilter(substring, caseSensitive)
	if err != nil {
		return fmt.Errorf("contains filter: %w", err)
	}
	chain.Add(f)
	return nil
}
