package output

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

type jsonEntry struct {
	Timestamp time.Time         `json:"timestamp"`
	Level     string            `json:"level"`
	Message   string            `json:"message"`
	Fields    map[string]string `json:"fields,omitempty"`
}

type jsonFormatter struct {
	w   io.Writer
	enc *json.Encoder
}

func (f *jsonFormatter) Write(entry parser.Entry) error {
	if f.enc == nil {
		f.enc = json.NewEncoder(f.w)
	}
	je := jsonEntry{
		Timestamp: entry.Timestamp,
		Level:     entry.Level,
		Message:   entry.Message,
		Fields:    entry.Fields,
	}
	if err := f.enc.Encode(je); err != nil {
		return fmt.Errorf("json encode: %w", err)
	}
	return nil
}

func (f *jsonFormatter) Flush() error {
	return nil
}
