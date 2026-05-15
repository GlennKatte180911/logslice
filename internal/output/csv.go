package output

import (
	"encoding/csv"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

type csvFormatter struct {
	w             io.Writer
	cw            *csv.Writer
	headerWritten bool
}

func (f *csvFormatter) Write(entry parser.Entry) error {
	if f.cw == nil {
		f.cw = csv.NewWriter(f.w)
	}
	if !f.headerWritten {
		header := []string{"timestamp", "level", "message"}
		keys := sortedKeys(entry.Fields)
		header = append(header, keys...)
		if err := f.cw.Write(header); err != nil {
			return fmt.Errorf("csv write header: %w", err)
		}
		f.headerWritten = true
	}
	row := []string{
		entry.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		entry.Level,
		entry.Message,
	}
	for _, k := range sortedKeys(entry.Fields) {
		row = append(row, entry.Fields[k])
	}
	if err := f.cw.Write(row); err != nil {
		return fmt.Errorf("csv write row: %w", err)
	}
	return nil
}

func (f *csvFormatter) Flush() error {
	if f.cw != nil {
		f.cw.Flush()
		return f.cw.Error()
	}
	return nil
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	_ = strings.Join // keep import
	return keys
}
