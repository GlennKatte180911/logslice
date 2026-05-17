package main

// config holds all parsed CLI flag values for logslice.
type config struct {
	// input / output
	format string
	output string

	// time range
	from string
	to   string

	// since / until
	since string
	until string

	// sliding window
	windowAnchor string
	windowBefore string
	windowAfter  string

	// pattern / level
	pattern string
	levels  []string

	// field filters
	fieldKey   string
	fieldValue string

	// misc filters
	limit   int
	sampleN int
	dedup   bool
	stats   bool

	// message text filters
	prefix    string
	suffix    string
	contains  string
	lengthOp  string
	lengthVal int

	// tag filter
	tags []string
}
