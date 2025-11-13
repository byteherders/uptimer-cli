package main

import "testing"

func TestFormatStatusLine(t *testing.T) {
	line := formatStatusLine("label", "1ms", "2ms", "1ms", "3ms", "graph")
	expect := "label | last 1ms        avg 2ms        min 1ms        max 3ms        |graph"
	if line != expect {
		t.Fatalf("unexpected format: %q", line)
	}
}
