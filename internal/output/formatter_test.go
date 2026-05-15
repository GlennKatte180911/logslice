package output_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/parser"
)

func makeEntry(ts string, level, msg string) parser.Entry {
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		panic(err)
	}
	return parser.Entry{
		Timestamp: t,
		Level:     level,
		Message:   msg,
		Fields:    map[string]string{"host": "srv1"},
	}
}

func TestNewFormatter_UnknownFormat(t *testing.T) {
	_, err := output.NewFormatter("xml", &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestTextFormatter_Write(t *testing.T) {
	var buf bytes.Buffer
	f, err := output.NewFormatter("text", &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entry := makeEntry("2024-01-15T10:00:00Z", "INFO", "hello world")
	if err := f.Write(entry); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if err := f.Flush(); err != nil {
		t.Fatalf("Flush: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "hello world") {
		t.Errorf("expected message in output, got: %q", got)
	}
}

func TestJSONFormatter_Write(t *testing.T) {
	var buf bytes.Buffer
	f, err := output.NewFormatter("json", &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entry := makeEntry("2024-01-15T10:00:00Z", "ERROR", "something failed")
	if err := f.Write(entry); err != nil {
		t.Fatalf("Write: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, `"level":"ERROR"`) {
		t.Errorf("expected JSON level field, got: %q", got)
	}
	if !strings.Contains(got, `"message":"something failed"`) {
		t.Errorf("expected JSON message field, got: %q", got)
	}
}

func TestCSVFormatter_Write(t *testing.T) {
	var buf bytes.Buffer
	f, err := output.NewFormatter("csv", &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entry := makeEntry("2024-01-15T10:00:00Z", "WARN", "disk low")
	if err := f.Write(entry); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if err := f.Flush(); err != nil {
		t.Fatalf("Flush: %v", err)
	}
	got := buf.String()
	lines := strings.Split(strings.TrimSpace(got), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected header + 1 row, got %d lines: %q", len(lines), got)
	}
	if !strings.HasPrefix(lines[0], "timestamp,level,message") {
		t.Errorf("unexpected header: %q", lines[0])
	}
	if !strings.Contains(lines[1], "disk low") {
		t.Errorf("expected message in row, got: %q", lines[1])
	}
}
