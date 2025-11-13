package main

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type statusUpdate struct {
	index int
	line  string
}

var dialTimeout = func(network, address string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout(network, address, timeout)
}

func runMonitors(ctx context.Context, cfg appConfig) {
	renderer := newRenderer(len(cfg.targets))
	for i := range cfg.targets {
		renderer.lines[i] = formatStatusLine(
			cfg.targets[i].label,
			"n/a",
			"n/a",
			"n/a",
			"n/a",
			strings.Repeat(" ", graphWidth),
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
		graph := strings.Repeat("!", graphWidth)

		if err == nil {
			history.Add(latency)
			stats := history.Stats()
			lastText = formatDuration(latency)
			avgText = formatDuration(stats.avg)
			minText = formatDuration(stats.min)
			maxText = formatDuration(stats.max)
			graph = renderGraph(history.Values(), graphWidth)
		} else {
			stats := history.Stats()
			if stats.count > 0 {
				avgText = formatDuration(stats.avg)
				minText = formatDuration(stats.min)
				maxText = formatDuration(stats.max)
			}
		}

		line := formatStatusLine(cfg.label, lastText, avgText, minText, maxText, graph)

		select {
		case <-ctx.Done():
			return
		case updates <- statusUpdate{index: idx, line: line}:
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func formatStatusLine(label, last, avg, min, max, graph string) string {
	return fmt.Sprintf("%s | last %-10s avg %-10s min %-10s max %-10s |%s", label, last, avg, min, max, graph)
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
