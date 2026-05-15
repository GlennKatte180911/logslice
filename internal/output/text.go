package output

import (
	"fmt"
	"io"

	"github.com/yourorg/logslice/internal/parser"
)

type textFormatter struct {
	w io.Writer
}

func (f *textFormatter) Write(entry parser.Entry) error {
	_, err := fmt.Fprintln(f.w, entry.String())
	return err
}

func (f *textFormatter) Flush() error {
	return nil
}
