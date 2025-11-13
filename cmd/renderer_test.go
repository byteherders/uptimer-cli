package main

import (
	"strings"
	"testing"
)

func TestComposeStatusLinesSingleLine(t *testing.T) {
	var receivedWidth int
	graph := func(width int) string {
		receivedWidth = width
		return strings.Repeat("x", width)
	}
	lines := composeStatusLines("label", "1ms", "2ms", "1ms", "3ms", graph, 200)
	if len(lines) != 1 {
		t.Fatalf("expected single line, got %d", len(lines))
	}
	if receivedWidth != graphWidth {
		t.Fatalf("expected graph width %d, got %d", graphWidth, receivedWidth)
	}
	if !strings.HasSuffix(lines[0], strings.Repeat("x", graphWidth)) {
		t.Fatalf("graph missing in %q", lines[0])
	}
}

func TestComposeStatusLinesWrapsGraph(t *testing.T) {
	var receivedWidth int
	graph := func(width int) string {
		receivedWidth = width
		return strings.Repeat("y", width)
	}
	lines := composeStatusLines("short", "1ms", "2ms", "1ms", "3ms", graph, 20)
	if len(lines) != 2 {
		t.Fatalf("expected two lines, got %d", len(lines))
	}
	if receivedWidth != 20 {
		t.Fatalf("expected wrap width 20, got %d", receivedWidth)
	}
	if lines[1] != strings.Repeat("y", 20) {
		t.Fatalf("graph line mismatch: %q", lines[1])
	}
}
