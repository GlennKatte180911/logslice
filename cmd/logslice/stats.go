package main

import (
	"fmt"
	"os"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/output"
)

// wrapWithCounter optionally wraps the chain's root filter in a CountFilter
// so that match statistics can be reported at the end of a run.
// It returns the counter (may be nil when stats are disabled).
func wrapWithCounter(ch *filter.Chain, enable bool) *filter.CountFilter {
	if !enable {
		return nil
	}
	// Chain itself satisfies Matcher — wrap the whole chain.
	cf, err := filter.NewCountFilter(ch)
	if err != nil {
		// Should never happen; chain is always non-nil.
		fmt.Fprintf(os.Stderr, "warning: could not create count filter: %v\n", err)
		return nil
	}
	return cf
}

// printStats writes a stats summary to stderr when counter is non-nil.
func printStats(counter *filter.CountFilter, written int64) {
	if counter == nil {
		return
	}
	s := output.Stats{
		Total:   counter.Total(),
		Matched: counter.Matched(),
		Written: written,
	}
	if err := output.WriteStats(os.Stderr, s); err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not write stats: %v\n", err)
	}
}
