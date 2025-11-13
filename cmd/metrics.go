package main

import (
	"fmt"
	"strings"
	"time"
)

func formatDuration(d time.Duration) string {
	if d <= 0 {
		return "0ms"
	}
	if d >= time.Second {
		return fmt.Sprintf("%.2fs", d.Seconds())
	}
	if d >= time.Millisecond {
		return fmt.Sprintf("%.2fms", float64(d)/float64(time.Millisecond))
	}
	return fmt.Sprintf("%.0fus", float64(d)/float64(time.Microsecond))
}

func padLabel(label string, width int) string {
	if width <= 0 {
		return label
	}
	if len(label) > width {
		if width <= 3 {
			return label[:width]
		}
		return label[:width-3] + "..."
	}
	if len(label) < width {
		return label + strings.Repeat(" ", width-len(label))
	}
	return label
}

func renderGraph(history []time.Duration, width int) string {
	if width <= 0 {
		return ""
	}
	result := make([]rune, width)
	for i := range result {
		result[i] = ' '
	}
	if len(history) == 0 {
		return string(result)
	}

	scale := maxDuration(history)
	if scale <= 0 {
		scale = time.Millisecond
	}

	levels := []rune{' ', '▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}
	start := 0
	if len(history) > width {
		start = len(history) - width
	}

	for i := start; i < len(history); i++ {
		ratio := float64(history[i]) / float64(scale)
		if ratio < 0 {
			ratio = 0
		}
		if ratio > 1 {
			ratio = 1
		}
		idx := int(ratio * float64(len(levels)-1))
		pos := width - (len(history) - i)
		if pos >= 0 && pos < width {
			result[pos] = levels[idx]
		}
	}
	return string(result)
}

func maxDuration(values []time.Duration) time.Duration {
	var max time.Duration
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return max
}
