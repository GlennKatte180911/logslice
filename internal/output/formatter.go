package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// Format defines the output format for log entries.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
)

// Formatter writes log entries to an io.Writer in a specific format.
type Formatter interface {
	Write(entry parser.Entry) error
	Flush() error
}

// NewFormatter returns a Formatter for the given format string.
func NewFormatter(format string, w io.Writer) (Formatter, error) {
	switch Format(strings.ToLower(format)) {
	case FormatText, "":
		return &textFormatter{w: w}, nil
	case FormatJSON:
		return &jsonFormatter{w: w}, nil
	case FormatCSV:
		return &csvFormatter{w: w, headerWritten: false}, nil
	default:
		return nil, fmt.Errorf("unknown format %q: must be one of text, json, csv", format)
	}
}
