package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/logslice/internal/output"
)

func TestWriteStats_Basic(t *testing.T) {
	var buf bytes.Buffer
	s := output.Stats{Total: 10, Matched: 4, Written: 4}
	if err := output.WriteStats(&buf, s); err != nil {
		t.Fatalf("WriteStats: %v", err)
	}
	got := buf.String()
	for _, want := range []string{"total=10", "matched=4", "written=4", "skipped=6"} {
		if !strings.Contains(got, want) {
			t.Errorf("output %q missing %q", got, want)
		}
	}
}

func TestWriteStats_AllMatch(t *testing.T) {
	var buf bytes.Buffer
	s := output.Stats{Total: 5, Matched: 5, Written: 5}
	output.WriteStats(&buf, s)
	if !strings.Contains(buf.String(), "skipped=0") {
		t.Errorf("expected skipped=0 in %q", buf.String())
	}
}

func TestWriteStats_NoneMatch(t *testing.T) {
	var buf bytes.Buffer
	s := output.Stats{Total: 7, Matched: 0, Written: 0}
	output.WriteStats(&buf, s)
	if !strings.Contains(buf.String(), "matched=0") {
		t.Errorf("expected matched=0 in %q", buf.String())
	}
	if !strings.Contains(buf.String(), "skipped=7") {
		t.Errorf("expected skipped=7 in %q", buf.String())
	}
}
