package main

import (
	"testing"
	"time"
)

func TestLatencyBufferStatsAndValues(t *testing.T) {
	buf := newLatencyBuffer(3)
	buf.Add(10 * time.Millisecond)
	buf.Add(20 * time.Millisecond)
	buf.Add(30 * time.Millisecond)

	stats := buf.Stats()
	if stats.count != 3 {
		t.Fatalf("expected count 3, got %d", stats.count)
	}
	if stats.avg != 20*time.Millisecond {
		t.Fatalf("expected avg 20ms, got %v", stats.avg)
	}
	if stats.min != 10*time.Millisecond {
		t.Fatalf("expected min 10ms, got %v", stats.min)
	}
	if stats.max != 30*time.Millisecond {
		t.Fatalf("expected max 30ms, got %v", stats.max)
	}

	values := buf.Values()
	expected := []time.Duration{10 * time.Millisecond, 20 * time.Millisecond, 30 * time.Millisecond}
	if len(values) != len(expected) {
		t.Fatalf("expected %d values, got %d", len(expected), len(values))
	}
	for i, v := range values {
		if v != expected[i] {
			t.Fatalf("value[%d] expected %v, got %v", i, expected[i], v)
		}
	}

	buf.Add(40 * time.Millisecond)
	stats = buf.Stats()
	if stats.count != 3 {
		t.Fatalf("after wrap expected count 3, got %d", stats.count)
	}
	if stats.min != 20*time.Millisecond || stats.max != 40*time.Millisecond {
		t.Fatalf("after wrap min/max mismatch: %+v", stats)
	}
}
