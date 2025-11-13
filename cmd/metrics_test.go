package main

import (
	"strings"
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	if got := formatDuration(1500 * time.Millisecond); got != "1.50s" {
		t.Fatalf("expected 1.50s, got %s", got)
	}
	if got := formatDuration(250 * time.Millisecond); got != "250.00ms" {
		t.Fatalf("expected 250.00ms, got %s", got)
	}
	if got := formatDuration(500 * time.Microsecond); got != "500us" {
		t.Fatalf("expected 500us, got %s", got)
	}
}

func TestPadLabel(t *testing.T) {
	if got := padLabel("short", 8); got != "short   " {
		t.Fatalf("expected padded string, got %q", got)
	}
	if got := padLabel("this-is-a-long-label", 10); got != "this-is..." {
		t.Fatalf("expected truncated label, got %q", got)
	}
}

func TestRenderGraph(t *testing.T) {
	history := []time.Duration{
		1 * time.Millisecond,
		2 * time.Millisecond,
		0,
	}
	graph := renderGraph(history, 3)
	if len([]rune(graph)) != 3 {
		t.Fatalf("expected graph width 3, got %d", len([]rune(graph)))
	}
	if strings.Count(graph, "█") != 1 {
		t.Fatalf("expected one max block in %q", graph)
	}
	if strings.TrimSpace(graph) == "" {
		t.Fatalf("graph should contain visible blocks")
	}
}
