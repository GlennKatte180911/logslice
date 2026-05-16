package main

import (
	"fmt"

	"github.com/user/logslice/internal/filter"
)

// suffixFlags holds the CLI flag values related to the suffix filter.
type suffixFlags struct {
	suffix          string
	caseInsensitive bool
}

// buildSuffixFilter constructs a suffixFilter from the provided flags and adds
// it to chain. It is a no-op when suffix is empty. Returns an error if the
// filter cannot be constructed.
func buildSuffixFilter(chain *filter.Chain, sf suffixFlags) error {
	if sf.suffix == "" {
		return nil
	}
	f, err := filter.NewSuffixFilter(sf.suffix, sf.caseInsensitive)
	if err != nil {
		return fmt.Errorf("suffix filter: %w", err)
	}
	chain.Add(f)
	return nil
}
