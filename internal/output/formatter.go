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

// SupportedFormats lists all valid format values accepted by NewFormatter.
var SupportedFormats = []Format{FormatText, FormatJSON, FormatCSV}

// Formatter writes log entries to an io.Writer in a specific format.
type Formatter interface {
	Write(entry parser.Entry) error
	Flush() error
}

// NewFormatter returns a Formatter for the given format string.
// An empty format string defaults to text output.
// Returns an error if the format is not one of: text, json, csv.
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

// ValidFormat reports whether the given format string is a supported output format.
func ValidFormat(format string) bool {
	f := Format(strings.ToLower(format))
	for _, supported := range SupportedFormats {
		if f == supported {
			return true
		}
	}
	return false
}
