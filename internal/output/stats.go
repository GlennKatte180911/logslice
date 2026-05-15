package output

import (
	"fmt"
	"io"
)

// Stats holds match statistics gathered during a run.
type Stats struct {
	Total   int64
	Matched int64
	Written int64
}

// WriteStats prints a summary line to w.
func WriteStats(w io.Writer, s Stats) error {
	_, err := fmt.Fprintf(
		w,
		"stats: total=%d matched=%d written=%d skipped=%d\n",
		s.Total,
		s.Matched,
		s.Written,
		s.Total-s.Matched,
	)
	return err
}
