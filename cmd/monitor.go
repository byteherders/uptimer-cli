package main

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

type statusUpdate struct {
	index int
	lines []string
}

var dialTimeout = func(network, address string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout(network, address, timeout)
}

func runMonitors(ctx context.Context, cfg appConfig) {
	renderer := newRenderer(len(cfg.targets))
	initialWidth := detectTerminalWidth()
	for i := range cfg.targets {
		graphFn := func(width int) string {
			if width <= 0 {
				width = graphWidth
			}
			return strings.Repeat(" ", width)
		}
		renderer.blocks[i] = composeStatusLines(
			cfg.targets[i].label,
			"n/a",
			"n/a",
			"n/a",
			"n/a",
			graphFn,
			initialWidth,
		)
	}
	renderer.Init()

	bufferSize := max(graphWidth, historySamples)
	updates := make(chan statusUpdate, len(cfg.targets))

	var wg sync.WaitGroup
	for i := range cfg.targets {
		wg.Add(1)
		go func(idx int, target targetConfig) {
			defer wg.Done()
			monitorTarget(ctx, idx, target, bufferSize, cfg.interval, cfg.timeout, updates)
		}(i, cfg.targets[i])
	}

	go func() {
		wg.Wait()
		close(updates)
	}()

	renderer.Run(ctx, updates)
}

func monitorTarget(ctx context.Context, idx int, cfg targetConfig, bufferSize int, interval, timeout time.Duration, updates chan<- statusUpdate) {
	history := newLatencyBuffer(bufferSize)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		latency, err := measureLatency(cfg.proto, cfg.address, timeout)

		lastText := "timeout"
		avgText := "n/a"
		minText := "n/a"
		maxText := "n/a"
		graphFn := func(width int) string {
			if width <= 0 {
				width = graphWidth
			}
			return strings.Repeat("!", width)
		}

		if err == nil {
			history.Add(latency)
			stats := history.Stats()
			lastText = formatDuration(latency)
			avgText = formatDuration(stats.avg)
			minText = formatDuration(stats.min)
			maxText = formatDuration(stats.max)
			values := history.Values()
			graphFn = func(width int) string {
				if width <= 0 {
					width = graphWidth
				}
				return renderGraph(values, width)
			}
		} else {
			stats := history.Stats()
			if stats.count > 0 {
				avgText = formatDuration(stats.avg)
				minText = formatDuration(stats.min)
				maxText = formatDuration(stats.max)
			}
		}

		lines := composeStatusLines(cfg.label, lastText, avgText, minText, maxText, graphFn, detectTerminalWidth())

		select {
		case <-ctx.Done():
			return
		case updates <- statusUpdate{index: idx, lines: lines}:
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

type graphBuilder func(width int) string

func composeStatusLines(label, last, avg, min, max string, graph graphBuilder, termWidth int) []string {
	statsLine := fmt.Sprintf("%s | last %-10s avg %-10s min %-10s max %-10s |", label, last, avg, min, max)
	if graph == nil {
		return []string{statsLine}
	}

	if termWidth > 0 && utf8.RuneCountInString(statsLine)+graphWidth > termWidth {
		width := termWidth
		if width <= 0 {
			width = graphWidth
		}
		return []string{statsLine, graph(width)}
	}

	return []string{statsLine + graph(graphWidth)}
}

func verifyConnectivity(proto, target string, timeout time.Duration) error {
	conn, err := dialTimeout(proto, target, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

func measureLatency(proto, target string, timeout time.Duration) (time.Duration, error) {
	start := time.Now()
	conn, err := dialTimeout(proto, target, timeout)
	elapsed := time.Since(start)
	if err != nil {
		return elapsed, err
	}
	conn.Close()
	return elapsed, nil
}
